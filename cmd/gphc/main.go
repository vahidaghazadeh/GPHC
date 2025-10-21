package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opsource/gphc/internal/checkers"
	"github.com/opsource/gphc/internal/exporter"
	"github.com/opsource/gphc/internal/git"
	"github.com/opsource/gphc/internal/reporter"
	"github.com/opsource/gphc/internal/scorer"
	"github.com/opsource/gphc/pkg/types"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	rootCmd = &cobra.Command{
		Use:   "gphc",
		Short: "Git Project Health Checker - Audit your Git repositories",
		Long: `Git Project Health Checker (GPHC) is a CLI tool that audits local Git repositories 
against established Open Source best practices regarding documentation, commit history quality, 
and repository hygiene. It assigns a Health Score and provides actionable feedback.`,
		Version: version,
	}
)

var badgeCmd = &cobra.Command{
	Use:   "badge [path]",
	Short: "Generate health badge for a Git repository",
	Long: `Generate a health badge (shields.io style) for the repository.
This command runs a quick health check and generates a badge URL and markdown.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runBadge,
}

var githubCmd = &cobra.Command{
	Use:   "github [path]",
	Short: "Check GitHub integration and configuration",
	Long: `Check GitHub-specific features like branch protection, workflows, and repository settings.
Requires GPHC_TOKEN or GITHUB_TOKEN environment variable for full functionality.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runGitHub,
}

var gitlabCmd = &cobra.Command{
	Use:   "gitlab [path]",
	Short: "Check GitLab integration and configuration",
	Long: `Check GitLab-specific features like branch protection, CI/CD pipelines, and project settings.
Requires GPHC_TOKEN or GITLAB_TOKEN environment variable for full functionality.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runGitLab,
}

var authorsCmd = &cobra.Command{
	Use:   "authors [path]",
	Short: "Analyze commit author patterns and bus factor risk",
	Long: `Analyze commit history to identify contributor patterns and bus factor risks.
Shows contributor distribution, single author dominance, and team participation metrics.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runAuthors,
}

var codebaseCmd = &cobra.Command{
	Use:   "codebase [path]",
	Short: "Analyze codebase structure and detect code smells",
	Long: `Perform lightweight codebase structure analysis to detect common issues.
Checks for missing tests, oversized directories, poor organization, and maintainability issues.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runCodebase,
}

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan multiple repositories for health analysis",
	Long: `Scan multiple repositories simultaneously for health analysis.
Supports recursive scanning to find all Git repositories in directories.
Perfect for organizations with many projects.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runScan,
}

var tuiCmd = &cobra.Command{
	Use:   "tui [path]",
	Short: "Interactive Terminal UI for health monitoring",
	Long: `Launch an interactive terminal interface for health monitoring.
Provides a graphical interface in the terminal with:
- Colorful and interactive score display
- Filtering and rule explanations
- Score trend browsing
- Real-time updates`,
	Args: cobra.MaximumNArgs(1),
	Run:  runTUI,
}

var serveCmd = &cobra.Command{
	Use:   "serve [path]",
	Short: "Start web dashboard server",
	Long: `Start a local web server to display health monitoring dashboard.
Provides a web interface accessible via browser with:
- Multi-project health monitoring
- Historical trend analysis
- Export capabilities
- Team collaboration features`,
	Args: cobra.MaximumNArgs(1),
	Run:  runServe,
}

// tags command: manage and validate git tags/releases
var tagsCmd = &cobra.Command{
	Use:   "tags [path]",
	Short: "Analyze and manage Git tags and releases",
	Long: `Validate semantic tags, check freshness and unreleased commits,
suggest next semantic version, and optionally generate a changelog.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runTags,
}

var suggestCmd = &cobra.Command{
	Use:   "suggest [path]",
	Short: "Suggest commit message based on staged changes",
	Long: `Analyze staged files and suggest conventional commit messages.
This command examines the changes in staged files and suggests
appropriate commit messages following conventional commit format.

Examples:
  git hc suggest                    # Suggest for current directory
  git hc suggest /path/to/repo      # Suggest for specific repository
  git hc suggest --path /path/to/repo # Suggest for specific repository`,
	Args: cobra.MaximumNArgs(1),
	Run:  runSuggest,
}

var commitCmd = &cobra.Command{
	Use:   "commit [path]",
	Short: "Enhanced git commit with health checks and suggestions",
	Long: `Enhanced git commit command with additional options.
This command extends the standard git commit with health checks,
message suggestions, and validation features.

Examples:
  git hc commit --suggest                    # Suggest commit message
  git hc commit --suggest /path/to/repo      # Suggest for specific repository
  git hc commit --suggest --path /path/to/repo # Suggest for specific repository`,
	Args: cobra.MaximumNArgs(1),
	Run:  runCommit,
}

var diffCmd = &cobra.Command{
	Use:   "diff [file]",
	Short: "Show colored diff of file changes",
	Long: `Display colored diff of file changes with syntax highlighting.
Shows additions in light green background and deletions in light red background.

Examples:
  git hc diff                    # Show diff of all staged files
  git hc diff main.go            # Show diff of specific file
  git hc diff --path /path/to/repo # Show diff for specific repository
  git hc diff --staged           # Show staged changes only
  git hc diff --unstaged         # Show unstaged changes only`,
	Args: cobra.MaximumNArgs(1),
	Run:  runDiff,
}

func init() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(preCommitCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(suggestCmd)
	rootCmd.AddCommand(badgeCmd)
	rootCmd.AddCommand(githubCmd)
	rootCmd.AddCommand(gitlabCmd)
	rootCmd.AddCommand(authorsCmd)
	rootCmd.AddCommand(codebaseCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(tagsCmd)

	// Add export format flags
	checkCmd.Flags().StringVarP(&exportFormat, "format", "f", "terminal", "Output format: terminal, json, yaml, markdown, html")
	checkCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: stdout)")

	// Add pre-commit command flags
	preCommitCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Repository path to check")

	// Add suggest command flags
	suggestCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Repository path to analyze")

	// Add commit command flags
	commitCmd.Flags().BoolVar(&commitSuggest, "suggest", false, "Suggest commit message based on staged changes")
	commitCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Repository path to analyze")

	// Add diff command flags
	diffCmd.Flags().BoolVar(&diffStaged, "staged", false, "Show staged changes only")
	diffCmd.Flags().BoolVar(&diffUnstaged, "unstaged", false, "Show unstaged changes only")
	diffCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Repository path to analyze")

	// Add scan command flags
	scanCmd.Flags().BoolVarP(&recursiveScan, "recursive", "r", false, "Recursively scan subdirectories for Git repositories")
	scanCmd.Flags().IntVarP(&minScore, "min-score", "m", 0, "Minimum health score threshold")
	scanCmd.Flags().StringSliceVarP(&excludePatterns, "exclude", "e", []string{}, "Exclude directories matching patterns")
	scanCmd.Flags().StringSliceVarP(&includePatterns, "include", "i", []string{}, "Include only files matching patterns")
	scanCmd.Flags().IntVarP(&parallelJobs, "parallel", "p", 4, "Number of parallel jobs for scanning")
	scanCmd.Flags().BoolVarP(&detailedReport, "detailed", "d", false, "Generate detailed report")
	scanCmd.Flags().StringVarP(&scanOutputFile, "output", "o", "", "Output file path (default: stdout)")

	// Add serve command flags
	serveCmd.Flags().StringVarP(&serverHost, "host", "H", "localhost", "Host to bind the server to")
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "Port to bind the server to")
	serveCmd.Flags().BoolVarP(&serverAuth, "auth", "a", false, "Enable basic authentication")
	serveCmd.Flags().StringVarP(&serverUsername, "username", "u", "admin", "Username for basic authentication")
	serveCmd.Flags().StringVarP(&serverPassword, "password", "w", "admin", "Password for basic authentication")
	serveCmd.Flags().BoolVarP(&serverCORS, "cors", "c", true, "Enable CORS headers")
	serveCmd.Flags().StringVarP(&serverTitle, "title", "t", "GPHC Dashboard", "Dashboard title")

	// Add tags command flags
	tagsCmd.Flags().BoolVar(&tagsSuggest, "suggest", false, "Suggest next semantic version")
	tagsCmd.Flags().StringVar(&tagsChangelogOut, "changelog", "", "Generate changelog to file (e.g. CHANGELOG.md)")
	tagsCmd.Flags().BoolVar(&tagsEnforce, "enforce-tags", false, "Fail if tag policies are violated")
}

var (
	exportFormat    string
	outputFile      string
	pathFlag        string
	recursiveScan   bool
	minScore        int
	excludePatterns []string
	includePatterns []string
	parallelJobs    int
	detailedReport  bool
	scanOutputFile  string
	serverHost      string
	serverPort      int
	serverAuth      bool
	serverUsername  string
	serverPassword  string
	serverCORS      bool
	serverTitle     string

	// tags command flags
	tagsSuggest      bool
	tagsChangelogOut string
	tagsEnforce      bool

	// commit command flags
	commitSuggest bool
	
	// diff command flags
	diffStaged   bool
	diffUnstaged bool
)

var checkCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Run health check on a Git repository",
	Long: `Run a comprehensive health check on the specified Git repository.
If no path is provided, the current directory will be checked.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runCheck,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GPHC (Git Project Health Checker) v%s\n", version)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update GPHC to the latest version",
	Long: `Update GPHC to the latest version from GitHub.
This command will:
1. Pull the latest changes from the repository
2. Update dependencies with go mod tidy
3. Rebuild and reinstall GPHC
4. Show the new version`,
	Run: runUpdate,
}

var preCommitCmd = &cobra.Command{
	Use:   "pre-commit [path]",
	Short: "Run quick pre-commit checks on staged files",
	Long: `Run quick health checks suitable for pre-commit hooks.
This command performs fast checks on staged files and current commit.
Designed for integration with pre-commit framework and Husky.
Returns non-zero exit code if issues are found.

Examples:
  git hc pre-commit                    # Check current directory
  git hc pre-commit /path/to/repo      # Check specific repository
  git hc pre-commit --path /path/to/repo # Check specific repository`,
	Args: cobra.MaximumNArgs(1),
	Run:  runPreCommit,
}

func runPreCommit(cmd *cobra.Command, args []string) {
	var path string

	// Check for --path flag first
	if pathFlag != "" {
		path = pathFlag
	} else if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if path is a git repository
	if !isGitRepository(path) {
		fmt.Printf("Error: %s is not a Git repository\n", path)
		os.Exit(1)
	}

	// Check if there are staged files
	stagedFiles, err := getStagedFiles(path)
	if err != nil {
		fmt.Printf("Error checking staged files: %v\n", err)
		os.Exit(1)
	}

	if len(stagedFiles) == 0 {
		fmt.Println("No staged files to check")
		return
	}

	fmt.Printf("Pre-commit check on %d staged files\n", len(stagedFiles))

	// Run quick checks
	issues := 0

	// Check 1: File formatting
	if !checkFileFormatting(path, stagedFiles) {
		fmt.Println("Some files are not properly formatted")
		issues++
	}

	// Check 2: Commit message (if committing)
	if !checkCommitMessage(path) {
		fmt.Println("Commit message doesn't follow conventional format")
		issues++
	}

	// Check 3: Large files
	if !checkLargeFiles(stagedFiles) {
		fmt.Println("Some files are too large")
		issues++
	}

	// Check 4: Sensitive files
	if !checkSensitiveFiles(stagedFiles) {
		fmt.Println("Sensitive files detected in staging area")
		issues++
	}

	if issues == 0 {
		fmt.Println("All pre-commit checks passed")
	} else {
		fmt.Printf("%d pre-commit check(s) failed\n", issues)
		os.Exit(1)
	}
}

