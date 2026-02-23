package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ahmetartuc/dighub/internal/config"
	"github.com/ahmetartuc/dighub/internal/dorks"
	"github.com/ahmetartuc/dighub/internal/logger"
	"github.com/schollz/progressbar/v3"
)

// Scanner handles the GitHub code search operations
type Scanner struct {
	config *config.Config
	logger *logger.Logger
	client *http.Client
	bar    *progressbar.ProgressBar
}

// Match represents a single search result
type Match struct {
	Dork       dorks.Dork
	URL        string
	Repository string
	Path       string
	Score      float64
	Timestamp  time.Time
}

// ScanResults holds the aggregated scan results
type ScanResults struct {
	Matches        []Match
	TotalDorks     int
	TotalMatches   int
	UniqueFiles    int
	HighPriority   int
	MediumPriority int
	LowPriority    int
	Duration       string
	OutputFile     string
	Target         string
	ScanDate       time.Time
}

// GitHubSearchResponse represents the GitHub API response
type GitHubSearchResponse struct {
	TotalCount int `json:"total_count"`
	Items      []struct {
		Name       string `json:"name"`
		Path       string `json:"path"`
		HTMLURL    string `json:"html_url"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
		Score float64 `json:"score"`
	} `json:"items"`
}

// New creates a new scanner instance
func New(cfg *config.Config, log *logger.Logger) (*Scanner, error) {
	return &Scanner{
		config: cfg,
		logger: log,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Scan performs the complete scanning operation
func (s *Scanner) Scan() (*ScanResults, error) {
	startTime := time.Now()

	// Get and filter dorks
	allDorks := dorks.GetDorks()
	filteredDorks := s.filterDorks(allDorks)

	if len(filteredDorks) == 0 {
		return nil, fmt.Errorf("no dorks match the specified filters")
	}

	s.logger.Info("Loaded %d dorks for scanning", len(filteredDorks))

	// Initialize progress bar
	if !s.config.Quiet {
		s.bar = progressbar.NewOptions(len(filteredDorks),
			progressbar.OptionSetDescription("Scanning"),
			progressbar.OptionSetWidth(40),
			progressbar.OptionShowCount(),
			progressbar.OptionShowIts(),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "█",
				SaucerPadding: "░",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)
	}

	// Scan with concurrency
	matches := s.scanConcurrent(filteredDorks)

	// Calculate statistics
	results := s.buildResults(matches, filteredDorks, startTime)

	return results, nil
}

// filterDorks applies all filters to the dork list
func (s *Scanner) filterDorks(allDorks []dorks.Dork) []dorks.Dork {
	filtered := allDorks

	// Filter by priority
	filtered = dorks.FilterByPriority(filtered, s.config.Priority)

	// Filter by include patterns
	filtered = dorks.FilterByInclude(filtered, s.config.IncludeDorks)

	// Filter by exclude patterns
	filtered = dorks.FilterByExclude(filtered, s.config.ExcludeDorks)

	return filtered
}

// scanConcurrent performs concurrent scanning with worker pool
func (s *Scanner) scanConcurrent(dorkList []dorks.Dork) []Match {
	var (
		matches []Match
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	// Create worker pool
	jobs := make(chan dorks.Dork, len(dorkList))

	// Start workers
	for i := 0; i < s.config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for dork := range jobs {
				// Perform search
				dorkMatches := s.searchDork(dork)

				// Add matches
				if len(dorkMatches) > 0 {
					mu.Lock()
					matches = append(matches, dorkMatches...)
					mu.Unlock()
				}

				// Update progress bar
				if s.bar != nil {
					s.bar.Add(1)
				}

				// Rate limiting delay
				time.Sleep(time.Duration(s.config.Delay) * time.Second)
			}
		}()
	}

	// Send jobs to workers
	for _, dork := range dorkList {
		jobs <- dork
	}
	close(jobs)

	// Wait for completion
	wg.Wait()

	if s.bar != nil {
		s.bar.Finish()
		fmt.Println() // New line after progress bar
	}

	return matches
}

// searchDork performs a single dork search
func (s *Scanner) searchDork(dork dorks.Dork) []Match {
	// Build query
	targetType := s.config.GetTargetType()
	target := s.config.GetTarget()
	query := fmt.Sprintf("%s:%s %s", targetType, target, dork.Pattern)

	s.logger.Debug("Searching: %s", query)

	// Perform search with retry
	var response *GitHubSearchResponse
	var err error

	for retries := 0; retries < 3; retries++ {
		response, err = s.performSearch(query)
		if err == nil {
			break
		}

		// Check if it's a rate limit error
		if err.Error() == "rate_limit" {
			s.handleRateLimit()
			continue
		}

		s.logger.Debug("Search failed (attempt %d/3): %v", retries+1, err)
		time.Sleep(time.Duration(retries+1) * time.Second)
	}

	if err != nil {
		s.logger.Debug("Search failed for dork '%s': %v", dork.Pattern, err)
		return nil
	}

	// Process results
	if response.TotalCount == 0 {
		s.logger.NoMatch(dork.Pattern)
		return nil
	}

	// Convert to matches
	matches := make([]Match, 0, len(response.Items))
	urls := make([]string, 0, len(response.Items))

	for _, item := range response.Items {
		match := Match{
			Dork:       dork,
			URL:        item.HTMLURL,
			Repository: item.Repository.FullName,
			Path:       item.Path,
			Score:      item.Score,
			Timestamp:  time.Now(),
		}
		matches = append(matches, match)
		urls = append(urls, item.HTMLURL)
	}

	// Log the match
	s.logger.Match(dork.Pattern, len(matches), urls)

	return matches
}

// performSearch executes the GitHub API search request
func (s *Scanner) performSearch(query string) (*GitHubSearchResponse, error) {
	urlStr := fmt.Sprintf("https://api.github.com/search/code?q=%s&per_page=100",
		url.QueryEscape(query))

	req, err := http.NewRequestWithContext(context.Background(), "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+s.config.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for rate limiting
	if resp.StatusCode == 403 {
		body, _ := io.ReadAll(resp.Body)
		if contains(string(body), "rate limit") {
			return nil, fmt.Errorf("rate_limit")
		}
		return nil, fmt.Errorf("API error 403: %s", string(body))
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var searchResp GitHubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}

// handleRateLimit waits until the rate limit is reset
func (s *Scanner) handleRateLimit() {
	s.logger.Debug("Rate limit hit, checking reset time...")

	req, _ := http.NewRequest("GET", "https://api.github.com/rate_limit", nil)
	req.Header.Set("Authorization", "token "+s.config.Token)

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Warning("Could not fetch rate limit info, waiting 60 seconds...")
		time.Sleep(60 * time.Second)
		return
	}
	defer resp.Body.Close()

	var data struct {
		Resources struct {
			Search struct {
				Reset int64 `json:"reset"`
			} `json:"search"`
		} `json:"resources"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		s.logger.Warning("Could not parse rate limit info, waiting 60 seconds...")
		time.Sleep(60 * time.Second)
		return
	}

	resetTime := time.Unix(data.Resources.Search.Reset, 0)
	waitDuration := time.Until(resetTime) + 5*time.Second

	s.logger.RateLimit(resetTime, waitDuration)
	time.Sleep(waitDuration)
}

// buildResults creates the final results structure
func (s *Scanner) buildResults(matches []Match, dorkList []dorks.Dork, startTime time.Time) *ScanResults {
	// Count unique files
	uniqueFiles := make(map[string]bool)
	for _, match := range matches {
		uniqueFiles[match.Repository+"/"+match.Path] = true
	}

	// Count by priority
	highCount := 0
	mediumCount := 0
	lowCount := 0

	for _, match := range matches {
		switch match.Dork.Priority {
		case dorks.PriorityHigh:
			highCount++
		case dorks.PriorityMedium:
			mediumCount++
		case dorks.PriorityLow:
			lowCount++
		}
	}

	return &ScanResults{
		Matches:        matches,
		TotalDorks:     len(dorkList),
		TotalMatches:   len(matches),
		UniqueFiles:    len(uniqueFiles),
		HighPriority:   highCount,
		MediumPriority: mediumCount,
		LowPriority:    lowCount,
		Duration:       time.Since(startTime).Round(time.Second).String(),
		Target:         s.config.GetTarget(),
		ScanDate:       time.Now(),
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// indexOf finds the index of a substring
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
