package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

// Logger handles all logging operations
type Logger struct {
	verbose bool
	quiet   bool
}

// New creates a new logger instance
func New(verbose, quiet bool) *Logger {
	return &Logger{
		verbose: verbose,
		quiet:   quiet,
	}
}

// Info logs informational messages
func (l *Logger) Info(format string, args ...interface{}) {
	if l.quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s %s\n",
		color.BlueString("[%s]", timestamp),
		color.CyanString("[INFO]"),
		msg,
	)
}

// Success logs success messages
func (l *Logger) Success(format string, args ...interface{}) {
	if l.quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s %s\n",
		color.BlueString("[%s]", timestamp),
		color.GreenString("[✓]"),
		color.GreenString(msg),
	)
}

// Warning logs warning messages
func (l *Logger) Warning(format string, args ...interface{}) {
	if l.quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s %s\n",
		color.BlueString("[%s]", timestamp),
		color.YellowString("[!]"),
		color.YellowString(msg),
	)
}

// Error logs error messages
func (l *Logger) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	fmt.Fprintf(os.Stderr, "%s %s %s\n",
		color.BlueString("[%s]", timestamp),
		color.RedString("[✗]"),
		color.RedString(msg),
	)
}

// Debug logs debug messages (only in verbose mode)
func (l *Logger) Debug(format string, args ...interface{}) {
	if !l.verbose || l.quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s %s\n",
		color.BlueString("[%s]", timestamp),
		color.MagentaString("[DEBUG]"),
		color.HiBlackString(msg),
	)
}

// Match logs matched dorks (always shown unless quiet)
func (l *Logger) Match(dork string, count int, urls []string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s %s (%s)\n",
		color.BlueString("[%s]", timestamp),
		color.GreenString("[+]"),
		color.GreenString("Dork matched: %s", dork),
		color.YellowString("%d results", count),
	)

	if !l.quiet && !l.verbose {
		// Show first 3 URLs in normal mode
		limit := 3
		if count < limit {
			limit = count
		}
		for i := 0; i < limit; i++ {
			fmt.Printf("    %s %s\n", color.HiBlackString("→"), urls[i])
		}
		if count > limit {
			fmt.Printf("    %s\n", color.HiBlackString("... and %d more", count-limit))
		}
	} else if l.verbose {
		// Show all URLs in verbose mode
		for _, url := range urls {
			fmt.Printf("    %s %s\n", color.HiBlackString("→"), url)
		}
	}
}

// NoMatch logs when a dork has no matches (only in verbose mode)
func (l *Logger) NoMatch(dork string) {
	if !l.verbose || l.quiet {
		return
	}
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s %s\n",
		color.BlueString("[%s]", timestamp),
		color.HiBlackString("[-]"),
		color.HiBlackString("No results: %s", dork),
	)
}

// Progress logs progress updates
func (l *Logger) Progress(current, total int) {
	if l.quiet {
		return
	}
	percentage := float64(current) / float64(total) * 100
	fmt.Printf("\r%s Progress: %d/%d (%.1f%%)  ",
		color.CyanString("[SCAN]"),
		current,
		total,
		percentage,
	)
}

// RateLimit logs rate limit warnings
func (l *Logger) RateLimit(resetTime time.Time, waitDuration time.Duration) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s Rate limit reached. Waiting until %s (%s)\n",
		color.BlueString("[%s]", timestamp),
		color.YellowString("[⏳]"),
		resetTime.Format("15:04:05"),
		waitDuration.Round(time.Second),
	)
}

// Fatal logs a fatal error and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Error(format, args...)
	os.Exit(1)
}