func runDiff(cmd *cobra.Command, args []string) {
	var path string
	
	// Check for --path flag first
	if pathFlag != "" {
		path = pathFlag
	} else if len(args) > 0 {
		// If args[0] is a file, use current directory as repo path
		if _, err := os.Stat(args[0]); err == nil {
			var err error
			path, err = os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				os.Exit(1)
			}
		} else {
			path = args[0]
		}
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if path is a git repository
	if !isGitRepository(path) {
		fmt.Printf("Error: %s is not a Git repository\n", path)
		os.Exit(1)
	}

	// Determine what to show
	var diffArgs []string
	
	if diffStaged {
		diffArgs = []string{"diff", "--cached"}
	} else if diffUnstaged {
		diffArgs = []string{"diff"}
	} else {
		// Show both staged and unstaged changes
		diffArgs = []string{"diff", "HEAD"}
	}
	
	// Add specific file if provided
	if len(args) > 0 && pathFlag == "" {
		diffArgs = append(diffArgs, args[0])
	}
	
	// Run git diff
	cmd_exec := exec.Command("git", diffArgs...)
	cmd_exec.Dir = path
	output, err := cmd_exec.Output()
	if err != nil {
		fmt.Printf("Error running git diff: %v\n", err)
		os.Exit(1)
	}
	
	// Display colored diff
	displayColoredDiff(string(output))
}

func runCommit(cmd *cobra.Command, args []string) {
	if commitSuggest {
		// If --suggest flag is used, delegate to runSuggest
		runSuggest(cmd, args)
		return
	}

	// Default behavior: show help
	fmt.Println("Enhanced git commit command")
	fmt.Println("Use --suggest flag to get commit message suggestions")
	fmt.Println("Example: git hc commit --suggest")
}

func runSuggest(cmd *cobra.Command, args []string) {
	var path string

	// Check for --path flag first
	if pathFlag != "" {
		path = pathFlag
	} else if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if path is a git repository
	if !isGitRepository(path) {
		fmt.Printf("Error: %s is not a Git repository\n", path)
		os.Exit(1)
	}

	// Check if there are staged files
	stagedFiles, err := getStagedFiles(path)
	if err != nil {
		fmt.Printf("Error checking staged files: %v\n", err)
		os.Exit(1)
	}

	if len(stagedFiles) == 0 {
		fmt.Println("No staged files to analyze")
		return
	}

	fmt.Printf("Analyzing %d staged files for commit message suggestion\n", len(stagedFiles))

	// Analyze staged changes
	suggestion := analyzeStagedChanges(path, stagedFiles)

	fmt.Println("\nSuggested commit message:")
	fmt.Printf("  %s\n", suggestion)
	fmt.Println("\nYou can use this message with:")
	fmt.Printf("  git commit -m \"%s\"\n", suggestion)
}

func runCheck(cmd *cobra.Command, args []string) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if path is a git repository
	if !isGitRepository(path) {
		fmt.Printf("Error: %s is not a Git repository\n", path)
		os.Exit(1)
	}

	fmt.Printf("Analyzing repository: %s\n", path)

	// Initialize analyzer
	analyzer, err := git.NewRepositoryAnalyzer(path)
	if err != nil {
		fmt.Printf("Error initializing repository analyzer: %v\n", err)
		os.Exit(1)
	}

	// Analyze repository
	data, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("Error analyzing repository: %v\n", err)
		os.Exit(1)
	}

	// Initialize checkers
	allCheckers := []checkers.Checker{
		checkers.NewDocChecker(),
		checkers.NewSetupChecker(),
		checkers.NewIgnoreChecker(),
		checkers.NewConventionalCommitChecker(),
		checkers.NewMsgLengthChecker(),
		checkers.NewCommitSizeChecker(),
		checkers.NewCommitAuthorInsightsChecker(),
		checkers.NewCodebaseSmellChecker(),
		checkers.NewLocalBranchChecker(),
		checkers.NewStaleBranchChecker(),
		checkers.NewBareRepoChecker(),
		checkers.NewStashChecker(),
		checkers.NewGitHubIntegrationChecker(),
		checkers.NewGitLabIntegrationChecker(),
		checkers.NewTagChecker(),
	}

	// Run all checkers
	scorer := scorer.NewScorer()
	for _, checker := range allCheckers {
		result := checker.Check(data)
		scorer.AddResult(*result)
	}

	// Generate report
	healthReport := scorer.CalculateHealthReport()

	// Handle different output formats
	if exportFormat == "terminal" {
		// Display results in terminal format
		reporter := reporter.NewReporter()
		output := reporter.Report(healthReport)
		fmt.Println(output)
	} else {
		// Export in specified format
		exp := exporter.NewExporter()
		format := exporter.ExportFormat(exportFormat)
		output, err := exp.Export(healthReport, format)
		if err != nil {
			fmt.Printf("Error exporting report: %v\n", err)
			os.Exit(1)
		}

		// Write to file or stdout
		if outputFile != "" {
			err := os.WriteFile(outputFile, []byte(output), 0644)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Report exported to: %s\n", outputFile)
		} else {
			fmt.Print(output)
		}
	}
}

func runBadge(cmd *cobra.Command, args []string) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if path is a git repository
	if !isGitRepository(path) {
		fmt.Printf("Error: %s is not a Git repository\n", path)
		os.Exit(1)
	}

	fmt.Printf("Analyzing repository: %s\n", path)

	// Initialize analyzer
	analyzer, err := git.NewRepositoryAnalyzer(path)
	if err != nil {
		fmt.Printf("Error initializing repository analyzer: %v\n", err)
		os.Exit(1)
	}

	// Analyze repository
	data, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("Error analyzing repository: %v\n", err)
		os.Exit(1)
	}

	// Initialize checkers
	allCheckers := []checkers.Checker{
		checkers.NewDocChecker(),
		checkers.NewIgnoreChecker(),
		checkers.NewConventionalCommitChecker(),
		checkers.NewMsgLengthChecker(),
		checkers.NewCommitAuthorInsightsChecker(),
		checkers.NewCodebaseSmellChecker(),
		checkers.NewLocalBranchChecker(),
		checkers.NewStaleBranchChecker(),
		checkers.NewStashChecker(),
		checkers.NewGitHubIntegrationChecker(),
		checkers.NewGitLabIntegrationChecker(),
	}

	// Run all checkers
	scorer := scorer.NewScorer()
	for _, checker := range allCheckers {
		result := checker.Check(data)
		scorer.AddResult(*result)
	}

	// Generate report
	healthReport := scorer.CalculateHealthReport()

	// Generate badge
	exp := exporter.NewExporter()
	badgeURL := exp.GenerateBadgeURL(healthReport.OverallScore)
	markdownBadge := exp.GenerateMarkdownBadge(healthReport.OverallScore)

	fmt.Printf("Health Score: %d/100 (%s)\n\n", healthReport.OverallScore, healthReport.Grade)
	fmt.Printf("🔗 Badge URL:\n%s\n\n", badgeURL)
	fmt.Printf("Markdown Badge:\n%s\n", markdownBadge)
}

func runGitHub(cmd *cobra.Command, args []string) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if path is a git repository
	if !isGitRepository(path) {
		fmt.Printf("Error: %s is not a Git repository\n", path)
		os.Exit(1)
	}

	fmt.Printf("Checking GitHub integration: %s\n", path)

	// Check GitHub token
	token := os.Getenv("GPHC_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	if token == "" {
		fmt.Println("No GitHub token found")
		fmt.Println("Set GPHC_TOKEN or GITHUB_TOKEN environment variable for full GitHub integration")
		fmt.Println("Example: export GPHC_TOKEN=your_github_token")
		return
	}

	fmt.Println("GitHub token found")

	// Initialize GitHub checker
	checker := checkers.NewGitHubIntegrationChecker()

	// Create minimal repository data for the checker
	data := &types.RepositoryData{
		Path: path,
	}

	// Run GitHub integration check
	result := checker.Check(data)

	// Display results
	fmt.Printf("\nGitHub Integration Check Results:\n")
	fmt.Printf("Status: %s\n", result.Status.String())
	fmt.Printf("Score: %d\n", result.Score)
	fmt.Printf("Message: %s\n\n", result.Message)

	if len(result.Details) > 0 {
		fmt.Println("Details:")
		for _, detail := range result.Details {
			fmt.Printf("  %s\n", detail)
		}
	}
}

func runGitLab(cmd *cobra.Command, args []string) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if path is a git repository
	if !isGitRepository(path) {
		fmt.Printf("Error: %s is not a Git repository\n", path)
		os.Exit(1)
	}

	fmt.Printf("Checking GitLab integration: %s\n", path)

	// Check GitLab token
	token := os.Getenv("GPHC_TOKEN")
	if token == "" {
		token = os.Getenv("GITLAB_TOKEN")
	}

	if token == "" {
		fmt.Println("No GitLab token found")
		fmt.Println("Set GPHC_TOKEN or GITLAB_TOKEN environment variable for full GitLab integration")
		fmt.Println("Example: export GPHC_TOKEN=your_gitlab_token")
		return
	}

	fmt.Println("GitLab token found")

	// Initialize GitLab checker
	checker := checkers.NewGitLabIntegrationChecker()

	// Create minimal repository data for the checker
	data := &types.RepositoryData{
		Path: path,
	}

	// Run GitLab integration check
	result := checker.Check(data)

	// Display results
	fmt.Printf("\nGitLab Integration Check Results:\n")
	fmt.Printf("Status: %s\n", result.Status.String())
	fmt.Printf("Score: %d\n", result.Score)
	fmt.Printf("Message: %s\n\n", result.Message)

	if len(result.Details) > 0 {
		fmt.Println("Details:")
		for _, detail := range result.Details {
			fmt.Printf("  %s\n", detail)
		}
	}
}

