package checkers

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// GitPolicyChecker validates Git security policies and configurations
type GitPolicyChecker struct {
	BaseChecker
}

// PolicyViolation represents a security policy violation
type PolicyViolation struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	File        string `json:"file,omitempty"`
	Line        int    `json:"line,omitempty"`
	Recommendation string `json:"recommendation"`
}

// SignatureStats represents commit signature statistics
type SignatureStats struct {
	TotalCommits    int     `json:"total_commits"`
	SignedCommits   int     `json:"signed_commits"`
	UnsignedCommits int     `json:"unsigned_commits"`
	SignatureRate   float64 `json:"signature_rate"`
	ValidSignatures int     `json:"valid_signatures"`
	InvalidSignatures int   `json:"invalid_signatures"`
}

// SensitiveFile represents a detected sensitive file
type SensitiveFile struct {
	Path        string `json:"path"`
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	InGitignore bool   `json:"in_gitignore"`
	InHistory   bool   `json:"in_history"`
}

// GitPolicyReport represents the complete policy validation report
type GitPolicyReport struct {
	Violations      []PolicyViolation `json:"violations"`
	SignatureStats  SignatureStats    `json:"signature_stats"`
	SensitiveFiles  []SensitiveFile   `json:"sensitive_files"`
	PushPolicies    []string          `json:"push_policies"`
	BranchProtection []string          `json:"branch_protection"`
	Score           int               `json:"score"`
}

// NewGitPolicyChecker creates a new GitPolicyChecker
func NewGitPolicyChecker() *GitPolicyChecker {
	return &GitPolicyChecker{
		BaseChecker: BaseChecker{
			id:   "GIT-POLICY",
			name: "Git Policy Validation",
		},
	}
}

// Check performs Git security policy validation
func (c *GitPolicyChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        c.ID(),
		Name:      c.Name(),
		Status:    types.StatusPass,
		Score:     100,
		Message:   "All Git security policies are properly configured",
		Details:   []string{},
		Category:  c.Category(),
		Timestamp: time.Now(),
	}

	// Initialize policy report
	report := &GitPolicyReport{
		Violations:      []PolicyViolation{},
		SignatureStats:  SignatureStats{},
		SensitiveFiles:  []SensitiveFile{},
		PushPolicies:    []string{},
		BranchProtection: []string{},
	}

	// Check Git configuration
	c.checkGitConfig(data.Path, report)

	// Check commit signatures
	c.checkCommitSignatures(data.Path, report)

	// Check sensitive files
	c.checkSensitiveFiles(data.Path, report)

	// Check push policies
	c.checkPushPolicies(data.Path, report)

	// Check branch protection
	c.checkBranchProtection(data.Path, report)

	// Calculate score based on violations
	score := c.calculateScore(report)
	result.Score = score

	// Update result based on findings
	if len(report.Violations) > 0 {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Found %d Git security policy violations", len(report.Violations))
	} else {
		result.Status = types.StatusPass
		result.Message = "All Git security policies are properly configured"
	}

	// Add detailed information
	result.Details = append(result.Details, fmt.Sprintf("Total Violations: %d", len(report.Violations)))
	result.Details = append(result.Details, fmt.Sprintf("Signature Rate: %.1f%%", report.SignatureStats.SignatureRate))
	result.Details = append(result.Details, fmt.Sprintf("Sensitive Files: %d", len(report.SensitiveFiles)))
	result.Details = append(result.Details, fmt.Sprintf("Push Policies: %d", len(report.PushPolicies)))
	result.Details = append(result.Details, fmt.Sprintf("Branch Protection: %d", len(report.BranchProtection)))

	return result
}

// checkGitConfig checks Git configuration for security issues
func (c *GitPolicyChecker) checkGitConfig(repoPath string, report *GitPolicyReport) {
	configPath := filepath.Join(repoPath, ".git", "config")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return
	}

	file, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Check for dangerous push settings
		if strings.Contains(line, "push.default") && strings.Contains(line, "matching") {
			report.Violations = append(report.Violations, PolicyViolation{
				Type:        "push_policy",
				Severity:    "medium",
				Description: "Push default is set to 'matching' which can be dangerous",
				File:        configPath,
				Line:        lineNum,
				Recommendation: "Set push.default to 'simple' or 'current'",
			})
		}

		// Check for unsafe credential storage
		if strings.Contains(line, "credential.helper") && strings.Contains(line, "store") {
			report.Violations = append(report.Violations, PolicyViolation{
				Type:        "credential_storage",
				Severity:    "high",
				Description: "Credentials are stored in plain text",
				File:        configPath,
				Line:        lineNum,
				Recommendation: "Use credential.helper=cache or credential.helper=osxkeychain",
			})
		}

		// Check for unsafe merge settings
		if strings.Contains(line, "merge.ours") && strings.Contains(line, "true") {
			report.Violations = append(report.Violations, PolicyViolation{
				Type:        "merge_policy",
				Severity:    "medium",
				Description: "Merge strategy 'ours' can hide conflicts",
				File:        configPath,
				Line:        lineNum,
				Recommendation: "Use merge strategy 'recursive' or 'resolve'",
			})
		}
	}
}

