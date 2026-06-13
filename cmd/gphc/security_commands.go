package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vahidaghazadeh/gphc/internal/checkers"
	"github.com/vahidaghazadeh/gphc/internal/git"
	"github.com/vahidaghazadeh/gphc/pkg/types"
	"gopkg.in/yaml.v2"
)

func runSecretsScan(cmd *cobra.Command, args []string) {
	// Get flags
	scanHistory, _ := cmd.Flags().GetBool("history")
	scanStashes, _ := cmd.Flags().GetBool("stashes")
	scanEntropy, _ := cmd.Flags().GetBool("entropy")
	minSeverity, _ := cmd.Flags().GetString("severity")
	minConfidence, _ := cmd.Flags().GetFloat64("confidence")
	format, _ := cmd.Flags().GetString("format")
	outputFile, _ := cmd.Flags().GetString("output")

	// Determine repository path
	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}

	// Check if it's a Git repository
	if !isGitRepository(repoPath) {
		fmt.Printf("Error: %s is not a Git repository\n", repoPath)
		os.Exit(1)
	}

	fmt.Printf("🔍 Scanning for secrets in Git history...\n")
	fmt.Printf("Repository: %s\n", repoPath)
	fmt.Printf("Scanning history: %v\n", scanHistory)
	fmt.Printf("Scanning stashes: %v\n", scanStashes)
	fmt.Printf("Entropy analysis: %v\n", scanEntropy)
	fmt.Printf("Minimum severity: %s\n", minSeverity)
	fmt.Printf("Minimum confidence: %.2f\n\n", minConfidence)

	// Run secret checker
	secretChecker := checkers.NewSecretChecker()

	// Create RepositoryData for the checker
	analyzer, err := git.NewRepositoryAnalyzer(repoPath)
	if err != nil {
		fmt.Printf("Error initializing repository analyzer: %v\n", err)
		os.Exit(1)
	}
	data, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("Error analyzing repository: %v\n", err)
		os.Exit(1)
	}

	result := secretChecker.CheckWithOptions(data, scanHistory, scanStashes, scanEntropy, minSeverity, minConfidence)

	switch format {
	case "json":
		outputSecurityResult(result, outputFile, "json")
	case "yaml":
		outputSecurityResult(result, outputFile, "yaml")
	default:
		for _, detail := range result.Details {
			fmt.Printf("• %s\n", detail)
		}
	}

	if result.Status == types.StatusFail && len(result.Details) > 0 {
		fmt.Printf("🚨 Secrets found in Git history!\n\n")

		// Display details
		for _, detail := range result.Details {
			fmt.Printf("• %s\n", detail)
		}

		// Show remediation
		fmt.Printf("\n🚨 CRITICAL: Secrets found in Git history!\n\n")
		fmt.Printf("Immediate Actions Required:\n")
		fmt.Printf("1. Rotate/revoke all exposed credentials immediately\n")
		fmt.Printf("2. Rewrite Git history to remove secrets\n")
		fmt.Printf("3. Notify team members about the exposure\n\n")

		fmt.Printf("Tools for History Rewriting:\n")
		fmt.Printf("- git filter-repo: https://github.com/newren/git-filter-repo\n")
		fmt.Printf("- BFG Repo-Cleaner: https://rtyley.github.io/bfg-repo-cleaner/\n\n")

		fmt.Printf("Commands:\n")
		fmt.Printf("# Using git filter-repo\n")
		fmt.Printf("git filter-repo --replace-text <(echo 'SECRET_VALUE==>REDACTED')\n\n")
		fmt.Printf("# Using BFG\n")
		fmt.Printf("java -jar bfg.jar --replace-text replacements.txt\n\n")

		fmt.Printf("After rewriting history:\n")
		fmt.Printf("git push --force-with-lease origin main\n")

		os.Exit(1)
	} else {
		fmt.Printf("✅ No secrets found.\n")
	}
}

func outputSecurityResult(result *types.CheckResult, outputFile, format string) {
	var (
		data []byte
		err  error
	)
	if format == "yaml" {
		data, err = yaml.Marshal(result)
	} else {
		data, err = json.MarshalIndent(result, "", "  ")
	}
	if err != nil {
		fmt.Printf("Error encoding result: %v\n", err)
		return
	}
	if outputFile == "" {
		fmt.Printf("%s\n", data)
		return
	}
	if err := os.WriteFile(outputFile, data, 0600); err != nil {
		fmt.Printf("Error writing result: %v\n", err)
	}
}