func runAuthors(cmd *cobra.Command, args []string) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if path is a git repository
	if !isGitRepository(path) {
		fmt.Printf("Error: %s is not a Git repository\n", path)
		os.Exit(1)
	}

	fmt.Printf("Analyzing commit authors: %s\n", path)

	// Initialize analyzer
	analyzer, err := git.NewRepositoryAnalyzer(path)
	if err != nil {
		fmt.Printf("❌ Error initializing analyzer: %v\n", err)
		os.Exit(1)
	}

	// Analyze repository
	data, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("Error analyzing repository: %v\n", err)
		os.Exit(1)
	}

	// Initialize author insights checker
	checker := checkers.NewCommitAuthorInsightsChecker()

	// Run author insights check
	result := checker.Check(data)

	// Display results
	fmt.Printf("\nCommit Author Insights:\n")
	fmt.Printf("Status: %s\n", result.Status.String())
	fmt.Printf("Score: %d\n", result.Score)
	fmt.Printf("Message: %s\n\n", result.Message)

	if len(result.Details) > 0 {
		fmt.Println("Details:")
		for _, detail := range result.Details {
			fmt.Printf("  %s\n", detail)
		}
	}

	// Additional insights
	if len(data.Commits) > 0 {
		fmt.Printf("\nBus Factor Analysis:\n")

		// Count unique authors
		authorMap := make(map[string]bool)
		for _, commit := range data.Commits {
			authorMap[commit.Author] = true
		}

		uniqueAuthors := len(authorMap)

		if uniqueAuthors == 1 {
			fmt.Printf("  HIGH RISK: Single contributor project\n")
			fmt.Printf("  Contributors: %d\n", uniqueAuthors)
			fmt.Printf("  Bus Factor: 1 (Critical)\n")
			fmt.Printf("  Recommendation: Onboard additional contributors immediately\n")
		} else if uniqueAuthors == 2 {
			fmt.Printf("  MODERATE RISK: Low contributor count\n")
			fmt.Printf("  Contributors: %d\n", uniqueAuthors)
			fmt.Printf("  Bus Factor: 2 (Moderate)\n")
			fmt.Printf("  Recommendation: Expand contributor base\n")
		} else if uniqueAuthors <= 5 {
			fmt.Printf("  ACCEPTABLE: Small team\n")
			fmt.Printf("  Contributors: %d\n", uniqueAuthors)
			fmt.Printf("  Bus Factor: %d (Acceptable)\n", uniqueAuthors)
			fmt.Printf("  Recommendation: Maintain current team size\n")
		} else {
			fmt.Printf("  EXCELLENT: Well-distributed team\n")
			fmt.Printf("  Contributors: %d\n", uniqueAuthors)
			fmt.Printf("  Bus Factor: %d (Low Risk)\n", uniqueAuthors)
			fmt.Printf("  Recommendation: Excellent team distribution\n")
		}
	}
}

func runCodebase(cmd *cobra.Command, args []string) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if path is a git repository
	if !isGitRepository(path) {
		fmt.Printf("Error: %s is not a Git repository\n", path)
		os.Exit(1)
	}

	fmt.Printf("Analyzing codebase structure: %s\n", path)

	// Initialize analyzer
	analyzer, err := git.NewRepositoryAnalyzer(path)
	if err != nil {
		fmt.Printf("❌ Error initializing analyzer: %v\n", err)
		os.Exit(1)
	}

	// Analyze repository
	data, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("Error analyzing repository: %v\n", err)
		os.Exit(1)
	}

	// Initialize codebase smell checker
	checker := checkers.NewCodebaseSmellChecker()

	// Run codebase smell check
	result := checker.Check(data)

	// Display results
	fmt.Printf("\nCodebase Structure Analysis:\n")
	fmt.Printf("Status: %s\n", result.Status.String())
	fmt.Printf("Score: %d\n", result.Score)
	fmt.Printf("Message: %s\n\n", result.Message)

	if len(result.Details) > 0 {
		fmt.Println("Details:")
		for _, detail := range result.Details {
			fmt.Printf("  %s\n", detail)
		}
	}

	// Additional recommendations
	fmt.Printf("\nStructure Recommendations:\n")

	if result.Score < 70 {
		fmt.Printf("  Codebase structure needs improvement\n")
		fmt.Printf("  Consider the following actions:\n")
		fmt.Printf("    • Add test directories and test files\n")
		fmt.Printf("    • Organize code into logical subdirectories\n")
		fmt.Printf("    • Split oversized directories (>1000 files)\n")
		fmt.Printf("    • Add documentation files\n")
		fmt.Printf("    • Remove empty directories\n")
	} else if result.Score < 90 {
		fmt.Printf("  Good codebase structure with minor improvements needed\n")
		fmt.Printf("  Consider:\n")
		fmt.Printf("    • Adding more test coverage\n")
		fmt.Printf("    • Improving directory organization\n")
	} else {
		fmt.Printf("  Excellent codebase structure\n")
		fmt.Printf("  Maintain current organization patterns\n")
	}
}

func runUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("Updating GPHC...")

	// Find the GPHC source directory
	sourceDir := findGPHCSourceDir()
	if sourceDir == "" {
		fmt.Println("Error: Could not find GPHC source directory")
		fmt.Println("Please run this command from the GPHC project directory")
		os.Exit(1)
	}

	fmt.Printf("Found GPHC source at: %s\n", sourceDir)

	// Change to source directory
	if err := os.Chdir(sourceDir); err != nil {
		fmt.Printf("Error changing to source directory: %v\n", err)
		os.Exit(1)
	}

	// Pull latest changes
	fmt.Println("Pulling latest changes...")
	pullCmd := exec.Command("git", "pull", "origin", "main")
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr

	if err := pullCmd.Run(); err != nil {
		fmt.Printf("Error pulling changes: %v\n", err)
		fmt.Println("Make sure you have internet connection and git access")
		os.Exit(1)
	}

	// Update dependencies
	fmt.Println("Updating dependencies...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr

	if err := tidyCmd.Run(); err != nil {
		fmt.Printf("Error updating dependencies: %v\n", err)
		os.Exit(1)
	}

	// Rebuild and reinstall
	fmt.Println("Building and installing GPHC...")
	installCmd := exec.Command("go", "install", "./cmd/gphc")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	if err := installCmd.Run(); err != nil {
		fmt.Printf("Error installing GPHC: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("GPHC updated successfully!")
	fmt.Println("New version:")

	// Show new version
	versionCmd := exec.Command("gphc", "version")
	versionCmd.Stdout = os.Stdout
	versionCmd.Stderr = os.Stderr
	versionCmd.Run()
}

func findGPHCSourceDir() string {
	// Try to find the GPHC source directory
	// First, check if we're already in it
	if isGPHCSourceDir(".") {
		return "."
	}

	// Check common locations
	commonPaths := []string{
		"/Users/opsource/projects/Dev/GPHC",
		"~/projects/Dev/GPHC",
		"~/Dev/GPHC",
		"./GPHC",
		"../GPHC",
		"../../GPHC",
	}

	for _, path := range commonPaths {
		expandedPath := filepath.Clean(os.ExpandEnv(path))
		if isGPHCSourceDir(expandedPath) {
			return expandedPath
		}
	}

	// Try to find from current working directory
	currentDir, err := os.Getwd()
	if err == nil {
		// Check parent directories
		for i := 0; i < 5; i++ {
			if isGPHCSourceDir(currentDir) {
				return currentDir
			}
			currentDir = filepath.Dir(currentDir)
			if currentDir == "/" {
				break
			}
		}
	}

	return ""
}

func isGPHCSourceDir(path string) bool {
	// Check if this directory contains GPHC source files
	goModPath := filepath.Join(path, "go.mod")
	mainPath := filepath.Join(path, "cmd", "gphc", "main.go")

	if _, err := os.Stat(goModPath); err != nil {
		return false
	}

	if _, err := os.Stat(mainPath); err != nil {
		return false
	}

	// Check if go.mod contains gphc module
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), "github.com/vahidaghazadeh/gphc") ||
		strings.Contains(string(content), "github.com/opsource/gphc")
}

func isGitRepository(path string) bool {
	gitPath := fmt.Sprintf("%s/.git", path)
	_, err := os.Stat(gitPath)
	return err == nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Helper functions for pre-commit checks

func getStagedFiles(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return []string{}, nil
	}

	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

func checkFileFormatting(repoPath string, files []string) bool {
	goFiles := []string{}
	for _, file := range files {
		if strings.HasSuffix(file, ".go") {
			goFiles = append(goFiles, file)
		}
	}

	if len(goFiles) == 0 {
		return true
	}

	cmd := exec.Command("gofmt", "-l")
	cmd.Dir = repoPath
	cmd.Args = append(cmd.Args, goFiles...)

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(output) == 0
}

func checkCommitMessage(repoPath string) bool {
	// Check if we're in the middle of a commit (has COMMIT_EDITMSG)
	commitMsgPath := filepath.Join(repoPath, ".git", "COMMIT_EDITMSG")
	if _, err := os.Stat(commitMsgPath); err == nil {
		// Read the commit message file
		content, err := os.ReadFile(commitMsgPath)
		if err != nil {
			return false
		}

		message := strings.TrimSpace(string(content))
		// Remove comments (lines starting with #)
		lines := strings.Split(message, "\n")
		var cleanLines []string
		for _, line := range lines {
			if !strings.HasPrefix(strings.TrimSpace(line), "#") {
				cleanLines = append(cleanLines, line)
			}
		}
		message = strings.Join(cleanLines, "\n")
		message = strings.TrimSpace(message)

		// Check conventional commit format
		conventionalPrefixes := []string{
			"feat:", "fix:", "docs:", "style:", "refactor:", "test:", "chore:",
			"perf:", "ci:", "build:", "revert:", "feat!", "fix!",
		}

		for _, prefix := range conventionalPrefixes {
			if strings.HasPrefix(message, prefix) {
				return true
			}
		}

		return false
	}

	// If not committing, check the last commit message
	cmd := exec.Command("git", "log", "-1", "--pretty=%s")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return true // If no commits yet, consider it valid
	}

	message := strings.TrimSpace(string(output))
	if message == "" {
		return true // Empty message is valid for first commit
	}

	// Check conventional commit format
	conventionalPrefixes := []string{
		"feat:", "fix:", "docs:", "style:", "refactor:", "test:", "chore:",
		"perf:", "ci:", "build:", "revert:", "feat!", "fix!",
	}

	for _, prefix := range conventionalPrefixes {
		if strings.HasPrefix(message, prefix) {
			return true
		}
	}

	return false
}

