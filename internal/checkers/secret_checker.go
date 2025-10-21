package checkers

import (
	"bufio"
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// SecretChecker checks for secrets in Git history
type SecretChecker struct {
	BaseChecker
}

// Secret represents a found secret
type Secret struct {
	Type        string    `json:"type"`
	Pattern     string    `json:"pattern"`
	File        string    `json:"file"`
	Commit      string    `json:"commit"`
	Line        int       `json:"line"`
	Content     string    `json:"content"`
	Severity    string    `json:"severity"`
	Confidence  float64   `json:"confidence"`
	Timestamp   time.Time `json:"timestamp"`
	Remediation string    `json:"remediation"`
}

// SecretPattern represents a regex pattern for secret detection
type SecretPattern struct {
	Name        string
	Pattern     *regexp.Regexp
	Severity    string
	Confidence  float64
	Description string
	Remediation string
}

// NewSecretChecker creates a new SecretChecker
func NewSecretChecker() *SecretChecker {
	return &SecretChecker{
		BaseChecker: NewBaseChecker(
			"Secret Scanning",
			"secret-scanning",
			types.CategoryStructure,
			25,
		),
	}
}

// Check performs secret scanning
func (c *SecretChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "secret-scanning",
		Name:      c.Name(),
		Status:    types.StatusPass,
		Score:     100,
		Message:   "No secrets found in Git history",
		Details:   []string{},
		Category:  types.CategoryStructure,
		Timestamp: time.Now(),
	}

	secrets, err := c.scanGitHistory(data.Path)
	if err != nil {
		result.Status = types.StatusFail
		result.Score = 0
		result.Message = fmt.Sprintf("Error scanning for secrets: %v", err)
		result.Details = []string{err.Error()}
		return result
	}

	if len(secrets) > 0 {
		result.Status = types.StatusFail
		result.Score = 0
		result.Message = fmt.Sprintf("Found %d secrets in Git history", len(secrets))
		
		// Add details about found secrets
		details := []string{
			fmt.Sprintf("Total secrets found: %d", len(secrets)),
			fmt.Sprintf("High severity secrets: %d", c.countHighSeveritySecrets(secrets)),
		}
		
		// Add individual secret details
		for i, secret := range secrets {
			details = append(details, fmt.Sprintf("%d. %s (%s) in %s:%d", 
				i+1, secret.Type, secret.Severity, secret.File, secret.Line))
		}
		
		result.Details = details
	} else {
		result.Details = []string{"No secrets found in Git history"}
	}

	return result
}

// scanGitHistory scans the entire Git history for secrets
func (c *SecretChecker) scanGitHistory(repoPath string) ([]Secret, error) {
	var secrets []Secret

	// Get all commits
	commits, err := c.getAllCommits(repoPath)
	if err != nil {
		return nil, err
	}

	// Get all stashes
	stashes, err := c.getAllStashes(repoPath)
	if err != nil {
		return nil, err
	}

	// Scan commits
	for _, commit := range commits {
		commitSecrets, err := c.scanCommit(repoPath, commit)
		if err != nil {
			continue // Skip failed commits
		}
		secrets = append(secrets, commitSecrets...)
	}

	// Scan stashes
	for _, stash := range stashes {
		stashSecrets, err := c.scanStash(repoPath, stash)
		if err != nil {
			continue // Skip failed stashes
		}
		secrets = append(secrets, stashSecrets...)
	}

	return secrets, nil
}

// getAllCommits gets all commit hashes
func (c *SecretChecker) getAllCommits(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "rev-list", "--all")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var commits []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		commits = append(commits, strings.TrimSpace(scanner.Text()))
	}

	return commits, nil
}

// getAllStashes gets all stash references
func (c *SecretChecker) getAllStashes(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "stash", "list", "--format=%H")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var stashes []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		stashes = append(stashes, strings.TrimSpace(scanner.Text()))
	}

	return stashes, nil
}

// scanCommit scans a specific commit for secrets
func (c *SecretChecker) scanCommit(repoPath string, commitHash string) ([]Secret, error) {
	var secrets []Secret

	// Get commit files
	cmd := exec.Command("git", "show", "--name-only", "--pretty=format:", commitHash)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(output), "\n")
	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}

		// Get file content for this commit
		fileSecrets, err := c.scanFileInCommit(repoPath, commitHash, file)
		if err != nil {
			continue
		}

		secrets = append(secrets, fileSecrets...)
	}

	return secrets, nil
}

// scanStash scans a specific stash for secrets
func (c *SecretChecker) scanStash(repoPath string, stashHash string) ([]Secret, error) {
	var secrets []Secret

	// Get stash files
	cmd := exec.Command("git", "stash", "show", "--name-only", stashHash)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(output), "\n")
	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}

		// Get file content for this stash
		fileSecrets, err := c.scanFileInStash(repoPath, stashHash, file)
		if err != nil {
			continue
		}

		secrets = append(secrets, fileSecrets...)
	}

	return secrets, nil
}

