package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vahidaghazadeh/gphc/internal/checkers"
	"github.com/vahidaghazadeh/gphc/internal/exporter"
	"github.com/vahidaghazadeh/gphc/internal/git"
	"github.com/vahidaghazadeh/gphc/internal/reporter"
	"github.com/vahidaghazadeh/gphc/internal/scorer"
	"github.com/vahidaghazadeh/gphc/pkg/config"
	"github.com/vahidaghazadeh/gphc/pkg/types"
)

var (
	version = "dev"
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
- Overview and detailed check results
- Manual refresh without leaving the interface
- Keyboard navigation`,
	Args: cobra.MaximumNArgs(1),
	Run:  runTUI,
}

var serveCmd = &cobra.Command{
	Use:   "serve [path]",
	Short: "Start web dashboard server",
	Long: `Start a local web server to display health monitoring dashboard.
Provides a web interface accessible via browser with:
- Repository health monitoring
- Git diff and tag information
- JSON export and API access
- Optional basic authentication`,
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
	Use:     "suggest [path]",
	Aliases: []string{"comment"},
	Short:   "Suggest commit message based on staged changes",
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

var securityCmd = &cobra.Command{
	Use:   "security",
	Short: "Security scanning and analysis",
	Long:  `Perform security scans including secret detection in Git history`,
}

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Scan for secrets in Git history",
	Long:  `Scan the entire Git history and stashes for exposed secrets and credentials`,
	Run:   runSecretsScan,
}

var dependenciesCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Scan transitive dependencies for vulnerabilities",
	Long: `Perform deep analysis of transitive dependencies to detect security vulnerabilities.
This includes both direct and indirect dependencies, helping identify supply chain attacks.

Examples:
  git hc security dependencies                    # Basic dependency scan
  git hc security dependencies --depth deep       # Deep transitive analysis
  git hc security dependencies --format json      # JSON output format
  git hc security dependencies --severity high    # Only show high/critical vulnerabilities`,
	Run: runDependenciesScan,
}

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Validate Git security policies and configurations",
	Long: `Validate Git security policies including commit signatures, push policies,
sensitive file detection, and branch protection settings.

Examples:
  git hc security policy                          # Basic policy validation
  git hc security policy --check-signing          # Focus on commit signatures
  git hc security policy --check-files            # Focus on sensitive files
  git hc security policy --format json            # JSON output format
  git hc security policy --severity high          # Only show high/critical issues`,
	Run: runPolicyValidation,
}

var binariesCmd = &cobra.Command{
	Use:   "binaries",
	Short: "Audit executable and large files for security risks",
	Long: `Scan repository for executable files, large files, and suspicious file types
that pose security risks or repository health issues.

Examples:
  git hc security binaries                        # Basic binary audit
  git hc security binaries --max-size 50mb        # Set custom size threshold
  git hc security binaries --check-history        # Include Git history scan
  git hc security binaries --format json          # JSON output format
  git hc security binaries --severity high        # Only show high/critical issues`,
	Run: runBinariesAudit,
}

