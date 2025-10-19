package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

func init() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(preCommitCmd)
	rootCmd.AddCommand(badgeCmd)
	rootCmd.AddCommand(githubCmd)
	rootCmd.AddCommand(gitlabCmd)
	rootCmd.AddCommand(authorsCmd)
	rootCmd.AddCommand(codebaseCmd)

	// Add export format flags
	checkCmd.Flags().StringVarP(&exportFormat, "format", "f", "terminal", "Output format: terminal, json, yaml, markdown, html")
	checkCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: stdout)")
}

var (
	exportFormat string
	outputFile   string
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
	Use:   "pre-commit",
	Short: "Run quick pre-commit checks on staged files",
	Long: `Run quick health checks suitable for pre-commit hooks.
This command performs fast checks on staged files and current commit.
Designed for integration with pre-commit framework and Husky.
Returns non-zero exit code if issues are found.`,
	Run: runPreCommit,
}

func runPreCommit(cmd *cobra.Command, args []string) {
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
	cmd := exec.Command("git", "log", "-1", "--pretty=%s")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	message := strings.TrimSpace(string(output))

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