func runDependenciesScan(cmd *cobra.Command, args []string) {
	// Get flags
	depth, _ := cmd.Flags().GetString("depth")
	minSeverity, _ := cmd.Flags().GetString("severity")
	format, _ := cmd.Flags().GetString("format")
	outputFile, _ := cmd.Flags().GetString("output")
	showTree, _ := cmd.Flags().GetBool("tree")
	directOnly, _ := cmd.Flags().GetBool("direct-only")

	// Determine repository path
	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}

	// Check if it's a Git repository
	if !isGitRepository(repoPath) {
		fmt.Printf("Error: %s is not a Git repository\n", repoPath)
		os.Exit(1)
	}

	fmt.Printf("🔍 Scanning transitive dependencies for vulnerabilities...\n")
	fmt.Printf("Repository: %s\n", repoPath)
	fmt.Printf("Scan depth: %s\n", depth)
	fmt.Printf("Minimum severity: %s\n", minSeverity)
	fmt.Printf("Direct dependencies only: %v\n", directOnly)
	fmt.Printf("Show dependency tree: %v\n\n", showTree)

	// Run transitive dependency checker
	depChecker := checkers.NewTransitiveDependencyChecker()

	// Create RepositoryData for the checker
	analyzer, err := git.NewRepositoryAnalyzer(repoPath)
	if err != nil {
		fmt.Printf("Error initializing repository analyzer: %v\n", err)
		os.Exit(1)
	}
	data, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("Error analyzing repository: %v\n", err)
		os.Exit(1)
	}

	inventoryResult, dependencyTree := depChecker.CheckWithOptions(data, directOnly, depth)
	result := depChecker.ScanVulnerabilities(data, minSeverity)
	result.Details = append(inventoryResult.Details, result.Details...)

	// Display results based on format
	switch format {
	case "json":
		outputDependenciesJSON(result, dependencyTree, outputFile)
	case "yaml":
		outputDependenciesYAML(result, dependencyTree, outputFile)
	default:
		outputDependenciesTable(result, dependencyTree, showTree, minSeverity)
	}

	// Show remediation if vulnerabilities found
	if result.Status == types.StatusFail {
		fmt.Printf("\n🚨 VULNERABILITIES FOUND!\n\n")
		fmt.Printf("Immediate Actions Required:\n")
		fmt.Printf("1. Update vulnerable dependencies to secure versions\n")
		fmt.Printf("2. Review dependency tree to identify root causes\n")
		fmt.Printf("3. Consider removing unnecessary dependencies\n")
		fmt.Printf("4. Implement dependency scanning in CI/CD pipeline\n\n")

		fmt.Printf("Tools for Dependency Management:\n")
		fmt.Printf("- npm audit fix (Node.js)\n")
		fmt.Printf("- go get -u (Go)\n")
		fmt.Printf("- pip install --upgrade (Python)\n")
		fmt.Printf("- cargo update (Rust)\n")
		fmt.Printf("- mvn versions:use-latest-releases (Java)\n\n")

		fmt.Printf("Prevention:\n")
		fmt.Printf("- Use dependency scanning tools in CI/CD\n")
		fmt.Printf("- Regularly update dependencies\n")
		fmt.Printf("- Use lock files (package-lock.json, go.sum, etc.)\n")
		fmt.Printf("- Monitor security advisories\n")

		os.Exit(1)
	} else if result.Status == types.StatusWarning {
		fmt.Printf("⚠️  %s\n", result.Message)
	} else {
		fmt.Printf("✅ Dependency scan completed.\n")
	}
}

// outputDependenciesTable outputs dependency scan results in table format
func outputDependenciesTable(result *types.CheckResult, tree *checkers.DependencyTree, showTree bool, minSeverity string) {
	fmt.Printf("📊 Dependency Scan Results\n")
	fmt.Printf("==========================\n\n")

	// Display details
	for _, detail := range result.Details {
		fmt.Printf("%s\n", detail)
	}

	fmt.Printf("Security Score: %d/100\n\n", result.Score)

	// Note: Dependency tree display would need to be implemented differently
	// since Details is now []string instead of map[string]interface{}
	if showTree {
		fmt.Printf("🌳 Dependency Tree\n")
		fmt.Printf("==================\n\n")
		if tree == nil || tree.Root == nil {
			fmt.Printf("No dependency tree available.\n")
		} else {
			printDependencyTree(tree.Root, 0)
		}
	}
}

