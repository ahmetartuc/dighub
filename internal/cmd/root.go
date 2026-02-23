package cmd

import (
	"fmt"
	"strings"

	"github.com/ahmetartuc/dighub/internal/config"
	"github.com/ahmetartuc/dighub/internal/logger"
	"github.com/ahmetartuc/dighub/internal/output"
	"github.com/ahmetartuc/dighub/internal/scanner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	cfg     *config.Config
	version string
)

var rootCmd = &cobra.Command{
	Use:   "dighub",
	Short: "Advanced GitHub Dorking & Secret Hunting Tool",
	Long: `Dighub is a powerful CLI tool that performs advanced GitHub dorking 
to detect exposed secrets, credentials, webhooks and sensitive files 
inside public repositories.`,
	Version: version,
	RunE:    run,
}

func init() {
	cfg = &config.Config{}

	rootCmd.Flags().StringVarP(&cfg.Organization, "org", "o", "", "GitHub organization to scan (required)")
	rootCmd.Flags().StringVarP(&cfg.Token, "token", "t", "", "GitHub Personal Access Token (required)")
	rootCmd.Flags().StringVarP(&cfg.User, "user", "u", "", "GitHub user to scan (alternative to org)")

	rootCmd.Flags().StringVarP(&cfg.OutputFormat, "output", "f", "terminal", "Output format: terminal, json, csv, html")
	rootCmd.Flags().StringVarP(&cfg.OutputFile, "out-file", "w", "", "Output file path (default: ./dighub-results.[format])")

	rootCmd.Flags().StringSliceVarP(&cfg.IncludeDorks, "include", "i", []string{}, "Include specific dorks (comma-separated)")
	rootCmd.Flags().StringSliceVarP(&cfg.ExcludeDorks, "exclude", "e", []string{}, "Exclude specific dorks (comma-separated)")
	rootCmd.Flags().StringVarP(&cfg.Priority, "priority", "p", "all", "Priority level: all, high, medium, low")

	rootCmd.Flags().IntVarP(&cfg.Workers, "workers", "W", 5, "Number of concurrent workers")
	rootCmd.Flags().IntVarP(&cfg.RateLimit, "rate-limit", "r", 30, "Requests per minute")
	rootCmd.Flags().IntVarP(&cfg.Delay, "delay", "d", 2, "Delay between requests (seconds)")

	rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().BoolVarP(&cfg.Quiet, "quiet", "q", false, "Quiet mode (only show matches)")
	rootCmd.Flags().BoolVarP(&cfg.NoColor, "no-color", "n", false, "Disable colored output")
	rootCmd.Flags().BoolVarP(&cfg.SaveProgress, "save-progress", "s", false, "Save progress to resume later")
	rootCmd.Flags().StringVar(&cfg.ProgressFile, "progress-file", ".dighub-progress.json", "Progress file path")

	rootCmd.MarkFlagRequired("token")
}

func Execute(ver string) error {
	version = ver
	rootCmd.Version = ver
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	if cfg.NoColor {
		color.NoColor = true
	}

	log := logger.New(cfg.Verbose, cfg.Quiet)

	if !cfg.Quiet {
		printBanner(version)
	}

	log.Info("Initializing scanner...")
	scn, err := scanner.New(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to initialize scanner: %w", err)
	}

	log.Info("Starting scan...")
	target := cfg.Organization
	if target == "" {
		target = cfg.User
	}
	log.Info("Target: %s", target)
	log.Info("Workers: %d", cfg.Workers)
	log.Info("Output format: %s", cfg.OutputFormat)

	results, err := scn.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	log.Info("Generating output...")
	outputHandler := output.New(cfg)
	if err := outputHandler.Write(results); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	if !cfg.Quiet {
		printSummary(results, log)
	}

	return nil
}

func printBanner(ver string) {
	banner := `
    ____  _       __  __      __  
   / __ \(_)___ _/ / / /_  __/ /_ 
  / / / / / __ '/ /_/ / / / / __ \
 / /_/ / / /_/ / __  / /_/ / /_/ /
/_____/_/\__, /_/ /_/\__,_/_.___/ 
        /____/                     
`
	fmt.Println(color.CyanString(banner))
	fmt.Printf("%s %s\n", color.YellowString("Version:"), ver)
	fmt.Printf("%s Advanced GitHub Dorking & Secret Hunting Tool\n", color.GreenString("=>"))
	fmt.Println(strings.Repeat("─", 50))
}

func printSummary(results *scanner.ScanResults, log *logger.Logger) {
	fmt.Println()
	fmt.Println(strings.Repeat("═", 50))
	fmt.Println(color.CyanString("SCAN SUMMARY"))
	fmt.Println(strings.Repeat("═", 50))
	
	fmt.Printf("%s %d\n", color.BlueString("Total Dorks Scanned:"), results.TotalDorks)
	fmt.Printf("%s %d\n", color.GreenString("Matches Found:"), results.TotalMatches)
	fmt.Printf("%s %d\n", color.YellowString("Unique Files:"), results.UniqueFiles)
	fmt.Printf("%s %d\n", color.RedString("High Priority:"), results.HighPriority)
	fmt.Printf("%s %d\n", color.MagentaString("Medium Priority:"), results.MediumPriority)
	fmt.Printf("%s %d\n", color.CyanString("Low Priority:"), results.LowPriority)
	fmt.Printf("%s %s\n", color.BlueString("Duration:"), results.Duration)
	
	fmt.Println(strings.Repeat("═", 50))
	
	if results.TotalMatches > 0 {
		log.Success("✓ Scan completed successfully!")
		if results.OutputFile != "" {
			log.Success("Results saved to: %s", results.OutputFile)
		}
	} else {
		log.Info("No sensitive data found in the target.")
	}
}