func init() {
	rootCmd.Version = effectiveVersion()
	secretsCmd.Flags().Bool("history", true, "Scan entire Git history for secrets")
	secretsCmd.Flags().Bool("stashes", true, "Scan Git stashes for secrets")
	secretsCmd.Flags().Bool("entropy", true, "Perform entropy analysis for random strings")
	secretsCmd.Flags().String("severity", "medium", "Minimum severity level (low, medium, high)")
	secretsCmd.Flags().Float64("confidence", 0.8, "Minimum confidence threshold (0.0-1.0)")
	secretsCmd.Flags().String("format", "table", "Output format (table, json, yaml)")
	secretsCmd.Flags().String("output", "", "Output file path")

	dependenciesCmd.Flags().String("depth", "deep", "Scan depth (shallow, deep)")
	dependenciesCmd.Flags().String("severity", "low", "Minimum severity level (low, medium, high, critical)")
	dependenciesCmd.Flags().String("format", "table", "Output format (table, json, yaml)")
	dependenciesCmd.Flags().String("output", "", "Output file path")
	dependenciesCmd.Flags().Bool("tree", true, "Show dependency tree structure")
	dependenciesCmd.Flags().Bool("direct-only", false, "Only check direct dependencies")

	policyCmd.Flags().Bool("check-signing", true, "Check commit signature verification")
	policyCmd.Flags().Bool("check-files", true, "Check for sensitive files")
	policyCmd.Flags().Bool("check-push", true, "Check push policies")
	policyCmd.Flags().Bool("check-branches", true, "Check branch protection")
	policyCmd.Flags().String("severity", "low", "Minimum severity level (low, medium, high, critical)")
	policyCmd.Flags().String("format", "table", "Output format (table, json, yaml)")
	policyCmd.Flags().String("output", "", "Output file path")

	binariesCmd.Flags().String("max-size", "10mb", "Maximum file size threshold (e.g., 10mb, 50mb, 100mb)")
	binariesCmd.Flags().Bool("check-history", true, "Check Git history for binary files")
	binariesCmd.Flags().Bool("check-executables", true, "Check for executable files")
	binariesCmd.Flags().Bool("check-large", true, "Check for large files")
	binariesCmd.Flags().Bool("check-suspicious", true, "Check for suspicious file types")
	binariesCmd.Flags().String("severity", "low", "Minimum severity level (low, medium, high, critical)")
	binariesCmd.Flags().String("format", "table", "Output format (table, json, yaml)")
	binariesCmd.Flags().String("output", "", "Output file path")

	securityCmd.AddCommand(secretsCmd)
	securityCmd.AddCommand(dependenciesCmd)
	securityCmd.AddCommand(policyCmd)
	securityCmd.AddCommand(binariesCmd)
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
	rootCmd.AddCommand(securityCmd)

	// Add export format flags
	checkCmd.Flags().StringVarP(&exportFormat, "format", "f", "terminal", "Output format: terminal, json, yaml, markdown, html")
	checkCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: stdout)")

	// Add pre-commit command flags
	preCommitCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Repository path to check")

	// Add suggest command flags
	suggestCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Repository path to analyze")

	// Add commit command flags
	commitCmd.Flags().Bool("suggest", false, "Suggest commit message based on staged changes (default behavior)")
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
	serveCmd.Flags().StringVarP(&serverUsername, "username", "u", "", "Username for basic authentication")
	serveCmd.Flags().StringVarP(&serverPassword, "password", "w", "", "Password for basic authentication")
	serveCmd.Flags().BoolVarP(&serverCORS, "cors", "c", false, "Enable permissive CORS headers")
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
	serverRepoPath  string

	// tags command flags
	tagsSuggest      bool
	tagsChangelogOut string
	tagsEnforce      bool

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
		fmt.Printf("GPHC (Git Project Health Checker) v%s\n", effectiveVersion())
	},
}

func effectiveVersion() string {
	if version != "" && version != "dev" {
		return strings.TrimPrefix(version, "v")
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		moduleVersion := strings.TrimPrefix(info.Main.Version, "v")
		if moduleVersion != "" && moduleVersion != "(devel)" {
			return moduleVersion
		}
	}
	return "dev"
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update GPHC to the latest version",
	Long: `Update GPHC to the latest version from GitHub.
This command will:
1. Install the latest published GPHC version with go install
2. Prefer the directory of the currently running gphc binary
3. Show how to verify the installed version`,
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
	if !checkLargeFiles(path, stagedFiles) {
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
	path, suggestion, stagedCount, ok := suggestedCommitForArgs(args)
	if !ok {
		return
	}

	printCommitSuggestion(stagedCount, suggestion)
	fmt.Print("\nCreate commit with this message? [y/N]: ")

	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Printf("\nCould not read confirmation: %v\n", err)
		return
	}

	if !isAffirmativeAnswer(answer) {
		fmt.Println("Commit cancelled.")
		return
	}

	if err := commitStagedChanges(path, suggestion); err != nil {
		fmt.Printf("Error creating commit: %v\n", err)
		os.Exit(1)
	}
}

func runSuggest(cmd *cobra.Command, args []string) {
	_, suggestion, stagedCount, ok := suggestedCommitForArgs(args)
	if !ok {
		return
	}
	printCommitSuggestion(stagedCount, suggestion)
}

func suggestedCommitForArgs(args []string) (string, string, int, bool) {
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
		return path, "", 0, false
	}

	suggestion := analyzeStagedChanges(path, stagedFiles)
	return path, suggestion, len(stagedFiles), true
}

func printCommitSuggestion(stagedCount int, suggestion string) {
	fmt.Printf("Staged changes: %d file(s)\n", stagedCount)
	fmt.Println("\nCommit suggestion:")
	fmt.Printf("  %s\n", suggestion)
	fmt.Println("\nCommit with:")
	fmt.Printf("  git commit -m %s\n", shellQuote(suggestion))
}

func isAffirmativeAnswer(answer string) bool {
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "y", "yes", "بله", "آره", "اره":
		return true
	default:
		return false
	}
}