// printDependencyTree recursively prints the dependency tree
func printDependencyTree(dep *checkers.Dependency, depth int) {
	indent := strings.Repeat("  ", depth)

	// Determine icon based on vulnerability status
	icon := "📦"
	if dep.Vulnerable {
		switch dep.Severity {
		case "critical":
			icon = "🚨"
		case "high":
			icon = "⚠️"
		case "medium":
			icon = "🔶"
		case "low":
			icon = "🔸"
		}
	}

	// Print dependency info
	fmt.Printf("%s%s %s@%s", indent, icon, dep.Name, dep.Version)
	if dep.Direct {
		fmt.Printf(" (direct)")
	}
	if dep.Vulnerable {
		fmt.Printf(" [%s]", strings.ToUpper(dep.Severity))
	}
	fmt.Printf("\n")

	// Print vulnerabilities
	for _, vuln := range dep.Vulnerabilities {
		fmt.Printf("%s  🔍 %s: %s (CVSS: %.1f)\n", indent, vuln.ID, vuln.Description, vuln.CVSS)
	}

	// Recursively print children
	for _, child := range dep.Children {
		printDependencyTree(child, depth+1)
	}
}

// outputDependenciesJSON outputs dependency scan results in JSON format
func outputDependenciesJSON(result *types.CheckResult, tree *checkers.DependencyTree, outputFile string) {
	payload := struct {
		Result *types.CheckResult       `json:"result"`
		Tree   *checkers.DependencyTree `json:"tree,omitempty"`
	}{Result: result, Tree: tree}
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			fmt.Printf("Error writing JSON file: %v\n", err)
			return
		}
		fmt.Printf("Results written to %s\n", outputFile)
	} else {
		fmt.Printf("%s\n", string(jsonData))
	}
}

// outputDependenciesYAML outputs dependency scan results in YAML format
func outputDependenciesYAML(result *types.CheckResult, tree *checkers.DependencyTree, outputFile string) {
	payload := struct {
		Result *types.CheckResult       `yaml:"result"`
		Tree   *checkers.DependencyTree `yaml:"tree,omitempty"`
	}{Result: result, Tree: tree}
	yamlData, err := yaml.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling YAML: %v\n", err)
		return
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, yamlData, 0644)
		if err != nil {
			fmt.Printf("Error writing YAML file: %v\n", err)
			return
		}
		fmt.Printf("Results written to %s\n", outputFile)
	} else {
		fmt.Printf("%s\n", string(yamlData))
	}
}

