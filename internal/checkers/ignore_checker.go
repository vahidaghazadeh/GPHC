package checkers

import (
	"strings"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// IgnoreChecker checks for proper .gitignore usage
type IgnoreChecker struct {
	BaseChecker
}

// NewIgnoreChecker creates a new gitignore checker
func NewIgnoreChecker() *IgnoreChecker {
	return &IgnoreChecker{
		BaseChecker: NewBaseChecker("Gitignore Checker", "IGNORE", types.CategoryDocs, 6),
	}
}

// Check performs the gitignore check
func (ic *IgnoreChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "IG-201",
		Name:      "Gitignore Configuration",
		Category:  types.CategoryDocs,
		Timestamp: time.Now(),
	}

	var details []string
	score := 0

	// Check if .gitignore exists
	if !data.HasGitignore {
		result.Status = types.StatusFail
		result.Score = 0
		result.Message = ".gitignore file is missing"
		result.Details = []string{"❌ .gitignore file is missing"}
		return result
	}

	score += 10
	details = append(details, "✅ .gitignore file is present")

	// Check for common patterns
	commonPatterns := []string{
		"node_modules",
		"__pycache__",
		"target/",
		".env",
		"*.log",
		".DS_Store",
		"*.tmp",
		"*.swp",
	}

	gitignoreContent := strings.ToLower(data.GitignoreContent)
	patternsFound := 0

	for _, pattern := range commonPatterns {
		if strings.Contains(gitignoreContent, pattern) {
			patternsFound++
			details = append(details, "✅ Contains "+pattern)
		}
	}

	// Calculate score based on patterns found
	if patternsFound >= len(commonPatterns)/2 {
		score += 10
		result.Status = types.StatusPass
		result.Message = ".gitignore is properly configured"
	} else {
		score += 5
		result.Status = types.StatusWarning
		result.Message = ".gitignore could be improved with more patterns"
	}

	result.Score = score
	result.Details = details

	return result
}