// scanFileInCommit scans a file in a specific commit
func (c *SecretChecker) scanFileInCommit(repoPath string, commitHash, filePath string) ([]Secret, error) {
	cmd := exec.Command("git", "show", commitHash+":"+filePath)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return c.scanContent(string(output), filePath, commitHash, "commit"), nil
}

// scanFileInStash scans a file in a specific stash
func (c *SecretChecker) scanFileInStash(repoPath string, stashHash, filePath string) ([]Secret, error) {
	cmd := exec.Command("git", "stash", "show", "-p", stashHash, "--", filePath)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return c.scanContent(string(output), filePath, stashHash, "stash"), nil
}

// scanContent scans content for secrets
func (c *SecretChecker) scanContent(content, filePath, ref, refType string) []Secret {
	var secrets []Secret
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		// Check against known patterns
		patternSecrets := c.checkPatterns(line, filePath, ref, refType, lineNum+1)
		secrets = append(secrets, patternSecrets...)

		// Check entropy for random-looking strings
		entropySecrets := c.checkEntropy(line, filePath, ref, refType, lineNum+1)
		secrets = append(secrets, entropySecrets...)
	}

	return secrets
}

// checkPatterns checks content against known secret patterns
func (c *SecretChecker) checkPatterns(line, filePath, ref, refType string, lineNum int) []Secret {
	var secrets []Secret
	patterns := c.getSecretPatterns()

	for _, pattern := range patterns {
		matches := pattern.Pattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 0 {
				secret := Secret{
					Type:        pattern.Name,
					Pattern:     pattern.Pattern.String(),
					File:        filePath,
					Commit:      ref,
					Line:        lineNum,
					Content:     match[0],
					Severity:    pattern.Severity,
					Confidence:  pattern.Confidence,
					Timestamp:   time.Now(),
					Remediation: pattern.Remediation,
				}
				secrets = append(secrets, secret)
			}
		}
	}

	return secrets
}

// checkEntropy checks for high-entropy strings that might be secrets
func (c *SecretChecker) checkEntropy(line, filePath, ref, refType string, lineNum int) []Secret {
	var secrets []Secret

	// Split line into words
	words := strings.Fields(line)
	for _, word := range words {
		// Skip short words and common patterns
		if len(word) < 16 || c.isCommonPattern(word) {
			continue
		}

		// Calculate entropy
		entropy := c.calculateEntropy(word)
		if entropy > 4.0 { // High entropy threshold
			secret := Secret{
				Type:        "High Entropy String",
				Pattern:     "entropy_analysis",
				File:        filePath,
				Commit:      ref,
				Line:        lineNum,
				Content:     word,
				Severity:    "medium",
				Confidence:  math.Min(entropy/6.0, 0.9), // Normalize to 0-0.9
				Timestamp:   time.Now(),
				Remediation: "Review this high-entropy string. If it's a secret, consider rewriting Git history.",
			}
			secrets = append(secrets, secret)
		}
	}

	return secrets
}