func runPolicyValidation(cmd *cobra.Command, args []string) {
	// Get flags
	checkSigning, _ := cmd.Flags().GetBool("check-signing")
	checkFiles, _ := cmd.Flags().GetBool("check-files")
	checkPush, _ := cmd.Flags().GetBool("check-push")
	checkBranches, _ := cmd.Flags().GetBool("check-branches")
	minSeverity, _ := cmd.Flags().GetString("severity")
	format, _ := cmd.Flags().GetString("format")
	outputFile, _ := cmd.Flags().GetString("output")

	// Determine repository path
	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}

	// Check if it's a Git repository
	if !isGitRepository(repoPath) {
		fmt.Printf("Error: %s is not a Git repository\n", repoPath)
		os.Exit(1)
	}

	fmt.Printf("🔍 Validating Git security policies...\n")
	fmt.Printf("Repository: %s\n", repoPath)
	fmt.Printf("Check signing: %v\n", checkSigning)
	fmt.Printf("Check files: %v\n", checkFiles)
	fmt.Printf("Check push policies: %v\n", checkPush)
	fmt.Printf("Check branch protection: %v\n", checkBranches)
	fmt.Printf("Minimum severity: %s\n\n", minSeverity)

	// Run Git policy checker
	policyChecker := checkers.NewGitPolicyChecker()

	// Create RepositoryData for the checker
	analyzer, err := git.NewRepositoryAnalyzer(repoPath)
	if err != nil {
		fmt.Printf("Error initializing repository analyzer: %v\n", err)
		os.Exit(1)
	}
	data, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("Error analyzing repository: %v\n", err)
		os.Exit(1)
	}

	result := policyChecker.CheckWithOptions(data, checkSigning, checkFiles, checkPush, checkBranches, minSeverity)

	// Process results
	if result.Status == types.StatusFail {
		fmt.Printf("❌ Policy validation found issues: %s\n", result.Message)
	} else {
		fmt.Printf("✅ Policy validation passed: %s\n", result.Message)
	}

	// Display results based on format
	switch format {
	case "json":
		outputPolicyJSON(result, outputFile)
	case "yaml":
		outputPolicyYAML(result, outputFile)
	default:
		outputPolicyTable(result, minSeverity)
	}

	// Show remediation if violations found
	if result.Status == types.StatusFail {
		fmt.Printf("\n🚨 POLICY VIOLATIONS FOUND!\n\n")
		fmt.Printf("Immediate Actions Required:\n")
		fmt.Printf("1. Review and fix policy violations\n")
		fmt.Printf("2. Enable commit signing for important commits\n")
		fmt.Printf("3. Add sensitive files to .gitignore\n")
		fmt.Printf("4. Configure branch protection rules\n")
		fmt.Printf("5. Review push policies and permissions\n\n")

		fmt.Printf("Security Best Practices:\n")
		fmt.Printf("- Enable GPG commit signing\n")
		fmt.Printf("- Use .gitignore for sensitive files\n")
		fmt.Printf("- Configure branch protection\n")
		fmt.Printf("- Use signed commits for releases\n")
		fmt.Printf("- Regular security policy audits\n")

		os.Exit(1)
	} else {
		fmt.Printf("✅ All Git security policies are properly configured!\n")
	}
}

// outputPolicyTable outputs policy validation results in table format
func outputPolicyTable(result *types.CheckResult, minSeverity string) {
	fmt.Printf("📊 Git Policy Validation Results\n")
	fmt.Printf("=================================\n\n")

	// Display details
	for _, detail := range result.Details {
		fmt.Printf("%s\n", detail)
	}

	fmt.Printf("Security Score: %d/100\n\n", result.Score)
}

// outputPolicyJSON outputs policy validation results in JSON format
func outputPolicyJSON(result *types.CheckResult, outputFile string) {
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			fmt.Printf("Error writing JSON file: %v\n", err)
			return
		}
		fmt.Printf("Results written to %s\n", outputFile)
	} else {
		fmt.Printf("%s\n", string(jsonData))
	}
}

// outputPolicyYAML outputs policy validation results in YAML format
func outputPolicyYAML(result *types.CheckResult, outputFile string) {
	yamlData, err := yaml.Marshal(result)
	if err != nil {
		fmt.Printf("Error marshaling YAML: %v\n", err)
		return
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, yamlData, 0644)
		if err != nil {
			fmt.Printf("Error writing YAML file: %v\n", err)
			return
		}
		fmt.Printf("Results written to %s\n", outputFile)
	} else {
		fmt.Printf("%s\n", string(yamlData))
	}
}

