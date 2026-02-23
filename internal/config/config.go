package config

import (
	"errors"
	"fmt"
	"strings"
)

// Config holds all configuration parameters
type Config struct {
	// Scan targets
	Organization string
	User         string
	Token        string

	// Output settings
	OutputFormat string
	OutputFile   string

	// Filtering
	IncludeDorks []string
	ExcludeDorks []string
	Priority     string

	// Performance
	Workers   int
	RateLimit int
	Delay     int

	// Flags
	Verbose      bool
	Quiet        bool
	NoColor      bool
	SaveProgress bool
	ProgressFile string
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Check if either org or user is provided
	if c.Organization == "" && c.User == "" {
		return errors.New("either --org or --user must be provided")
	}

	// Check if both are provided
	if c.Organization != "" && c.User != "" {
		return errors.New("cannot use both --org and --user at the same time")
	}

	// Validate token
	if c.Token == "" {
		return errors.New("--token is required")
	}

	// Validate token format
	if !strings.HasPrefix(c.Token, "ghp_") &&
		!strings.HasPrefix(c.Token, "gho_") &&
		!strings.HasPrefix(c.Token, "ghu_") &&
		!strings.HasPrefix(c.Token, "ghs_") {
		return errors.New("invalid GitHub token format (must start with ghp_, gho_, ghu_, or ghs_)")
	}

	// Validate output format
	validFormats := map[string]bool{
		"terminal": true,
		"json":     true,
		"csv":      true,
		"html":     true,
	}
	if !validFormats[c.OutputFormat] {
		return fmt.Errorf("invalid output format: %s (valid: terminal, json, csv, html)", c.OutputFormat)
	}

	// Validate priority
	validPriorities := map[string]bool{
		"all":    true,
		"high":   true,
		"medium": true,
		"low":    true,
	}
	if !validPriorities[c.Priority] {
		return fmt.Errorf("invalid priority: %s (valid: all, high, medium, low)", c.Priority)
	}

	// Validate workers
	if c.Workers < 1 || c.Workers > 20 {
		return errors.New("workers must be between 1 and 20")
	}

	// Validate rate limit
	if c.RateLimit < 1 || c.RateLimit > 60 {
		return errors.New("rate-limit must be between 1 and 60")
	}

	// Validate delay
	if c.Delay < 0 || c.Delay > 10 {
		return errors.New("delay must be between 0 and 10 seconds")
	}

	// Set default output file if not specified
	if c.OutputFile == "" && c.OutputFormat != "terminal" {
		c.OutputFile = fmt.Sprintf("./dighub-results.%s", c.OutputFormat)
	}

	return nil
}

// GetTarget returns the scan target (org or user)
func (c *Config) GetTarget() string {
	if c.Organization != "" {
		return c.Organization
	}
	return c.User
}

// GetTargetType returns the target type
func (c *Config) GetTargetType() string {
	if c.Organization != "" {
		return "org"
	}
	return "user"
}