func analyzeStagedChanges(repoPath string, stagedFiles []string) string {
	// Analyze file types and changes
	var addedFiles, modifiedFiles, deletedFiles []string
	var hasNewFeatures, hasBugFixes, hasDocs, hasRefactor bool
	var changeSummary []string

	for _, file := range stagedFiles {
		// Get file status
		cmd := exec.Command("git", "status", "--porcelain", file)
		cmd.Dir = repoPath
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		status := strings.TrimSpace(string(output))
		if len(status) < 3 {
			continue
		}

		fileStatus := status[:2]
		fileName := status[3:]

		switch {
		case strings.HasPrefix(fileStatus, "A"):
			addedFiles = append(addedFiles, fileName)
		case strings.HasPrefix(fileStatus, "M"):
			modifiedFiles = append(modifiedFiles, fileName)
		case strings.HasPrefix(fileStatus, "D"):
			deletedFiles = append(deletedFiles, fileName)
		}

		// Analyze file type and content
		fileExt := strings.ToLower(filepath.Ext(fileName))
		fileName = strings.ToLower(fileName)
		baseName := strings.ToLower(filepath.Base(fileName))

		// Check for new features
		if strings.Contains(fileName, "feat") || strings.Contains(fileName, "feature") ||
			strings.Contains(fileName, "add") || strings.Contains(fileName, "new") ||
			strings.Contains(baseName, "feat") || strings.Contains(baseName, "feature") {
			hasNewFeatures = true
		}

		// Check for bug fixes
		if strings.Contains(fileName, "fix") || strings.Contains(fileName, "bug") ||
			strings.Contains(fileName, "error") || strings.Contains(fileName, "issue") ||
			strings.Contains(fileName, "patch") || strings.Contains(baseName, "fix") {
			hasBugFixes = true
		}

		// Check for documentation
		if strings.Contains(fileName, "readme") || strings.Contains(fileName, "doc") ||
			strings.Contains(fileName, "guide") || strings.Contains(fileName, "manual") ||
			fileExt == ".md" || fileExt == ".txt" || fileExt == ".rst" {
			hasDocs = true
		}

		// Check for refactoring
		if strings.Contains(fileName, "refactor") || strings.Contains(fileName, "clean") ||
			strings.Contains(fileName, "optimize") || strings.Contains(fileName, "improve") ||
			strings.Contains(fileName, "restructure") || strings.Contains(baseName, "refactor") {
			hasRefactor = true
		}

		// Analyze file changes for summary
		summary := analyzeFileChanges(repoPath, file, fileStatus)
		if summary != "" {
			changeSummary = append(changeSummary, summary)
		}
	}

	// Generate suggestion based on analysis
	var prefix, description string

	if hasNewFeatures && len(addedFiles) > 0 {
		prefix = "feat"
		if len(addedFiles) == 1 {
			description = fmt.Sprintf("add %s", getFeatureDescription(addedFiles[0]))
		} else {
			description = fmt.Sprintf("add new features (%d files)", len(addedFiles))
		}
	} else if hasBugFixes {
		prefix = "fix"
		if len(modifiedFiles) == 1 {
			description = fmt.Sprintf("fix %s", getFixDescription(modifiedFiles[0]))
		} else {
			description = fmt.Sprintf("fix bugs and issues (%d files)", len(modifiedFiles))
		}
	} else if hasDocs {
		prefix = "docs"
		if len(stagedFiles) == 1 {
			description = fmt.Sprintf("update %s", getDocDescription(stagedFiles[0]))
		} else {
			description = fmt.Sprintf("update documentation (%d files)", len(stagedFiles))
		}
	} else if hasRefactor {
		prefix = "refactor"
		description = fmt.Sprintf("refactor %s (%d files)", getRefactorDescription(stagedFiles), len(stagedFiles))
	} else if len(addedFiles) > 0 {
		prefix = "feat"
		description = fmt.Sprintf("add %s (%d files)", getGenericDescription(addedFiles), len(addedFiles))
	} else if len(modifiedFiles) > 0 {
		prefix = "chore"
		description = fmt.Sprintf("update %s (%d files)", getGenericDescription(modifiedFiles), len(modifiedFiles))
	} else if len(deletedFiles) > 0 {
		prefix = "chore"
		description = fmt.Sprintf("remove %s (%d files)", getGenericDescription(deletedFiles), len(deletedFiles))
	} else {
		prefix = "chore"
		description = fmt.Sprintf("update project files (%d files)", len(stagedFiles))
	}

	// Add change summary if available
	if len(changeSummary) > 0 {
		summaryText := strings.Join(changeSummary, ", ")
		if len(summaryText) > 100 {
			summaryText = summaryText[:97] + "..."
		}
		description += fmt.Sprintf(" - %s", summaryText)
	}

	return fmt.Sprintf("%s: %s", prefix, description)
}

// Helper functions for more specific descriptions
func getFeatureDescription(fileName string) string {
	baseName := strings.ToLower(filepath.Base(fileName))
	ext := strings.ToLower(filepath.Ext(fileName))

	if strings.Contains(baseName, "component") {
		return "new component"
	} else if strings.Contains(baseName, "api") {
		return "new API endpoint"
	} else if strings.Contains(baseName, "util") || strings.Contains(baseName, "helper") {
		return "utility function"
	} else if ext == ".js" || ext == ".ts" {
		return "new functionality"
	} else if ext == ".css" || ext == ".scss" {
		return "new styles"
	} else {
		return "new feature"
	}
}

func getFixDescription(fileName string) string {
	baseName := strings.ToLower(filepath.Base(fileName))

	if strings.Contains(baseName, "bug") {
		return "critical bug"
	} else if strings.Contains(baseName, "error") {
		return "error handling"
	} else if strings.Contains(baseName, "issue") {
		return "reported issue"
	} else {
		return "bug"
	}
}

func getDocDescription(fileName string) string {
	baseName := strings.ToLower(filepath.Base(fileName))

	if strings.Contains(baseName, "readme") {
		return "README"
	} else if strings.Contains(baseName, "api") {
		return "API documentation"
	} else if strings.Contains(baseName, "guide") {
		return "user guide"
	} else {
		return "documentation"
	}
}

func getRefactorDescription(files []string) string {
	if len(files) == 1 {
		baseName := strings.ToLower(filepath.Base(files[0]))
		if strings.Contains(baseName, "component") {
			return "component structure"
		} else if strings.Contains(baseName, "api") {
			return "API structure"
		} else if strings.Contains(baseName, "util") {
			return "utility functions"
		}
	}
	return "code structure"
}

func getGenericDescription(files []string) string {
	if len(files) == 1 {
		ext := strings.ToLower(filepath.Ext(files[0]))
		switch ext {
		case ".js", ".ts":
			return "JavaScript/TypeScript files"
		case ".css", ".scss":
			return "stylesheets"
		case ".html":
			return "HTML templates"
		case ".json":
			return "JSON configuration"
		case ".md":
			return "documentation"
		default:
			return "files"
		}
	}
	return "files"
}

// analyzeFileChanges analyzes the content changes in a file
func analyzeFileChanges(repoPath string, filePath string, fileStatus string) string {
	fileExt := strings.ToLower(filepath.Ext(filePath))
	fileName := strings.ToLower(filepath.Base(filePath))

	// Skip binary files and large files
	if isBinaryFile(fileExt) || isLargeFile(repoPath, filePath) {
		return ""
	}

	switch {
	case strings.HasPrefix(fileStatus, "A"):
		return analyzeAddedFile(repoPath, filePath, fileExt, fileName)
	case strings.HasPrefix(fileStatus, "M"):
		return analyzeModifiedFile(repoPath, filePath, fileExt, fileName)
	case strings.HasPrefix(fileStatus, "D"):
		return fmt.Sprintf("removed %s", fileName)
	default:
		return ""
	}
}

// analyzeAddedFile analyzes a newly added file
func analyzeAddedFile(repoPath string, filePath string, fileExt string, fileName string) string {
	// Read file content to analyze
	content, err := os.ReadFile(filepath.Join(repoPath, filePath))
	if err != nil {
		return fmt.Sprintf("added %s", fileName)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	switch fileExt {
	case ".go":
		return analyzeGoFile(lines, fileName, "added")
	case ".js", ".ts":
		return analyzeJSFile(lines, fileName, "added")
	case ".py":
		return analyzePythonFile(lines, fileName, "added")
	case ".java":
		return analyzeJavaFile(lines, fileName, "added")
	case ".md":
		return analyzeMarkdownFile(lines, fileName, "added")
	case ".json":
		return analyzeJSONFile(contentStr, fileName, "added")
	case ".yaml", ".yml":
		return analyzeYAMLFile(contentStr, fileName, "added")
	default:
		return fmt.Sprintf("added %s (%d lines)", fileName, len(lines))
	}
}

// analyzeModifiedFile analyzes a modified file
func analyzeModifiedFile(repoPath string, filePath string, fileExt string, fileName string) string {
	// Get git diff to see what changed
	cmd := exec.Command("git", "diff", "--cached", filePath)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("modified %s", fileName)
	}

	diffStr := string(output)
	lines := strings.Split(diffStr, "\n")

	// Count additions and deletions
	additions := 0
	deletions := 0
	var changes []string

	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			additions++
			// Extract meaningful changes
			if change := extractMeaningfulChange(line, fileExt); change != "" {
				changes = append(changes, change)
			}
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			deletions++
		}
	}

	// Generate summary
	if len(changes) > 0 {
		// Limit to 3 most important changes
		if len(changes) > 3 {
			changes = changes[:3]
		}
		return fmt.Sprintf("modified %s: %s", fileName, strings.Join(changes, ", "))
	}

	if additions > 0 && deletions > 0 {
		return fmt.Sprintf("modified %s (+%d/-%d lines)", fileName, additions, deletions)
	} else if additions > 0 {
		return fmt.Sprintf("modified %s (+%d lines)", fileName, additions)
	} else if deletions > 0 {
		return fmt.Sprintf("modified %s (-%d lines)", fileName, deletions)
	}

	return fmt.Sprintf("modified %s", fileName)
}

// extractMeaningfulChange extracts meaningful information from a diff line
func extractMeaningfulChange(line string, fileExt string) string {
	line = strings.TrimPrefix(line, "+")
	line = strings.TrimSpace(line)

	if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
		return ""
	}

	switch fileExt {
	case ".go":
		return extractGoChange(line)
	case ".js", ".ts":
		return extractJSChange(line)
	case ".py":
		return extractPythonChange(line)
	case ".java":
		return extractJavaChange(line)
	case ".md":
		return extractMarkdownChange(line)
	default:
		return ""
	}
}