func runBinariesAudit(cmd *cobra.Command, args []string) {
	// Get flags
	maxSizeStr, _ := cmd.Flags().GetString("max-size")
	checkHistory, _ := cmd.Flags().GetBool("check-history")
	checkExecutables, _ := cmd.Flags().GetBool("check-executables")
	checkLarge, _ := cmd.Flags().GetBool("check-large")
	checkSuspicious, _ := cmd.Flags().GetBool("check-suspicious")
	minSeverity, _ := cmd.Flags().GetString("severity")
	format, _ := cmd.Flags().GetString("format")
	outputFile, _ := cmd.Flags().GetString("output")

	// Parse max size
	maxSizeMB := parseSizeToMB(maxSizeStr)

	// Determine repository path
	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}

	// Check if it's a Git repository
	if !isGitRepository(repoPath) {
		fmt.Printf("Error: %s is not a Git repository\n", repoPath)
		os.Exit(1)
	}

	fmt.Printf("🔍 Auditing executable and large files...\n")
	fmt.Printf("Repository: %s\n", repoPath)
	fmt.Printf("Max size threshold: %s (%.1f MB)\n", maxSizeStr, maxSizeMB)
	fmt.Printf("Check history: %v\n", checkHistory)
	fmt.Printf("Check executables: %v\n", checkExecutables)
	fmt.Printf("Check large files: %v\n", checkLarge)
	fmt.Printf("Check suspicious files: %v\n", checkSuspicious)
	fmt.Printf("Minimum severity: %s\n\n", minSeverity)

	// Run binary file checker
	binaryChecker := checkers.NewBinaryFileChecker()

	// Create RepositoryData for the checker
	analyzer, err := git.NewRepositoryAnalyzer(repoPath)
	if err != nil {
		fmt.Printf("Error initializing repository analyzer: %v\n", err)
		os.Exit(1)
	}
	data, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("Error analyzing repository: %v\n", err)
		os.Exit(1)
	}

	result := binaryChecker.CheckWithSeverity(data, checkExecutables, checkLarge, checkSuspicious, checkHistory, maxSizeMB, minSeverity)

	// Process results
	if result.Status == types.StatusFail {
		fmt.Printf("❌ Binary audit found issues: %s\n", result.Message)
	} else {
		fmt.Printf("✅ Binary audit passed: %s\n", result.Message)
	}

	// Display results based on format
	switch format {
	case "json":
		outputBinariesJSON(result, outputFile)
	case "yaml":
		outputBinariesYAML(result, outputFile)
	default:
		outputBinariesTable(result, minSeverity)
	}

	// Show remediation if violations found
	if result.Status == types.StatusFail {
		fmt.Printf("\n🚨 BINARY FILE ISSUES FOUND!\n\n")
		fmt.Printf("Immediate Actions Required:\n")
		fmt.Printf("1. Review and remove unnecessary binary files\n")
		fmt.Printf("2. Use Git LFS for large files\n")
		fmt.Printf("3. Add binary file patterns to .gitignore\n")
		fmt.Printf("4. Remove suspicious files from repository\n")
		fmt.Printf("5. Clean up Git history if needed\n\n")

		fmt.Printf("Best Practices:\n")
		fmt.Printf("- Use Git LFS for files > 100MB\n")
		fmt.Printf("- Avoid committing executable files\n")
		fmt.Printf("- Use .gitignore for binary patterns\n")
		fmt.Printf("- Regular binary file audits\n")
		fmt.Printf("- Use package managers for dependencies\n")

		os.Exit(1)
	} else {
		fmt.Printf("✅ No suspicious binary or large files found!\n")
	}
}

// parseSizeToMB parses size string to MB
func parseSizeToMB(sizeStr string) float64 {
	sizeStr = strings.ToLower(strings.TrimSpace(sizeStr))

	// Remove common suffixes
	sizeStr = strings.TrimSuffix(sizeStr, "mb")
	sizeStr = strings.TrimSuffix(sizeStr, "m")

	// Parse number
	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 10.0 // Default 10MB
	}

	return size
}

// outputBinariesTable outputs binary audit results in table format
func outputBinariesTable(result *types.CheckResult, minSeverity string) {
	fmt.Printf("📊 Binary File Audit Results\n")
	fmt.Printf("============================\n\n")

	// Display all details including file listings
	for _, detail := range result.Details {
		fmt.Printf("%s\n", detail)
	}

	fmt.Printf("\nSecurity Score: %d/100\n\n", result.Score)
}

// outputBinariesJSON outputs binary audit results in JSON format
func outputBinariesJSON(result *types.CheckResult, outputFile string) {
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			fmt.Printf("Error writing JSON file: %v\n", err)
			return
		}
		fmt.Printf("Results written to %s\n", outputFile)
	} else {
		fmt.Printf("%s\n", string(jsonData))
	}
}

// outputBinariesYAML outputs binary audit results in YAML format
func outputBinariesYAML(result *types.CheckResult, outputFile string) {
	yamlData, err := yaml.Marshal(result)
	if err != nil {
		fmt.Printf("Error marshaling YAML: %v\n", err)
		return
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, yamlData, 0644)
		if err != nil {
			fmt.Printf("Error writing YAML file: %v\n", err)
			return
		}
		fmt.Printf("Results written to %s\n", outputFile)
	} else {
		fmt.Printf("%s\n", string(yamlData))
	}
}