// checkCommitSignatures checks commit signature verification
func (c *GitPolicyChecker) checkCommitSignatures(repoPath string, report *GitPolicyReport) {
	// Get commit signature statistics
	cmd := exec.Command("git", "log", "--pretty=format:%G?%n", "--all")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(output), "\n")
	totalCommits := 0
	signedCommits := 0
	validSignatures := 0
	invalidSignatures := 0

	for _, line := range lines {
		if line == "" {
			continue
		}
		totalCommits++

		switch line {
		case "G": // Good signature
			signedCommits++
			validSignatures++
		case "B": // Bad signature
			signedCommits++
			invalidSignatures++
		case "U": // Good signature with unknown validity
			signedCommits++
		case "X": // Good signature that expired
			signedCommits++
		case "Y": // Good signature made by an expired key
			signedCommits++
		case "R": // Good signature made by a revoked key
			signedCommits++
		case "E": // Signature could not be checked
			signedCommits++
		case "N": // No signature
			// Unsigned commit
		}
	}

	report.SignatureStats = SignatureStats{
		TotalCommits:     totalCommits,
		SignedCommits:    signedCommits,
		UnsignedCommits:  totalCommits - signedCommits,
		SignatureRate:    float64(signedCommits) / float64(totalCommits) * 100,
		ValidSignatures:   validSignatures,
		InvalidSignatures: invalidSignatures,
	}

	// Add violations based on signature rate
	if report.SignatureStats.SignatureRate < 50.0 {
		report.Violations = append(report.Violations, PolicyViolation{
			Type:        "signature_policy",
			Severity:    "high",
			Description: fmt.Sprintf("Low signature rate: %.1f%%", report.SignatureStats.SignatureRate),
			Recommendation: "Enable commit signing and require signatures for important commits",
		})
	} else if report.SignatureStats.SignatureRate < 80.0 {
		report.Violations = append(report.Violations, PolicyViolation{
			Type:        "signature_policy",
			Severity:    "medium",
			Description: fmt.Sprintf("Moderate signature rate: %.1f%%", report.SignatureStats.SignatureRate),
			Recommendation: "Consider enabling commit signing for more commits",
		})
	}

	// Check for invalid signatures
	if invalidSignatures > 0 {
		report.Violations = append(report.Violations, PolicyViolation{
			Type:        "signature_validation",
			Severity:    "high",
			Description: fmt.Sprintf("Found %d invalid signatures", invalidSignatures),
			Recommendation: "Review and fix invalid signatures",
		})
	}
}

// checkSensitiveFiles checks for sensitive files in repository
func (c *GitPolicyChecker) checkSensitiveFiles(repoPath string, report *GitPolicyReport) {
	sensitivePatterns := map[string]SensitiveFile{
		".env": {
			Type:        "environment",
			Severity:    "high",
			Description: "Environment variables file",
		},
		".env.local": {
			Type:        "environment",
			Severity:    "high",
			Description: "Local environment variables file",
		},
		".env.production": {
			Type:        "environment",
			Severity:    "critical",
			Description: "Production environment variables file",
		},
		"kubeconfig": {
			Type:        "kubernetes",
			Severity:    "critical",
			Description: "Kubernetes configuration file",
		},
		".kube/config": {
			Type:        "kubernetes",
			Severity:    "critical",
			Description: "Kubernetes configuration file",
		},
		"id_rsa": {
			Type:        "ssh_key",
			Severity:    "critical",
			Description: "SSH private key",
		},
		"id_dsa": {
			Type:        "ssh_key",
			Severity:    "critical",
			Description: "SSH private key",
		},
		"id_ed25519": {
			Type:        "ssh_key",
			Severity:    "critical",
			Description: "SSH private key",
		},
		"*.pem": {
			Type:        "certificate",
			Severity:    "high",
			Description: "Certificate file",
		},
		"*.key": {
			Type:        "private_key",
			Severity:    "critical",
			Description: "Private key file",
		},
		"*.p12": {
			Type:        "certificate",
			Severity:    "high",
			Description: "PKCS#12 certificate file",
		},
		"*.pfx": {
			Type:        "certificate",
			Severity:    "high",
			Description: "PKCS#12 certificate file",
		},
		"config.json": {
			Type:        "configuration",
			Severity:    "medium",
			Description: "Configuration file",
		},
		"secrets.json": {
			Type:        "secrets",
			Severity:    "critical",
			Description: "Secrets file",
		},
		"credentials.json": {
			Type:        "credentials",
			Severity:    "critical",
			Description: "Credentials file",
		},
	}

	// Check current working directory
	c.scanDirectoryForSensitiveFiles(repoPath, sensitivePatterns, report)

	// Check Git history for sensitive files
	c.checkGitHistoryForSensitiveFiles(repoPath, sensitivePatterns, report)

	// Check .gitignore for sensitive file patterns
	c.checkGitignoreForSensitiveFiles(repoPath, sensitivePatterns, report)
}