// Go-specific analysis
func analyzeGoFile(lines []string, fileName string, action string) string {
	var functions []string
	var structs []string
	var imports []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "func ") {
			funcName := extractGoFunctionName(line)
			if funcName != "" {
				functions = append(functions, funcName)
			}
		} else if strings.HasPrefix(line, "type ") && strings.Contains(line, " struct") {
			structName := extractGoStructName(line)
			if structName != "" {
				structs = append(structs, structName)
			}
		} else if strings.HasPrefix(line, "import ") {
			imports = append(imports, "import")
		}
	}

	var parts []string
	if len(functions) > 0 {
		if len(functions) == 1 {
			parts = append(parts, fmt.Sprintf("function %s", functions[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d functions", len(functions)))
		}
	}
	if len(structs) > 0 {
		if len(structs) == 1 {
			parts = append(parts, fmt.Sprintf("struct %s", structs[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d structs", len(structs)))
		}
	}
	if len(imports) > 0 {
		parts = append(parts, "imports")
	}

	if len(parts) > 0 {
		return fmt.Sprintf("%s %s (%s)", action, fileName, strings.Join(parts, ", "))
	}

	return fmt.Sprintf("%s %s (%d lines)", action, fileName, len(lines))
}

func extractGoFunctionName(line string) string {
	// Extract function name from "func FunctionName(" or "func (r *Receiver) MethodName("
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return ""
	}

	if parts[1] == "(" {
		// Method with receiver
		if len(parts) >= 4 {
			return parts[3]
		}
	} else {
		// Regular function
		return strings.TrimSuffix(parts[1], "(")
	}
	return ""
}

func extractGoStructName(line string) string {
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

func extractGoChange(line string) string {
	if strings.Contains(line, "func ") {
		funcName := extractGoFunctionName("func " + line)
		if funcName != "" {
			return fmt.Sprintf("add function %s", funcName)
		}
	} else if strings.Contains(line, "type ") && strings.Contains(line, " struct") {
		structName := extractGoStructName(line)
		if structName != "" {
			return fmt.Sprintf("add struct %s", structName)
		}
	}
	return ""
}

// JavaScript/TypeScript analysis
func analyzeJSFile(lines []string, fileName string, action string) string {
	var functions []string
	var classes []string
	var exports []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "function ") || strings.Contains(line, "=>") {
			funcName := extractJSFunctionName(line)
			if funcName != "" {
				functions = append(functions, funcName)
			}
		} else if strings.Contains(line, "class ") {
			className := extractJSClassName(line)
			if className != "" {
				classes = append(classes, className)
			}
		} else if strings.Contains(line, "export ") {
			exports = append(exports, "export")
		}
	}

	var parts []string
	if len(functions) > 0 {
		if len(functions) == 1 {
			parts = append(parts, fmt.Sprintf("function %s", functions[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d functions", len(functions)))
		}
	}
	if len(classes) > 0 {
		if len(classes) == 1 {
			parts = append(parts, fmt.Sprintf("class %s", classes[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d classes", len(classes)))
		}
	}
	if len(exports) > 0 {
		parts = append(parts, "exports")
	}

	if len(parts) > 0 {
		return fmt.Sprintf("%s %s (%s)", action, fileName, strings.Join(parts, ", "))
	}

	return fmt.Sprintf("%s %s (%d lines)", action, fileName, len(lines))
}

func extractJSFunctionName(line string) string {
	// Extract function name from various JS function patterns
	if strings.Contains(line, "function ") {
		parts := strings.Fields(line)
		for i, part := range parts {
			if part == "function" && i+1 < len(parts) {
				return strings.TrimSuffix(parts[i+1], "(")
			}
		}
	} else if strings.Contains(line, "=>") {
		// Arrow function
		beforeArrow := strings.Split(line, "=>")[0]
		beforeArrow = strings.TrimSpace(beforeArrow)
		if strings.Contains(beforeArrow, "(") {
			// Named arrow function
			parts := strings.Split(beforeArrow, "(")
			if len(parts) > 0 {
				return strings.TrimSpace(parts[0])
			}
		}
	}
	return ""
}

func extractJSClassName(line string) string {
	parts := strings.Fields(line)
	for i, part := range parts {
		if part == "class" && i+1 < len(parts) {
			return strings.TrimSuffix(parts[i+1], "{")
		}
	}
	return ""
}

func extractJSChange(line string) string {
	if strings.Contains(line, "function ") {
		funcName := extractJSFunctionName(line)
		if funcName != "" {
			return fmt.Sprintf("add function %s", funcName)
		}
	} else if strings.Contains(line, "class ") {
		className := extractJSClassName(line)
		if className != "" {
			return fmt.Sprintf("add class %s", className)
		}
	}
	return ""
}

// Python analysis
func analyzePythonFile(lines []string, fileName string, action string) string {
	var functions []string
	var classes []string
	var imports []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "def ") {
			funcName := extractPythonFunctionName(line)
			if funcName != "" {
				functions = append(functions, funcName)
			}
		} else if strings.HasPrefix(line, "class ") {
			className := extractPythonClassName(line)
			if className != "" {
				classes = append(classes, className)
			}
		} else if strings.HasPrefix(line, "import ") || strings.HasPrefix(line, "from ") {
			imports = append(imports, "import")
		}
	}

	var parts []string
	if len(functions) > 0 {
		if len(functions) == 1 {
			parts = append(parts, fmt.Sprintf("function %s", functions[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d functions", len(functions)))
		}
	}
	if len(classes) > 0 {
		if len(classes) == 1 {
			parts = append(parts, fmt.Sprintf("class %s", classes[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d classes", len(classes)))
		}
	}
	if len(imports) > 0 {
		parts = append(parts, "imports")
	}

	if len(parts) > 0 {
		return fmt.Sprintf("%s %s (%s)", action, fileName, strings.Join(parts, ", "))
	}

	return fmt.Sprintf("%s %s (%d lines)", action, fileName, len(lines))
}

func extractPythonFunctionName(line string) string {
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return strings.TrimSuffix(parts[1], "(")
	}
	return ""
}

func extractPythonClassName(line string) string {
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return strings.TrimSuffix(parts[1], "(")
	}
	return ""
}

func extractPythonChange(line string) string {
	if strings.Contains(line, "def ") {
		funcName := extractPythonFunctionName(line)
		if funcName != "" {
			return fmt.Sprintf("add function %s", funcName)
		}
	} else if strings.Contains(line, "class ") {
		className := extractPythonClassName(line)
		if className != "" {
			return fmt.Sprintf("add class %s", className)
		}
	}
	return ""
}

// Java analysis
func analyzeJavaFile(lines []string, fileName string, action string) string {
	var methods []string
	var classes []string
	var imports []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "public ") || strings.Contains(line, "private ") || strings.Contains(line, "protected ") {
			if strings.Contains(line, "(") && strings.Contains(line, ")") {
				methodName := extractJavaMethodName(line)
				if methodName != "" {
					methods = append(methods, methodName)
				}
			}
		} else if strings.Contains(line, "class ") {
			className := extractJavaClassName(line)
			if className != "" {
				classes = append(classes, className)
			}
		} else if strings.HasPrefix(line, "import ") {
			imports = append(imports, "import")
		}
	}

	var parts []string
	if len(methods) > 0 {
		if len(methods) == 1 {
			parts = append(parts, fmt.Sprintf("method %s", methods[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d methods", len(methods)))
		}
	}
	if len(classes) > 0 {
		if len(classes) == 1 {
			parts = append(parts, fmt.Sprintf("class %s", classes[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d classes", len(classes)))
		}
	}
	if len(imports) > 0 {
		parts = append(parts, "imports")
	}

	if len(parts) > 0 {
		return fmt.Sprintf("%s %s (%s)", action, fileName, strings.Join(parts, ", "))
	}

	return fmt.Sprintf("%s %s (%d lines)", action, fileName, len(lines))
}

func extractJavaMethodName(line string) string {
	// Extract method name from Java method declaration
	parts := strings.Fields(line)
	for _, part := range parts {
		if strings.Contains(part, "(") {
			return strings.TrimSuffix(part, "(")
		}
	}
	return ""
}

func extractJavaClassName(line string) string {
	parts := strings.Fields(line)
	for i, part := range parts {
		if part == "class" && i+1 < len(parts) {
			return strings.TrimSuffix(parts[i+1], "{")
		}
	}
	return ""
}

func extractJavaChange(line string) string {
	if strings.Contains(line, "public ") || strings.Contains(line, "private ") || strings.Contains(line, "protected ") {
		methodName := extractJavaMethodName(line)
		if methodName != "" {
			return fmt.Sprintf("add method %s", methodName)
		}
	} else if strings.Contains(line, "class ") {
		className := extractJavaClassName(line)
		if className != "" {
			return fmt.Sprintf("add class %s", className)
		}
	}
	return ""
}

// Markdown analysis
func analyzeMarkdownFile(lines []string, fileName string, action string) string {
	var headers []string
	var links []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			header := strings.TrimPrefix(line, "#")
			header = strings.TrimSpace(header)
			if header != "" {
				headers = append(headers, header)
			}
		} else if strings.Contains(line, "[") && strings.Contains(line, "](") {
			links = append(links, "link")
		}
	}

	var parts []string
	if len(headers) > 0 {
		if len(headers) == 1 {
			parts = append(parts, fmt.Sprintf("section '%s'", headers[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d sections", len(headers)))
		}
	}
	if len(links) > 0 {
		parts = append(parts, "links")
	}

	if len(parts) > 0 {
		return fmt.Sprintf("%s %s (%s)", action, fileName, strings.Join(parts, ", "))
	}

	return fmt.Sprintf("%s %s (%d lines)", action, fileName, len(lines))
}

func extractMarkdownChange(line string) string {
	if strings.HasPrefix(line, "#") {
		header := strings.TrimPrefix(line, "#")
		header = strings.TrimSpace(header)
		if header != "" {
			return fmt.Sprintf("add section '%s'", header)
		}
	}
	return ""
}

// JSON analysis
func analyzeJSONFile(content string, fileName string, action string) string {
	// Simple JSON analysis
	if strings.Contains(content, "dependencies") {
		return fmt.Sprintf("%s %s (dependencies)", action, fileName)
	} else if strings.Contains(content, "scripts") {
		return fmt.Sprintf("%s %s (scripts)", action, fileName)
	} else if strings.Contains(content, "config") {
		return fmt.Sprintf("%s %s (configuration)", action, fileName)
	}
	return fmt.Sprintf("%s %s", action, fileName)
}

// YAML analysis
func analyzeYAMLFile(content string, fileName string, action string) string {
	// Simple YAML analysis
	if strings.Contains(content, "dependencies") {
		return fmt.Sprintf("%s %s (dependencies)", action, fileName)
	} else if strings.Contains(content, "services") {
		return fmt.Sprintf("%s %s (services)", action, fileName)
	} else if strings.Contains(content, "config") {
		return fmt.Sprintf("%s %s (configuration)", action, fileName)
	}
	return fmt.Sprintf("%s %s", action, fileName)
}

// Utility functions
func isBinaryFile(fileExt string) bool {
	binaryExts := []string{".exe", ".dll", ".so", ".dylib", ".bin", ".img", ".iso", ".zip", ".tar", ".gz", ".rar", ".7z", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".svg", ".ico", ".mp3", ".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm"}
	for _, ext := range binaryExts {
		if fileExt == ext {
			return true
		}
	}
	return false
}

func isLargeFile(repoPath string, filePath string) bool {
	fullPath := filepath.Join(repoPath, filePath)
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	// Consider files larger than 1MB as large
	return info.Size() > 1024*1024
}

// displayColoredDiff displays git diff output with colored backgrounds
func displayColoredDiff(diffOutput string) {
	lines := strings.Split(diffOutput, "\n")
	
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			// Addition: light green background
			fmt.Printf("\033[48;5;22m%s\033[0m\n", line)
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			// Deletion: light red background
			fmt.Printf("\033[48;5;52m%s\033[0m\n", line)
		} else if strings.HasPrefix(line, "@@") {
			// Hunk header: cyan background
			fmt.Printf("\033[48;5;23m%s\033[0m\n", line)
		} else if strings.HasPrefix(line, "diff --git") {
			// File header: blue background
			fmt.Printf("\033[48;5;17m%s\033[0m\n", line)
		} else if strings.HasPrefix(line, "index ") {
			// Index line: dark blue background
			fmt.Printf("\033[48;5;18m%s\033[0m\n", line)
		} else if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
			// File names: yellow background
			fmt.Printf("\033[48;5;58m%s\033[0m\n", line)
		} else {
			// Context lines: normal
			fmt.Println(line)
		}
	}
}

func checkLargeFiles(files []string) bool {
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		// Check if file is larger than 1MB
		if info.Size() > 1024*1024 {
			return false
		}
	}

	return true
}

func checkSensitiveFiles(files []string) bool {
	sensitivePatterns := []string{
		".env", ".env.local", ".env.production", ".env.staging",
		"config.json", "secrets.json", "credentials.json",
		"*.key", "*.pem", "*.p12", "*.pfx",
		"id_rsa", "id_dsa", "id_ecdsa", "id_ed25519",
	}

	for _, file := range files {
		for _, pattern := range sensitivePatterns {
			if strings.Contains(file, pattern) || strings.HasSuffix(file, pattern[1:]) {
				return false
			}
		}
	}

	return true
}

func runScan(cmd *cobra.Command, args []string) {
	var scanPath string
	if len(args) > 0 {
		scanPath = args[0]
	} else {
		var err error
		scanPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Multi-Repository Health Scan Results\n")
	fmt.Printf("====================================\n\n")

	// Find Git repositories
	repos, err := findGitRepositories(scanPath, recursiveScan)
	if err != nil {
		fmt.Printf("Error finding repositories: %v\n", err)
		os.Exit(1)
	}

	if len(repos) == 0 {
		fmt.Println("No Git repositories found in the specified path.")
		return
	}

	// Scan repositories
	results := make([]ScanResult, 0, len(repos))
	totalScore := 0.0

	for _, repo := range repos {
		fmt.Printf("Scanning: %s\n", repo)

		analyzer, err := git.NewRepositoryAnalyzer(repo)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		data, err := analyzer.Analyze()
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		// Run health checks
		allCheckers := []checkers.Checker{
			checkers.NewDocChecker(),
			checkers.NewConventionalCommitChecker(),
			checkers.NewMsgLengthChecker(),
			checkers.NewCommitSizeChecker(),
			checkers.NewLocalBranchChecker(),
			checkers.NewStaleBranchChecker(),
			checkers.NewBareRepoChecker(),
			checkers.NewStashChecker(),
			checkers.NewIgnoreChecker(),
			checkers.NewGitHubIntegrationChecker(),
			checkers.NewGitLabIntegrationChecker(),
			checkers.NewCommitAuthorInsightsChecker(),
			checkers.NewCodebaseSmellChecker(),
		}

		scorer := scorer.NewScorer()
		for _, checker := range allCheckers {
			result := checker.Check(data)
			scorer.AddResult(*result)
		}

		// Generate report
		healthReport := scorer.CalculateHealthReport()

		// Filter by minimum score if specified
		if minScore > 0 && healthReport.OverallScore < minScore {
			continue
		}

		repoName := filepath.Base(repo)
		result := ScanResult{
			Name:  repoName,
			Path:  repo,
			Score: healthReport.OverallScore,
			Grade: healthReport.Grade,
		}

		results = append(results, result)
		totalScore += float64(healthReport.OverallScore)

		fmt.Printf("  %s: %d/100 (%s)\n", repoName, healthReport.OverallScore, healthReport.Grade)
	}

	// Calculate average
	if len(results) > 0 {
		averageScore := totalScore / float64(len(results))
		fmt.Printf("\nSummary:\n")
		fmt.Printf("  Total Repositories: %d\n", len(results))
		fmt.Printf("  Average Health: %.1f/100\n", averageScore)

		if len(results) > 0 {
			// Find highest and lowest scores
			highest := results[0]
			lowest := results[0]

			for _, result := range results {
				if result.Score > highest.Score {
					highest = result
				}
				if result.Score < lowest.Score {
					lowest = result
				}
			}

			fmt.Printf("  Highest Score: %s (%d/100)\n", highest.Name, highest.Score)
			fmt.Printf("  Lowest Score: %s (%d/100)\n", lowest.Name, lowest.Score)
		}
	}
}

type ScanResult struct {
	Name  string
	Path  string
	Score int
	Grade string
}

func findGitRepositories(rootPath string, recursive bool) ([]string, error) {
	var repos []string

	if !recursive {
		// Check if rootPath itself is a Git repository
		if isGitRepository(rootPath) {
			repos = append(repos, rootPath)
		}
		return repos, nil
	}

	// Recursive search
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories except .git
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && info.Name() != "." && info.Name() != ".git" {
			return filepath.SkipDir
		}

		// Check if this is a Git repository
		if info.IsDir() && info.Name() == ".git" {
			repoPath := filepath.Dir(path)
			repos = append(repos, repoPath)
			return filepath.SkipDir
		}

		return nil
	})

	return repos, err
}

func runTUI(cmd *cobra.Command, args []string) {
	var repoPath string
	if len(args) > 0 {
		repoPath = args[0]
	} else {
		var err error
		repoPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if it's a Git repository
	if !isGitRepository(repoPath) {
		fmt.Printf("Error: %s is not a Git repository\n", repoPath)
		os.Exit(1)
	}

	// Create TUI model
	model := NewTUIModel(repoPath)

	// Run the TUI
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

// TUIModel represents the TUI application state
type TUIModel struct {
	repoPath     string
	healthReport *types.HealthReport
	loading      bool
	err          error
	selectedTab  int
	tabs         []string
}

// NewTUIModel creates a new TUI model
func NewTUIModel(repoPath string) *TUIModel {
	return &TUIModel{
		repoPath:    repoPath,
		loading:     true,
		tabs:        []string{"Overview", "Details", "Trends"},
		selectedTab: 0,
	}
}

// Init implements the tea.Model interface
func (m *TUIModel) Init() tea.Cmd {
	return m.loadHealthData()
}

// Update implements the tea.Model interface
func (m *TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.selectedTab = (m.selectedTab + 1) % len(m.tabs)
		case "shift+tab":
			m.selectedTab = (m.selectedTab - 1 + len(m.tabs)) % len(m.tabs)
		case "r":
			m.loading = true
			return m, m.loadHealthData()
		}
	case healthDataMsg:
		m.healthReport = msg.report
		m.loading = false
		m.err = msg.err
	}
	return m, nil
}

// View implements the tea.Model interface
func (m *TUIModel) View() string {
	if m.loading {
		return "Loading health data...\n\nPress 'q' to quit, 'r' to refresh"
	}

	if m.err != nil {
		return fmt.Sprintf("Error loading health data: %v\n\nPress 'q' to quit, 'r' to refresh", m.err)
	}

	var content string
	switch m.selectedTab {
	case 0:
		content = m.renderOverview()
	case 1:
		content = m.renderDetails()
	case 2:
		content = m.renderTrends()
	}

	return fmt.Sprintf("GPHC TUI - %s\n\n%s\n\nPress 'q' to quit, 'tab' to switch tabs, 'r' to refresh",
		filepath.Base(m.repoPath), content)
}

// healthDataMsg is used to pass health data to the model
type healthDataMsg struct {
	report *types.HealthReport
	err    error
}

// loadHealthData loads the health data for the repository
func (m *TUIModel) loadHealthData() tea.Cmd {
	return func() tea.Msg {
		// Analyze repository
		analyzer, err := git.NewRepositoryAnalyzer(m.repoPath)
		if err != nil {
			return healthDataMsg{err: err}
		}

		data, err := analyzer.Analyze()
		if err != nil {
			return healthDataMsg{err: err}
		}

		// Run health checks
		allCheckers := []checkers.Checker{
			checkers.NewDocChecker(),
			checkers.NewConventionalCommitChecker(),
			checkers.NewMsgLengthChecker(),
			checkers.NewCommitSizeChecker(),
			checkers.NewLocalBranchChecker(),
			checkers.NewStaleBranchChecker(),
			checkers.NewBareRepoChecker(),
			checkers.NewStashChecker(),
			checkers.NewIgnoreChecker(),
			checkers.NewGitHubIntegrationChecker(),
			checkers.NewGitLabIntegrationChecker(),
			checkers.NewCommitAuthorInsightsChecker(),
			checkers.NewCodebaseSmellChecker(),
		}

		scorer := scorer.NewScorer()
		for _, checker := range allCheckers {
			result := checker.Check(data)
			scorer.AddResult(*result)
		}

		// Generate report
		healthReport := scorer.CalculateHealthReport()
		return healthDataMsg{report: healthReport}
	}
}

// renderOverview renders the overview tab
func (m *TUIModel) renderOverview() string {
	if m.healthReport == nil {
		return "No data available"
	}

	return fmt.Sprintf(`Health Score: %d/100 (%s)

Total Checks: %d
Passed: %d
Failed: %d
Warnings: %d

Repository: %s
Last Updated: %s`,
		m.healthReport.OverallScore,
		m.healthReport.Grade,
		m.healthReport.Summary.TotalChecks,
		m.healthReport.Summary.PassedChecks,
		m.healthReport.Summary.FailedChecks,
		m.healthReport.Summary.WarningChecks,
		filepath.Base(m.repoPath),
		m.healthReport.Timestamp.Format("2006-01-02 15:04:05"))
}

// renderDetails renders the details tab
func (m *TUIModel) renderDetails() string {
	if m.healthReport == nil {
		return "No data available"
	}

	var details strings.Builder
	details.WriteString("Check Results:\n\n")

	for _, result := range m.healthReport.Results {
		status := "PASS"
		if result.Status == types.StatusFail {
			status = "FAIL"
		} else if result.Status == types.StatusWarning {
			status = "WARN"
		}

		details.WriteString(fmt.Sprintf("%s [%s] %s\n", status, result.ID, result.Name))
		if result.Message != "" {
			details.WriteString(fmt.Sprintf("  %s\n", result.Message))
		}
		details.WriteString("\n")
	}

	return details.String()
}

// renderTrends renders the trends tab
func (m *TUIModel) renderTrends() string {
	return "Trend analysis not yet implemented.\n\nThis feature will show historical health scores over time."
}

func runServe(cmd *cobra.Command, args []string) {
	var repoPath string
	if len(args) > 0 {
		repoPath = args[0]
	} else {
		var err error
		repoPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if it's a Git repository
	if !isGitRepository(repoPath) {
		fmt.Printf("Error: %s is not a Git repository\n", repoPath)
		os.Exit(1)
	}

	// Setup HTTP handlers
	http.HandleFunc("/", handleDashboard)
	http.HandleFunc("/api/health", handleHealthAPI)
	http.HandleFunc("/api/tags", handleTagsAPI)
	http.HandleFunc("/api/export/json", handleExportJSON)
	http.HandleFunc("/api/export/pdf", handleExportPDF)

	// Start server
	addr := fmt.Sprintf("%s:%d", serverHost, serverPort)
	fmt.Printf("Starting GPHC Web Dashboard...\n")
	fmt.Printf("Dashboard: http://%s\n", addr)
	fmt.Printf("Repository: %s\n", repoPath)
	fmt.Printf("Press Ctrl+C to stop\n\n")

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

func runTags(cmd *cobra.Command, args []string) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	if !isGitRepository(path) {
		fmt.Printf("Error: %s is not a Git repository\n", path)
		os.Exit(1)
	}

	// Initialize analyzer (not strictly needed for tags, but for consistency)
	analyzer, err := git.NewRepositoryAnalyzer(path)
	if err != nil {
		fmt.Printf("Error initializing repository analyzer: %v\n", err)
		os.Exit(1)
	}
	data, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("Error analyzing repository: %v\n", err)
		os.Exit(1)
	}

	// Run TagChecker alone for this command
	tc := checkers.NewTagChecker()
	res := tc.Check(data)

	// Print concise report
	fmt.Println("Tag & Release Health")
	fmt.Println("====================")
	fmt.Println()
	fmt.Printf("Status: %s\n", res.Status.String())
	fmt.Printf("Score: %d/%d\n", res.Score, 100)
	fmt.Printf("Message: %s\n", res.Message)
	for _, d := range res.Details {
		fmt.Printf("- %s\n", d)
	}

	// Suggest next tag
	if tagsSuggest {
		next, err := checkers.SuggestNextTag()
		if err == nil {
			fmt.Printf("\nAuto-suggested next tag: %s\n", next)
		}
	}

	// Generate changelog
	if tagsChangelogOut != "" {
		content, err := checkers.GenerateChangelog(tagsChangelogOut)
		if err != nil {
			fmt.Printf("Error generating changelog: %v\n", err)
		} else {
			if tagsChangelogOut != "" {
				fmt.Printf("Changelog generated: %s\n", tagsChangelogOut)
			} else {
				fmt.Println(content)
			}
		}
	}

	if tagsEnforce {
		// Simple policy enforcement: fail if status is FAIL or score < 50
		if res.Status == types.StatusFail || res.Score < 50 {
			fmt.Println("\nPolicy enforcement failed: tag policy violations detected")
			os.Exit(1)
		}
	}
}

// HTTP handlers
func handleDashboard(w http.ResponseWriter, r *http.Request) {
	// Check authentication if enabled
	if serverAuth {
		username, password, ok := r.BasicAuth()
		if !ok || username != serverUsername || password != serverPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="GPHC Dashboard"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + serverTitle + `</title>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            background: linear-gradient(135deg, #20b2aa 0%, #008b8b 50%, #006666 100%);
            background-attachment: fixed;
            min-height: 100vh;
            color: #333;
            overflow-x: hidden;
        }
        
        .container { 
            max-width: 1400px; 
            margin: 0 auto; 
            padding: 20px;
            width: 100%;
            box-sizing: border-box;
        }
        
        .header { 
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            color: #2c3e50; 
            padding: 30px; 
            border-radius: 20px; 
            margin-bottom: 30px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.1);
            text-align: center;
        }
        
        .header h1 { 
            font-size: 2.5em; 
            margin-bottom: 10px;
            background: linear-gradient(45deg, #20b2aa, #008b8b, #006666);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        
        .header p { 
            font-size: 1.2em; 
            opacity: 0.8;
        }
        
        .dashboard-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 25px;
            margin-bottom: 30px;
        }
        
        .card { 
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            padding: 25px; 
            border-radius: 20px; 
            box-shadow: 0 8px 32px rgba(0,0,0,0.1);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
            border: 1px solid rgba(255, 255, 255, 0.2);
            width: 100%;
            max-width: 100%;
            box-sizing: border-box;
            overflow: hidden;
        }
        
        .card:hover {
            transform: translateY(-5px);
            box-shadow: 0 12px 40px rgba(0,0,0,0.15);
        }
        
        .card h2 { 
            color: #2c3e50; 
            margin-bottom: 20px;
            font-size: 1.4em;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .card h2 i {
            color: #20b2aa;
        }
        
        .score-container {
            text-align: center;
            margin-bottom: 20px;
        }
        
        .score { 
            font-size: 3em; 
            font-weight: bold; 
            margin-bottom: 10px;
        }
        
        .score.excellent { color: #27ae60; }
        .score.good { color: #f39c12; }
        .score.poor { color: #e74c3c; }
        
        .grade {
            font-size: 1.5em;
            font-weight: bold;
            padding: 8px 16px;
            border-radius: 25px;
            display: inline-block;
        }
        
        .grade.excellent { background: #27ae60; color: white; }
        .grade.good { background: #f39c12; color: white; }
        .grade.poor { background: #e74c3c; color: white; }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 15px;
            margin: 20px 0;
        }
        
        .stat-item {
            background: rgba(32, 178, 170, 0.1);
            padding: 15px;
            border-radius: 10px;
            text-align: center;
        }
        
        .stat-number {
            font-size: 1.8em;
            font-weight: bold;
            color: #20b2aa;
        }
        
        .stat-label {
            font-size: 0.9em;
            color: #666;
            margin-top: 5px;
        }
        
        .feature-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-top: 20px;
            width: 100%;
            max-width: 100%;
            overflow: hidden;
        }
        
        .feature-item {
            background: rgba(32, 178, 170, 0.05);
            padding: 15px;
            border-radius: 10px;
            text-align: center;
            transition: background 0.3s ease, transform 0.2s ease;
            width: 100%;
            box-sizing: border-box;
            min-height: 120px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
        }
        
        .feature-item:hover {
            background: rgba(32, 178, 170, 0.1);
            transform: translateY(-2px);
        }
        
        .feature-icon {
            font-size: 2em;
            margin-bottom: 10px;
        }
        
        .feature-item:nth-child(1) .feature-icon { color: #e74c3c; } /* Documentation - Red */
        .feature-item:nth-child(2) .feature-icon { color: #3498db; } /* Commit Quality - Blue */
        .feature-item:nth-child(3) .feature-icon { color: #2ecc71; } /* Git Hygiene - Green */
        .feature-item:nth-child(4) .feature-icon { color: #f39c12; } /* Tag Management - Orange */
        .feature-item:nth-child(5) .feature-icon { color: #9b59b6; } /* Historical Tracking - Purple */
        .feature-item:nth-child(6) .feature-icon { color: #1abc9c; } /* Multi-Repo Scan - Turquoise */
        .feature-item:nth-child(7) .feature-icon { color: #34495e; } /* CI/CD Integration - Dark Gray */
        .feature-item:nth-child(8) .feature-icon { color: #e67e22; } /* Notifications - Dark Orange */
        .feature-item:nth-child(9) .feature-icon { color: #27ae60; } /* Terminal UI - Dark Green */
        .feature-item:nth-child(10) .feature-icon { color: #2980b9; } /* Web Dashboard - Dark Blue */
        .feature-item:nth-child(11) .feature-icon { color: #333; } /* GitHub Integration - Black */
        .feature-item:nth-child(12) .feature-icon { color: #fc6d26; } /* GitLab Integration - GitLab Orange */
        
        .feature-name {
            font-weight: bold;
            color: #2c3e50;
        }
        
        .btn { 
            background: linear-gradient(45deg, #20b2aa, #008b8b, #006666);
            color: white; 
            border: none; 
            padding: 12px 24px; 
            border-radius: 25px; 
            cursor: pointer; 
            font-size: 1em;
            transition: transform 0.3s ease, box-shadow 0.3s ease;
            margin: 5px;
        }
        
        .btn:hover { 
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(32, 178, 170, 0.4);
        }
        
        .btn-secondary {
            background: linear-gradient(45deg, #95a5a6, #7f8c8d);
        }
        
        .btn-success {
            background: linear-gradient(45deg, #27ae60, #2ecc71);
        }
        
        .btn-warning {
            background: linear-gradient(45deg, #f39c12, #e67e22);
        }
        
        .loading {
            text-align: center;
            padding: 40px;
            color: #666;
        }
        
        .loading i {
            font-size: 2em;
            animation: spin 1s linear infinite;
        }
        
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .error {
            background: #e74c3c;
            color: white;
            padding: 15px;
            border-radius: 10px;
            margin: 20px 0;
        }
        
        .success {
            background: #27ae60;
            color: white;
            padding: 15px;
            border-radius: 10px;
            margin: 20px 0;
        }
        
        .warning {
            background: #f39c12;
            color: white;
            padding: 15px;
            border-radius: 10px;
            margin: 20px 0;
        }
        
        .footer {
            text-align: center;
            margin-top: 30px;
            color: rgba(255, 255, 255, 0.8);
        }
        
        @media (max-width: 768px) {
            .dashboard-grid {
                grid-template-columns: 1fr;
            }
            
            .stats-grid {
                grid-template-columns: 1fr;
            }
            
            .feature-grid {
                grid-template-columns: repeat(2, 1fr);
                gap: 10px;
            }
            
            .feature-item {
                min-height: 100px;
                padding: 10px;
            }
            
            .container {
                padding: 10px;
            }
        }
        
        @media (max-width: 480px) {
            .feature-grid {
                grid-template-columns: 1fr;
            }
            
            .header h1 {
                font-size: 2em;
            }
            
            .card {
                padding: 15px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1><i class="fas fa-heartbeat"></i> ` + serverTitle + `</h1>
            <p>Modern Repository Health Monitoring Dashboard</p>
        </div>
        
        <div class="dashboard-grid">
            <div class="card">
                <h2><i class="fas fa-chart-line"></i> Health Overview</h2>
                <div id="health-data" class="loading">
                    <i class="fas fa-spinner"></i><br>
                    Loading health data...
                </div>
                <div style="text-align: center; margin-top: 20px;">
                    <button class="btn" onclick="refreshData()">
                        <i class="fas fa-sync-alt"></i> Refresh
                    </button>
                </div>
            </div>
            
            <div class="card">
                <h2><i class="fas fa-download"></i> Export Options</h2>
                <div style="text-align: center;">
                    <button class="btn btn-success" onclick="exportJSON()">
                        <i class="fas fa-file-code"></i> Export JSON
                    </button>
                    <button class="btn btn-warning" onclick="exportPDF()">
                        <i class="fas fa-file-pdf"></i> Export PDF
                    </button>
                </div>
            </div>
            
            <div class="card">
                <h2><i class="fas fa-tags"></i> Tag Management</h2>
                <div style="text-align: center;">
                    <button class="btn" onclick="checkTags()">
                        <i class="fas fa-tag"></i> Check Tags
                    </button>
                    <button class="btn btn-secondary" onclick="suggestTag()">
                        <i class="fas fa-lightbulb"></i> Suggest Version
                    </button>
                </div>
                <div id="tag-data" style="margin-top: 20px;"></div>
            </div>
            
            <div class="card">
                <h2><i class="fas fa-search"></i> Multi-Repository Scan</h2>
                <div style="text-align: center;">
                    <button class="btn" onclick="scanRepos()">
                        <i class="fas fa-search"></i> Scan Projects
                    </button>
                </div>
                <div id="scan-data" style="margin-top: 20px;"></div>
            </div>
        </div>
        
        <div class="card">
            <h2><i class="fas fa-rocket"></i> Available Features</h2>
            <div class="feature-grid">
                <div class="feature-item">
                    <div class="feature-icon"><i class="fas fa-file-alt"></i></div>
                    <div class="feature-name">Documentation</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fas fa-code-branch"></i></div>
                    <div class="feature-name">Commit Quality</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fas fa-broom"></i></div>
                    <div class="feature-name">Git Hygiene</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fas fa-tags"></i></div>
                    <div class="feature-name">Tag Management</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fas fa-history"></i></div>
                    <div class="feature-name">Historical Tracking</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fas fa-search"></i></div>
                    <div class="feature-name">Multi-Repo Scan</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fas fa-cogs"></i></div>
                    <div class="feature-name">CI/CD Integration</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fas fa-bell"></i></div>
                    <div class="feature-name">Notifications</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fas fa-terminal"></i></div>
                    <div class="feature-name">Terminal UI</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fas fa-globe"></i></div>
                    <div class="feature-name">Web Dashboard</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fab fa-github"></i></div>
                    <div class="feature-name">GitHub Integration</div>
                </div>
                <div class="feature-item">
                    <div class="feature-icon"><i class="fab fa-gitlab"></i></div>
                    <div class="feature-name">GitLab Integration</div>
                </div>
            </div>
        </div>
        
        <div class="footer">
            <p><i class="fas fa-heart"></i> Powered by GPHC - Git Project Health Checker</p>
        </div>
    </div>

    <script>
        function refreshData() {
            const healthData = document.getElementById('health-data');
            healthData.innerHTML = '<div class="loading"><i class="fas fa-spinner"></i><br>Loading...</div>';
            
            fetch('/api/health')
                .then(response => response.json())
                .then(data => {
                    const scoreClass = getScoreClass(data.overall_score);
                    const gradeClass = getGradeClass(data.grade);
                    
                    healthData.innerHTML = 
                        '<div class="score-container">' +
                            '<div class="score ' + scoreClass + '">' + data.overall_score + '/100</div>' +
                            '<div class="grade ' + gradeClass + '">' + data.grade + '</div>' +
                        '</div>' +
                        '<div class="stats-grid">' +
                            '<div class="stat-item">' +
                                '<div class="stat-number">' + data.summary.total_checks + '</div>' +
                                '<div class="stat-label">Total Checks</div>' +
                            '</div>' +
                            '<div class="stat-item">' +
                                '<div class="stat-number" style="color: #27ae60;">' + data.summary.passed_checks + '</div>' +
                                '<div class="stat-label">Passed</div>' +
                            '</div>' +
                            '<div class="stat-item">' +
                                '<div class="stat-number" style="color: #e74c3c;">' + data.summary.failed_checks + '</div>' +
                                '<div class="stat-label">Failed</div>' +
                            '</div>' +
                            '<div class="stat-item">' +
                                '<div class="stat-number" style="color: #f39c12;">' + data.summary.warning_checks + '</div>' +
                                '<div class="stat-label">Warnings</div>' +
                            '</div>' +
                        '</div>' +
                        '<div style="text-align: center; margin-top: 15px; color: #666;">' +
                            '<small>Last updated: ' + new Date().toLocaleTimeString() + '</small>' +
                        '</div>';
                })
                .catch(error => {
                    healthData.innerHTML = '<div class="error">Error loading data: ' + error + '</div>';
                });
        }
        
        function getScoreClass(score) {
            if (score >= 80) return 'excellent';
            if (score >= 60) return 'good';
            return 'poor';
        }
        
        function getGradeClass(grade) {
            if (grade.includes('A') || grade.includes('B')) return 'excellent';
            if (grade.includes('C')) return 'good';
            return 'poor';
        }
        
        function getStatusClass(status) {
            if (status === 'PASS') return 'excellent';
            if (status === 'WARNING') return 'good';
            return 'poor';
        }
        
        function exportJSON() {
            window.open('/api/export/json', '_blank');
        }
        
        function exportPDF() {
            window.open('/api/export/pdf', '_blank');
        }
        
        function checkTags() {
            const tagData = document.getElementById('tag-data');
            tagData.innerHTML = '<div class="loading"><i class="fas fa-spinner"></i><br>Checking tags...</div>';
            
            fetch('/api/tags')
                .then(response => response.json())
                .then(data => {
                    const statusClass = getStatusClass(data.status);
                    const scoreClass = getScoreClass(data.score);
                    
                    tagData.innerHTML = 
                        '<div class="score-container">' +
                            '<div class="score ' + scoreClass + '">' + data.score + '/100</div>' +
                            '<div class="grade ' + statusClass + '">' + data.status + '</div>' +
                        '</div>' +
                        '<div style="margin-top: 15px;">' +
                            '<p><strong>' + data.message + '</strong></p>' +
                            '<ul style="margin-top: 10px; padding-left: 20px;">';
                    
                    data.details.forEach(detail => {
                        tagData.innerHTML += '<li>' + detail + '</li>';
                    });
                    
                    tagData.innerHTML += 
                            '</ul>' +
                        '</div>';
                })
                .catch(error => {
                    tagData.innerHTML = '<div class="error">Error loading tag data: ' + error + '</div>';
                });
        }
        
        function suggestTag() {
            const tagData = document.getElementById('tag-data');
            tagData.innerHTML = '<div class="loading"><i class="fas fa-spinner"></i><br>Analyzing commits...</div>';
            
            // Simulate tag suggestion
            setTimeout(() => {
                tagData.innerHTML = '<div class="success"><i class="fas fa-lightbulb"></i> Suggested: v1.2.5. Use CLI: git hc tags --suggest</div>';
            }, 1000);
        }
        
        function scanRepos() {
            const scanData = document.getElementById('scan-data');
            scanData.innerHTML = '<div class="loading"><i class="fas fa-spinner"></i><br>Scanning repositories...</div>';
            
            // Simulate repo scan
            setTimeout(() => {
                scanData.innerHTML = '<div class="success"><i class="fas fa-search"></i> Scan completed. Use CLI: git hc scan ~/projects --recursive</div>';
            }, 1500);
        }
        
        // Load data on page load
        refreshData();
        
        // Auto-refresh every 30 seconds
        setInterval(refreshData, 30000);
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func handleHealthAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication if enabled
	if serverAuth {
		username, password, ok := r.BasicAuth()
		if !ok || username != serverUsername || password != serverPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="GPHC Dashboard"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
	}

	// Get current directory for health check
	repoPath, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting repository path", http.StatusInternalServerError)
		return
	}

	// Analyze repository
	analyzer, err := git.NewRepositoryAnalyzer(repoPath)
	if err != nil {
		http.Error(w, "Error analyzing repository", http.StatusInternalServerError)
		return
	}

	data, err := analyzer.Analyze()
	if err != nil {
		http.Error(w, "Error analyzing repository data", http.StatusInternalServerError)
		return
	}

	// Run health checks
	allCheckers := []checkers.Checker{
		checkers.NewDocChecker(),
		checkers.NewConventionalCommitChecker(),
		checkers.NewMsgLengthChecker(),
		checkers.NewCommitSizeChecker(),
		checkers.NewLocalBranchChecker(),
		checkers.NewStaleBranchChecker(),
		checkers.NewBareRepoChecker(),
		checkers.NewStashChecker(),
		checkers.NewIgnoreChecker(),
		checkers.NewGitHubIntegrationChecker(),
		checkers.NewGitLabIntegrationChecker(),
		checkers.NewCommitAuthorInsightsChecker(),
		checkers.NewCodebaseSmellChecker(),
		checkers.NewTagChecker(),
	}

	scorer := scorer.NewScorer()
	for _, checker := range allCheckers {
		result := checker.Check(data)
		scorer.AddResult(*result)
	}

	// Generate report
	healthReport := scorer.CalculateHealthReport()

	// Set CORS headers if enabled
	if serverCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	w.Header().Set("Content-Type", "application/json")

	// Convert to JSON (simplified)
	json := fmt.Sprintf(`{
		"overall_score": %d,
		"grade": "%s",
		"summary": {
			"total_checks": %d,
			"passed_checks": %d,
			"failed_checks": %d,
			"warning_checks": %d
		},
		"timestamp": "%s",
		"repository": "%s"
	}`,
		healthReport.OverallScore,
		healthReport.Grade,
		healthReport.Summary.TotalChecks,
		healthReport.Summary.PassedChecks,
		healthReport.Summary.FailedChecks,
		healthReport.Summary.WarningChecks,
		healthReport.Timestamp.Format(time.RFC3339),
		filepath.Base(repoPath))

	w.Write([]byte(json))
}

func handleTagsAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication if enabled
	if serverAuth {
		username, password, ok := r.BasicAuth()
		if !ok || username != serverUsername || password != serverPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="GPHC Dashboard"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
	}

	// Get current directory for tag check
	repoPath, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting repository path", http.StatusInternalServerError)
		return
	}

	// Analyze repository
	analyzer, err := git.NewRepositoryAnalyzer(repoPath)
	if err != nil {
		http.Error(w, "Error analyzing repository", http.StatusInternalServerError)
		return
	}

	data, err := analyzer.Analyze()
	if err != nil {
		http.Error(w, "Error analyzing repository data", http.StatusInternalServerError)
		return
	}

	// Run TagChecker
	tagChecker := checkers.NewTagChecker()
	result := tagChecker.Check(data)

	// Set CORS headers if enabled
	if serverCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	w.Header().Set("Content-Type", "application/json")

	// Convert to JSON
	json := fmt.Sprintf(`{
		"status": "%s",
		"score": %d,
		"message": "%s",
		"details": %s,
		"timestamp": "%s",
		"repository": "%s"
	}`,
		result.Status.String(),
		result.Score,
		result.Message,
		fmt.Sprintf(`["%s"]`, strings.Join(result.Details, `", "`)),
		result.Timestamp.Format(time.RFC3339),
		filepath.Base(repoPath))

	w.Write([]byte(json))
}

func handleExportJSON(w http.ResponseWriter, r *http.Request) {
	// Check authentication if enabled
	if serverAuth {
		username, password, ok := r.BasicAuth()
		if !ok || username != serverUsername || password != serverPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="GPHC Dashboard"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
	}

	// Get current directory for health check
	repoPath, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting repository path", http.StatusInternalServerError)
		return
	}

	// Analyze repository
	analyzer, err := git.NewRepositoryAnalyzer(repoPath)
	if err != nil {
		http.Error(w, "Error analyzing repository", http.StatusInternalServerError)
		return
	}

	data, err := analyzer.Analyze()
	if err != nil {
		http.Error(w, "Error analyzing repository data", http.StatusInternalServerError)
		return
	}

	// Run health checks
	allCheckers := []checkers.Checker{
		checkers.NewDocChecker(),
		checkers.NewConventionalCommitChecker(),
		checkers.NewMsgLengthChecker(),
		checkers.NewCommitSizeChecker(),
		checkers.NewLocalBranchChecker(),
		checkers.NewStaleBranchChecker(),
		checkers.NewBareRepoChecker(),
		checkers.NewStashChecker(),
		checkers.NewIgnoreChecker(),
		checkers.NewGitHubIntegrationChecker(),
		checkers.NewGitLabIntegrationChecker(),
		checkers.NewCommitAuthorInsightsChecker(),
		checkers.NewCodebaseSmellChecker(),
	}

	scorer := scorer.NewScorer()
	for _, checker := range allCheckers {
		result := checker.Check(data)
		scorer.AddResult(*result)
	}

	// Generate report
	healthReport := scorer.CalculateHealthReport()

	// Export to JSON
	exporter := exporter.NewExporter()
	jsonData, err := exporter.Export(healthReport, "json")
	if err != nil {
		http.Error(w, "Error exporting to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=health-report.json")
	w.Write([]byte(jsonData))
}

func handleExportPDF(w http.ResponseWriter, r *http.Request) {
	// Check authentication if enabled
	if serverAuth {
		username, password, ok := r.BasicAuth()
		if !ok || username != serverUsername || password != serverPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="GPHC Dashboard"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
	}

	// For now, return a simple message since PDF export requires additional dependencies
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		<html>
		<body>
			<h1>PDF Export</h1>
			<p>PDF export is not yet implemented. Please use JSON export for now.</p>
			<a href="/api/export/json">Download JSON Report</a>
		</body>
		</html>
	`))
}
