package config

import (
	"errors"
	"fmt"
)

type Config struct {
	Organization string
	User         string
	Token        string

	OutputFormat string
	OutputFile   string

	IncludeDorks []string
	ExcludeDorks []string
	Priority     string

	Workers   int
	RateLimit int
	Delay     int

	Verbose      bool
	Quiet        bool
	NoColor      bool
	SaveProgress bool
	ProgressFile string
}

func (c *Config) Validate() error {
	if c.Organization == "" && c.User == "" {
		return errors.New("either --org or --user must be provided")
	}

	if c.Organization != "" && c.User != "" {
		return errors.New("cannot use both --org and --user at the same time")
	}

	if c.Token == "" {
		return errors.New("--token is required")
	}

	validFormats := map[string]bool{
		"terminal": true,
		"json":     true,
		"csv":      true,
		"html":     true,
	}
	if !validFormats[c.OutputFormat] {
		return fmt.Errorf("invalid output format: %s (valid: terminal, json, csv, html)", c.OutputFormat)
	}

	validPriorities := map[string]bool{
		"all":    true,
		"high":   true,
		"medium": true,
		"low":    true,
	}
	if !validPriorities[c.Priority] {
		return fmt.Errorf("invalid priority: %s (valid: all, high, medium, low)", c.Priority)
	}

	if c.Workers < 1 || c.Workers > 20 {
		return errors.New("workers must be between 1 and 20")
	}

	if c.RateLimit < 1 || c.RateLimit > 60 {
		return errors.New("rate-limit must be between 1 and 60")
	}

	if c.Delay < 0 || c.Delay > 10 {
		return errors.New("delay must be between 0 and 10 seconds")
	}

	if c.OutputFile == "" && c.OutputFormat != "terminal" {
		c.OutputFile = fmt.Sprintf("./dighub-results.%s", c.OutputFormat)
	}

	return nil
}

func (c *Config) GetTarget() string {
	if c.Organization != "" {
		return c.Organization
	}
	return c.User
}

func (c *Config) GetTargetType() string {
	if c.Organization != "" {
		return "org"
	}
	return "user"
}