func commitStagedChanges(path, message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
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

	healthReport, err := buildHealthReport(path)
	if err != nil {
		fmt.Printf("Error running health check: %v\n", err)
		return
	}

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

func loadRepositoryConfig(repoPath string) (*config.Config, error) {
	configPath := filepath.Join(repoPath, "gphc.yml")
	if _, err := os.Stat(configPath); err == nil {
		return config.LoadConfig(configPath)
	}
	return config.DefaultConfig(), nil
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

	installDir := ""
	if exePath, err := os.Executable(); err == nil {
		installDir = filepath.Dir(exePath)
	}

	installCmd := exec.Command("go", "install", "github.com/vahidaghazadeh/gphc/cmd/gphc@latest")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if installDir != "" {
		installCmd.Env = append(os.Environ(), "GOBIN="+installDir)
		fmt.Printf("Installing into: %s\n", installDir)
	}

	if err := installCmd.Run(); err != nil {
		fmt.Printf("Error installing GPHC: %v\n", err)
		return
	}

	fmt.Println("GPHC updated successfully!")
	fmt.Println("Run 'gphc version' to verify the installed version.")
}

func isGitRepository(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	return cmd.Run() == nil
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

func checkLargeFiles(repoPath string, files []string) bool {
	for _, file := range files {
		info, err := os.Stat(filepath.Join(repoPath, file))
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
		base := filepath.Base(file)
		for _, pattern := range sensitivePatterns {
			matched, _ := filepath.Match(pattern, base)
			if matched || base == pattern {
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

	repos = filterRepositories(repos, includePatterns, excludePatterns)
	if parallelJobs < 1 {
		parallelJobs = 1
	}
	if parallelJobs > len(repos) {
		parallelJobs = len(repos)
	}

	type scanOutcome struct {
		result ScanResult
		err    error
	}
	jobs := make(chan string)
	outcomes := make(chan scanOutcome)
	for worker := 0; worker < parallelJobs; worker++ {
		go func() {
			for repo := range jobs {
				report, err := buildHealthReport(repo)
				if err != nil {
					outcomes <- scanOutcome{result: ScanResult{Path: repo}, err: err}
					continue
				}
				outcomes <- scanOutcome{result: ScanResult{
					Name:  filepath.Base(repo),
					Path:  repo,
					Score: report.OverallScore,
					Grade: report.Grade,
				}}
			}
		}()
	}
	go func() {
		for _, repo := range repos {
			jobs <- repo
		}
		close(jobs)
	}()

	results := make([]ScanResult, 0, len(repos))
	for range repos {
		outcome := <-outcomes
		if outcome.err != nil {
			fmt.Printf("Error scanning %s: %v\n", outcome.result.Path, outcome.err)
			continue
		}
		if minScore == 0 || outcome.result.Score >= minScore {
			results = append(results, outcome.result)
		}
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Path < results[j].Path })

	totalScore := 0.0
	for _, result := range results {
		totalScore += float64(result.Score)
		if detailedReport {
			fmt.Printf("%s\n  Path: %s\n  Score: %d/100 (%s)\n", result.Name, result.Path, result.Score, result.Grade)
		} else {
			fmt.Printf("%s: %d/100 (%s)\n", result.Name, result.Score, result.Grade)
		}
	}

	if scanOutputFile != "" {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Printf("Error encoding scan results: %v\n", err)
			return
		}
		if err := os.WriteFile(scanOutputFile, data, 0644); err != nil {
			fmt.Printf("Error writing scan results: %v\n", err)
			return
		}
		fmt.Printf("Results written to %s\n", scanOutputFile)
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

func filterRepositories(repos, includes, excludes []string) []string {
	filtered := make([]string, 0, len(repos))
	for _, repo := range repos {
		if len(includes) > 0 && !matchesRepositoryPattern(repo, includes) {
			continue
		}
		if matchesRepositoryPattern(repo, excludes) {
			continue
		}
		filtered = append(filtered, repo)
	}
	return filtered
}

func matchesRepositoryPattern(repo string, patterns []string) bool {
	for _, pattern := range patterns {
		baseMatch, _ := filepath.Match(pattern, filepath.Base(repo))
		pathMatch, _ := filepath.Match(pattern, repo)
		if baseMatch || pathMatch || strings.Contains(repo, pattern) {
			return true
		}
	}
	return false
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

	if serverAuth && (serverUsername == "" || serverPassword == "") {
		fmt.Println("Error: --auth requires non-empty --username and --password")
		return
	}

	absolutePath, err := filepath.Abs(repoPath)
	if err != nil {
		fmt.Printf("Error resolving repository path: %v\n", err)
		return
	}
	serverRepoPath = absolutePath

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleDashboard)
	mux.HandleFunc("/api/health", handleHealthAPI)
	mux.HandleFunc("/api/tags", handleTagsAPI)
	mux.HandleFunc("/api/diff", handleDiffAPI)
	mux.HandleFunc("/api/export/json", handleExportJSON)

	// Start server
	addr := fmt.Sprintf("%s:%d", serverHost, serverPort)
	fmt.Printf("Starting GPHC Web Dashboard...\n")
	fmt.Printf("Dashboard: http://%s\n", addr)
	fmt.Printf("Repository: %s\n", serverRepoPath)
	fmt.Printf("Press Ctrl+C to stop\n\n")

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
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
		next, err := checkers.SuggestNextTag(path)
		if err == nil {
			fmt.Printf("\nAuto-suggested next tag: %s\n", next)
		}
	}

	// Generate changelog
	if tagsChangelogOut != "" {
		content, err := checkers.GenerateChangelog(path, tagsChangelogOut)
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