// getSecretPatterns returns known secret patterns
func (c *SecretChecker) getSecretPatterns() []SecretPattern {
	return []SecretPattern{
		{
			Name:        "AWS Access Key",
			Pattern:     regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
			Severity:    "high",
			Confidence:  0.95,
			Description: "AWS Access Key ID",
			Remediation: "Rotate AWS access key immediately and rewrite Git history using git filter-repo or BFG.",
		},
		{
			Name:        "AWS Secret Key",
			Pattern:     regexp.MustCompile(`[A-Za-z0-9/+=]{40}`),
			Severity:    "high",
			Confidence:  0.9,
			Description: "AWS Secret Access Key",
			Remediation: "Rotate AWS secret key immediately and rewrite Git history.",
		},
		{
			Name:        "GitHub Token",
			Pattern:     regexp.MustCompile(`ghp_[A-Za-z0-9]{36}`),
			Severity:    "high",
			Confidence:  0.95,
			Description: "GitHub Personal Access Token",
			Remediation: "Revoke GitHub token immediately and rewrite Git history.",
		},
		{
			Name:        "GitHub App Token",
			Pattern:     regexp.MustCompile(`ghs_[A-Za-z0-9]{36}`),
			Severity:    "high",
			Confidence:  0.95,
			Description: "GitHub App Token",
			Remediation: "Revoke GitHub app token immediately and rewrite Git history.",
		},
		{
			Name:        "GitLab Token",
			Pattern:     regexp.MustCompile(`glpat-[A-Za-z0-9_-]{20}`),
			Severity:    "high",
			Confidence:  0.95,
			Description: "GitLab Personal Access Token",
			Remediation: "Revoke GitLab token immediately and rewrite Git history.",
		},
		{
			Name:        "Slack Token",
			Pattern:     regexp.MustCompile(`xox[baprs]-[A-Za-z0-9-]+`),
			Severity:    "high",
			Confidence:  0.9,
			Description: "Slack Bot Token",
			Remediation: "Revoke Slack token immediately and rewrite Git history.",
		},
		{
			Name:        "Discord Token",
			Pattern:     regexp.MustCompile(`[MN][A-Za-z\d]{23}\.[\w-]{6}\.[\w-]{27}`),
			Severity:    "high",
			Confidence:  0.9,
			Description: "Discord Bot Token",
			Remediation: "Regenerate Discord token immediately and rewrite Git history.",
		},
		{
			Name:        "JWT Token",
			Pattern:     regexp.MustCompile(`eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*`),
			Severity:    "medium",
			Confidence:  0.8,
			Description: "JSON Web Token",
			Remediation: "Review JWT token. If it contains sensitive data, regenerate and rewrite Git history.",
		},
		{
			Name:        "Base64 Encoded Secret",
			Pattern:     regexp.MustCompile(`[A-Za-z0-9+/]{40,}={0,2}`),
			Severity:    "medium",
			Confidence:  0.7,
			Description: "Base64 encoded string (potential secret)",
			Remediation: "Review base64 string. If it's a secret, rewrite Git history.",
		},
		{
			Name:        "Private Key",
			Pattern:     regexp.MustCompile(`-----BEGIN [A-Z ]+ PRIVATE KEY-----`),
			Severity:    "high",
			Confidence:  0.95,
			Description: "Private Key",
			Remediation: "Generate new private key immediately and rewrite Git history.",
		},
		{
			Name:        "API Key",
			Pattern:     regexp.MustCompile(`(?i)(api[_-]?key|apikey)[\s]*[:=][\s]*['"]?([A-Za-z0-9_-]{20,})['"]?`),
			Severity:    "high",
			Confidence:  0.85,
			Description: "Generic API Key",
			Remediation: "Rotate API key immediately and rewrite Git history.",
		},
		{
			Name:        "Password",
			Pattern:     regexp.MustCompile(`(?i)(password|passwd|pwd)[\s]*[:=][\s]*['"]?([A-Za-z0-9@#$%^&+=]{8,})['"]?`),
			Severity:    "high",
			Confidence:  0.8,
			Description: "Password field",
			Remediation: "Change password immediately and rewrite Git history.",
		},
	}
}

// calculateEntropy calculates Shannon entropy of a string
func (c *SecretChecker) calculateEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	// Count character frequencies
	freq := make(map[rune]int)
	for _, char := range s {
		freq[char]++
	}

	// Calculate entropy
	entropy := 0.0
	length := float64(len(s))
	for _, count := range freq {
		p := float64(count) / length
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}

	return entropy
}

// isCommonPattern checks if a string matches common non-secret patterns
func (c *SecretChecker) isCommonPattern(s string) bool {
	commonPatterns := []string{
		"http://", "https://", "ftp://", "file://",
		"www.", ".com", ".org", ".net", ".io",
		"localhost", "127.0.0.1", "0.0.0.0",
		"true", "false", "null", "undefined",
		"version", "v1.0", "v2.0", "latest",
	}

	for _, pattern := range commonPatterns {
		if strings.Contains(strings.ToLower(s), pattern) {
			return true
		}
	}

	return false
}

// countHighSeveritySecrets counts high severity secrets
func (c *SecretChecker) countHighSeveritySecrets(secrets []Secret) int {
	count := 0
	for _, secret := range secrets {
		if secret.Severity == "high" {
			count++
		}
	}
	return count
}

// generateRemediation generates remediation steps
func (c *SecretChecker) generateRemediation(secrets []Secret) string {
	remediation := "🚨 CRITICAL: Secrets found in Git history!\n\n"
	remediation += "Immediate Actions Required:\n"
	remediation += "1. Rotate/revoke all exposed credentials immediately\n"
	remediation += "2. Rewrite Git history to remove secrets\n"
	remediation += "3. Notify team members about the exposure\n\n"
	
	remediation += "Tools for History Rewriting:\n"
	remediation += "- git filter-repo: https://github.com/newren/git-filter-repo\n"
	remediation += "- BFG Repo-Cleaner: https://rtyley.github.io/bfg-repo-cleaner/\n\n"
	
	remediation += "Commands:\n"
	remediation += "# Using git filter-repo\n"
	remediation += "git filter-repo --replace-text <(echo 'SECRET_VALUE==>REDACTED')\n\n"
	remediation += "# Using BFG\n"
	remediation += "java -jar bfg.jar --replace-text replacements.txt\n\n"
	
	remediation += "After rewriting history:\n"
	remediation += "git push --force-with-lease origin main\n"
	
	return remediation
}