// scanDirectoryForSensitiveFiles scans current directory for sensitive files
func (c *GitPolicyChecker) scanDirectoryForSensitiveFiles(repoPath string, patterns map[string]SensitiveFile, report *GitPolicyReport) {
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip .git directory
		if strings.Contains(path, ".git") {
			return nil
		}

		// Check if file matches sensitive patterns
		relPath, _ := filepath.Rel(repoPath, path)
		for pattern, fileInfo := range patterns {
			if matched, _ := filepath.Match(pattern, info.Name()); matched {
				sensitiveFile := fileInfo
				sensitiveFile.Path = relPath
				sensitiveFile.InGitignore = c.isInGitignore(repoPath, relPath)
				sensitiveFile.InHistory = false // Will be checked separately

				report.SensitiveFiles = append(report.SensitiveFiles, sensitiveFile)

				// Add violation if not in .gitignore
				if !sensitiveFile.InGitignore {
					report.Violations = append(report.Violations, PolicyViolation{
						Type:        "sensitive_file",
						Severity:    sensitiveFile.Severity,
						Description: fmt.Sprintf("Sensitive file found: %s", relPath),
						File:        relPath,
						Recommendation: fmt.Sprintf("Add %s to .gitignore", pattern),
					})
				}
			}
		}

		return nil
	})

	if err != nil {
		// Handle error silently
	}
}

// checkGitHistoryForSensitiveFiles checks Git history for sensitive files
func (c *GitPolicyChecker) checkGitHistoryForSensitiveFiles(repoPath string, patterns map[string]SensitiveFile, report *GitPolicyReport) {
	// Get list of all files that have ever been tracked
	cmd := exec.Command("git", "log", "--name-only", "--pretty=format:", "--all")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return
	}

	files := strings.Split(string(output), "\n")
	seenFiles := make(map[string]bool)

	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" || seenFiles[file] {
			continue
		}
		seenFiles[file] = true

		// Check if file matches sensitive patterns
		for pattern, fileInfo := range patterns {
			if matched, _ := filepath.Match(pattern, filepath.Base(file)); matched {
				sensitiveFile := fileInfo
				sensitiveFile.Path = file
				sensitiveFile.InGitignore = c.isInGitignore(repoPath, file)
				sensitiveFile.InHistory = true

				// Check if this file is already in our list
				found := false
				for i, existingFile := range report.SensitiveFiles {
					if existingFile.Path == file {
						report.SensitiveFiles[i].InHistory = true
						found = true
						break
					}
				}

				if !found {
					report.SensitiveFiles = append(report.SensitiveFiles, sensitiveFile)
				}

				// Add violation for files in history
				report.Violations = append(report.Violations, PolicyViolation{
					Type:        "sensitive_file_history",
					Severity:    sensitiveFile.Severity,
					Description: fmt.Sprintf("Sensitive file found in Git history: %s", file),
					File:        file,
					Recommendation: "Remove from Git history using git filter-repo or BFG",
				})
			}
		}
	}
}

