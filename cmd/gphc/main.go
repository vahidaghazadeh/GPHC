package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/opsource/gphc/internal/checkers"
	"github.com/opsource/gphc/internal/git"
	"github.com/opsource/gphc/internal/reporter"
	"github.com/opsource/gphc/internal/scorer"
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

func init() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
}

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
2. Rebuild and reinstall GPHC
3. Show the new version`,
	Run: runUpdate,
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

	fmt.Printf("🔍 Analyzing repository: %s\n", path)

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
		checkers.NewLocalBranchChecker(),
		checkers.NewStaleBranchChecker(),
		checkers.NewBareRepoChecker(),
		checkers.NewStashChecker(),
	}

	// Run all checkers
	scorer := scorer.NewScorer()
	for _, checker := range allCheckers {
		result := checker.Check(data)
		scorer.AddResult(*result)
	}

	// Generate report
	healthReport := scorer.CalculateHealthReport()

	// Display results
	reporter := reporter.NewReporter()
	output := reporter.Report(healthReport)
	fmt.Println(output)
}

func runUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("🔄 Updating GPHC...")

	// Find the GPHC source directory
	sourceDir := findGPHCSourceDir()
	if sourceDir == "" {
		fmt.Println("❌ Error: Could not find GPHC source directory")
		fmt.Println("💡 Please run this command from the GPHC project directory")
		os.Exit(1)
	}

	fmt.Printf("📁 Found GPHC source at: %s\n", sourceDir)

	// Change to source directory
	if err := os.Chdir(sourceDir); err != nil {
		fmt.Printf("❌ Error changing to source directory: %v\n", err)
		os.Exit(1)
	}

	// Pull latest changes
	fmt.Println("📥 Pulling latest changes...")
	pullCmd := exec.Command("git", "pull", "origin", "main")
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr

	if err := pullCmd.Run(); err != nil {
		fmt.Printf("❌ Error pulling changes: %v\n", err)
		fmt.Println("💡 Make sure you have internet connection and git access")
		os.Exit(1)
	}

	// Rebuild and reinstall
	fmt.Println("🔨 Building and installing GPHC...")
	installCmd := exec.Command("go", "install", "./cmd/gphc")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	if err := installCmd.Run(); err != nil {
		fmt.Printf("❌ Error installing GPHC: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ GPHC updated successfully!")
	fmt.Println("📊 New version:")

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
