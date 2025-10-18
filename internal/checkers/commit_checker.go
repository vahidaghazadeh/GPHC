package checkers

import (
	"fmt"
	"regexp"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// ConventionalCommitChecker checks adherence to conventional commit format
type ConventionalCommitChecker struct {
	BaseChecker
}

// NewConventionalCommitChecker creates a new conventional commit checker
func NewConventionalCommitChecker() *ConventionalCommitChecker {
	return &ConventionalCommitChecker{
		BaseChecker: NewBaseChecker("Conventional Commit Checker", "CONV", types.CategoryCommits, 7),
	}
}

// Check performs the conventional commit check
func (ccc *ConventionalCommitChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "CHQ-301",
		Name:      "Conventional Commit Format",
		Category:  types.CategoryCommits,
		Timestamp: time.Now(),
	}

	if len(data.Commits) == 0 {
		result.Status = types.StatusWarning
		result.Score = 0
		result.Message = "No commits found to analyze"
		result.Details = []string{"⚠️ No commits available for analysis"}
		return result
	}

	// Conventional commit pattern: type(scope): subject
	pattern := regexp.MustCompile(`^(feat|fix|docs|style|refactor|perf|test|chore|build|ci|revert)(\(.+\))?: .+`)

	validCommits := 0
	totalCommits := len(data.Commits)
	var invalidCommits []string

	for _, commit := range data.Commits {
		if pattern.MatchString(commit.Subject) {
			validCommits++
		} else {
			invalidCommits = append(invalidCommits, commit.Subject)
		}
	}

	percentage := float64(validCommits) / float64(totalCommits) * 100
	score := int(percentage)

	var details []string
	details = append(details, fmt.Sprintf("📊 %d of %d commits follow conventional format (%.1f%%)",
		validCommits, totalCommits, percentage))

	if len(invalidCommits) > 0 {
		details = append(details, "❌ Non-standard commits:")
		for i, commit := range invalidCommits {
			if i < 3 { // Show only first 3 invalid commits
				details = append(details, "   - "+commit)
			}
		}
		if len(invalidCommits) > 3 {
			details = append(details, fmt.Sprintf("   ... and %d more", len(invalidCommits)-3))
		}
	}

	result.Score = score
	result.Details = details

	if percentage >= 80 {
		result.Status = types.StatusPass
		result.Message = "Most commits follow conventional format"
	} else if percentage >= 50 {
		result.Status = types.StatusWarning
		result.Message = "Some commits don't follow conventional format"
	} else {
		result.Status = types.StatusFail
		result.Message = "Many commits don't follow conventional format"
	}

	return result
}

// MsgLengthChecker checks commit message length
type MsgLengthChecker struct {
	BaseChecker
}

// NewMsgLengthChecker creates a new message length checker
func NewMsgLengthChecker() *MsgLengthChecker {
	return &MsgLengthChecker{
		BaseChecker: NewBaseChecker("Message Length Checker", "LENGTH", types.CategoryCommits, 6),
	}
}

// Check performs the message length check
func (mlc *MsgLengthChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "CHQ-302",
		Name:      "Commit Message Length",
		Category:  types.CategoryCommits,
		Timestamp: time.Now(),
	}

	if len(data.Commits) == 0 {
		result.Status = types.StatusWarning
		result.Score = 0
		result.Message = "No commits found to analyze"
		result.Details = []string{"⚠️ No commits available for analysis"}
		return result
	}

	maxLength := 72
	validCommits := 0
	totalLength := 0
	var longCommits []string

	for _, commit := range data.Commits {
		length := len(commit.Subject)
		totalLength += length

		if length <= maxLength {
			validCommits++
		} else {
			longCommits = append(longCommits, commit.Subject)
		}
	}

	avgLength := float64(totalLength) / float64(len(data.Commits))
	percentage := float64(validCommits) / float64(len(data.Commits)) * 100
	score := int(percentage)

	var details []string
	details = append(details, fmt.Sprintf("📏 Average commit message length: %.1f characters", avgLength))
	details = append(details, fmt.Sprintf("✅ %d of %d commits are within %d character limit",
		validCommits, len(data.Commits), maxLength))

	if len(longCommits) > 0 {
		details = append(details, "⚠️ Long commit messages:")
		for i, commit := range longCommits {
			if i < 3 { // Show only first 3 long commits
				details = append(details, fmt.Sprintf("   - %s (%d chars)", commit, len(commit)))
			}
		}
		if len(longCommits) > 3 {
			details = append(details, fmt.Sprintf("   ... and %d more", len(longCommits)-3))
		}
	}

	result.Score = score
	result.Details = details

	if percentage >= 90 {
		result.Status = types.StatusPass
		result.Message = "Commit message length is compliant"
	} else if percentage >= 70 {
		result.Status = types.StatusWarning
		result.Message = "Some commit messages are too long"
	} else {
		result.Status = types.StatusFail
		result.Message = "Many commit messages exceed recommended length"
	}

	return result
}

// CommitSizeChecker checks for oversized commits
type CommitSizeChecker struct {
	BaseChecker
}

// NewCommitSizeChecker creates a new commit size checker
func NewCommitSizeChecker() *CommitSizeChecker {
	return &CommitSizeChecker{
		BaseChecker: NewBaseChecker("Commit Size Checker", "SIZE", types.CategoryCommits, 5),
	}
}

// Check performs the commit size check
func (csc *CommitSizeChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "CHQ-303",
		Name:      "Commit Size Analysis",
		Category:  types.CategoryCommits,
		Timestamp: time.Now(),
	}

	if len(data.Commits) == 0 {
		result.Status = types.StatusWarning
		result.Score = 0
		result.Message = "No commits found to analyze"
		result.Details = []string{"⚠️ No commits available for analysis"}
		return result
	}

	totalLines := 0
	largeCommits := 0
	maxLines := 500 // Threshold for large commits

	for _, commit := range data.Commits {
		totalLines += commit.LinesAdded + commit.LinesDeleted
		if commit.LinesAdded+commit.LinesDeleted > maxLines {
			largeCommits++
		}
	}

	avgLines := float64(totalLines) / float64(len(data.Commits))
	percentage := float64(len(data.Commits)-largeCommits) / float64(len(data.Commits)) * 100
	score := int(percentage)

	var details []string
	details = append(details, fmt.Sprintf("📊 Average commit size: %.1f lines", avgLines))
	details = append(details, fmt.Sprintf("✅ %d of %d commits are reasonably sized",
		len(data.Commits)-largeCommits, len(data.Commits)))

	if largeCommits > 0 {
		details = append(details, fmt.Sprintf("⚠️ %d commits exceed %d lines (may indicate 'God Commits')",
			largeCommits, maxLines))
	}

	result.Score = score
	result.Details = details

	if percentage >= 80 {
		result.Status = types.StatusPass
		result.Message = "Average commit size is moderate"
	} else if percentage >= 60 {
		result.Status = types.StatusWarning
		result.Message = "Some commits are quite large"
	} else {
		result.Status = types.StatusFail
		result.Message = "Many commits are oversized"
	}

	return result
}