// checkGitignoreForSensitiveFiles checks if sensitive files are properly ignored
func (c *GitPolicyChecker) checkGitignoreForSensitiveFiles(repoPath string, patterns map[string]SensitiveFile, report *GitPolicyReport) {
	gitignorePath := filepath.Join(repoPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		report.Violations = append(report.Violations, PolicyViolation{
			Type:        "gitignore_missing",
			Severity:    "medium",
			Description: ".gitignore file is missing",
			Recommendation: "Create .gitignore file with sensitive file patterns",
		})
		return
	}

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		return
	}

	gitignoreContent := string(content)
	missingPatterns := []string{}

	for pattern := range patterns {
		if !strings.Contains(gitignoreContent, pattern) {
			missingPatterns = append(missingPatterns, pattern)
		}
	}

	if len(missingPatterns) > 0 {
		report.Violations = append(report.Violations, PolicyViolation{
			Type:        "gitignore_incomplete",
			Severity:    "medium",
			Description: fmt.Sprintf("Missing patterns in .gitignore: %s", strings.Join(missingPatterns, ", ")),
			File:        ".gitignore",
			Recommendation: "Add missing patterns to .gitignore",
		})
	}
}

// isInGitignore checks if a file is ignored by .gitignore
func (c *GitPolicyChecker) isInGitignore(repoPath, filePath string) bool {
	cmd := exec.Command("git", "check-ignore", filePath)
	cmd.Dir = repoPath
	err := cmd.Run()
	return err == nil
}

// checkPushPolicies checks Git push policies
func (c *GitPolicyChecker) checkPushPolicies(repoPath string, report *GitPolicyReport) {
	// Check for force push settings
	cmd := exec.Command("git", "config", "--get", "push.default")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err == nil {
		pushDefault := strings.TrimSpace(string(output))
		report.PushPolicies = append(report.PushPolicies, fmt.Sprintf("push.default: %s", pushDefault))

		if pushDefault == "matching" {
			report.Violations = append(report.Violations, PolicyViolation{
				Type:        "push_policy",
				Severity:    "medium",
				Description: "Push default is set to 'matching'",
				Recommendation: "Set push.default to 'simple' for safer pushes",
			})
		}
	}

	// Check for force push permissions
	cmd = exec.Command("git", "config", "--get", "receive.denyNonFastForwards")
	cmd.Dir = repoPath
	output, err = cmd.Output()
	if err == nil {
		denyNonFastForwards := strings.TrimSpace(string(output))
		if denyNonFastForwards != "true" {
			report.Violations = append(report.Violations, PolicyViolation{
				Type:        "push_policy",
				Severity:    "high",
				Description: "Force pushes are not denied",
				Recommendation: "Set receive.denyNonFastForwards=true to prevent force pushes",
			})
		}
	}
}

// checkBranchProtection checks branch protection settings
func (c *GitPolicyChecker) checkBranchProtection(repoPath string, report *GitPolicyReport) {
	// Get list of branches
	cmd := exec.Command("git", "branch", "-r")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return
	}

	branches := strings.Split(string(output), "\n")
	protectedBranches := []string{"main", "master", "develop", "production"}

	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch == "" {
			continue
		}

		// Remove origin/ prefix
		branchName := strings.TrimPrefix(branch, "origin/")

		// Check if it's a protected branch
		for _, protected := range protectedBranches {
			if branchName == protected {
				report.BranchProtection = append(report.BranchProtection, branchName)

				// Check if branch has protection rules
				cmd := exec.Command("git", "config", "--get", fmt.Sprintf("branch.%s.protection", branchName))
				cmd.Dir = repoPath
				_, err := cmd.Output()
				if err != nil {
					report.Violations = append(report.Violations, PolicyViolation{
						Type:        "branch_protection",
						Severity:    "high",
						Description: fmt.Sprintf("Protected branch '%s' has no protection rules", branchName),
						Recommendation: "Configure branch protection rules for important branches",
					})
				}
				break
			}
		}
	}
}

// calculateScore calculates security score based on violations
func (c *GitPolicyChecker) calculateScore(report *GitPolicyReport) int {
	score := 100

	// Deduct points based on violation severity
	for _, violation := range report.Violations {
		switch violation.Severity {
		case "critical":
			score -= 20
		case "high":
			score -= 15
		case "medium":
			score -= 10
		case "low":
			score -= 5
		}
	}

	// Deduct points for low signature rate
	if report.SignatureStats.SignatureRate < 50.0 {
		score -= 15
	} else if report.SignatureStats.SignatureRate < 80.0 {
		score -= 10
	}

	// Deduct points for sensitive files
	for _, file := range report.SensitiveFiles {
		if !file.InGitignore {
			switch file.Severity {
			case "critical":
				score -= 15
			case "high":
				score -= 10
			case "medium":
				score -= 5
			}
		}
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

// GetCategory returns the category of this checker
func (c *GitPolicyChecker) GetCategory() string {
	return "Security"
}
