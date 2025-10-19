package checkers

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// AuthorStats represents statistics for a commit author
type AuthorStats struct {
	Name      string
	Email     string
	Commits   int
	Percentage float64
}

// CommitAuthorInsightsChecker analyzes commit author patterns
type CommitAuthorInsightsChecker struct {
	BaseChecker
}

// NewCommitAuthorInsightsChecker creates a new commit author insights checker
func NewCommitAuthorInsightsChecker() *CommitAuthorInsightsChecker {
	return &CommitAuthorInsightsChecker{
		BaseChecker: BaseChecker{
			id:       "CAI-701",
			name:     "Commit Author Insights",
			category: types.CategoryCommits,
		},
	}
}

// Check performs commit author analysis
func (c *CommitAuthorInsightsChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        c.ID(),
		Name:      c.Name(),
		Category:  c.Category(),
		Status:    types.StatusPass,
		Score:     0,
		Message:   "Commit author analysis completed",
		Details:   []string{},
		Timestamp: time.Now(),
	}

	// Check if we have commits to analyze
	if len(data.Commits) == 0 {
		result.Status = types.StatusWarning
		result.Score = 0
		result.Message = "No commits found for analysis"
		result.Details = append(result.Details, "Repository has no commit history")
		return result
	}

	// Analyze commit authors
	authorStats := c.analyzeAuthors(data.Commits)
	
	// Calculate total commits
	totalCommits := len(data.Commits)
	
	// Sort authors by commit count (descending)
	sort.Slice(authorStats, func(i, j int) bool {
		return authorStats[i].Commits > authorStats[j].Commits
	})

	// Generate insights
	result.Details = append(result.Details, fmt.Sprintf("Contributors: %d", len(authorStats)))
	result.Details = append(result.Details, fmt.Sprintf("Total commits: %d", totalCommits))

	// Show top contributors (up to 5)
	topContributors := authorStats
	if len(authorStats) > 5 {
		topContributors = authorStats[:5]
	}

	for i, author := range topContributors {
		rank := i + 1
		result.Details = append(result.Details, fmt.Sprintf("  %d. %s (%d commits, %.1f%%)", 
			rank, author.Name, author.Commits, author.Percentage))
	}

	// Check for single author dominance
	score := 100
	if len(authorStats) > 0 {
		topAuthor := authorStats[0]
		
		// Single author dominance warning (>70%)
		if topAuthor.Percentage > 70 {
			result.Status = types.StatusWarning
			score = 60
			result.Message = "Single author dominance detected"
			result.Details = append(result.Details, fmt.Sprintf("Single Author Dominance Detected (>70%%)"))
			result.Details = append(result.Details, fmt.Sprintf("Top author: %s (%.1f%%)", topAuthor.Name, topAuthor.Percentage))
			result.Details = append(result.Details, "Consider encouraging more team participation")
		} else if topAuthor.Percentage > 50 {
			result.Status = types.StatusPass
			score = 80
			result.Message = "Good contributor distribution with some dominance"
			result.Details = append(result.Details, fmt.Sprintf("Top author: %s (%.1f%%)", topAuthor.Name, topAuthor.Percentage))
			result.Details = append(result.Details, "Consider encouraging more balanced contributions")
		} else {
			result.Status = types.StatusPass
			score = 100
			result.Message = "Excellent contributor distribution"
			result.Details = append(result.Details, fmt.Sprintf("Top author: %s (%.1f%%)", topAuthor.Name, topAuthor.Percentage))
			result.Details = append(result.Details, "Well-distributed contributions across team")
		}

		// Check for single contributor project
		if len(authorStats) == 1 {
			result.Status = types.StatusWarning
			score = 40
			result.Message = "Single contributor project - high bus factor risk"
			result.Details = append(result.Details, "Single Contributor Project")
			result.Details = append(result.Details, "High risk of 'Bus Factor' - project depends on one person")
			result.Details = append(result.Details, "Consider onboarding additional contributors")
		}

		// Check for very low contributor count
		if len(authorStats) == 2 {
			result.Status = types.StatusWarning
			score = 60
			result.Message = "Low contributor count - moderate bus factor risk"
			result.Details = append(result.Details, "Low Contributor Count")
			result.Details = append(result.Details, "Moderate risk of 'Bus Factor'")
			result.Details = append(result.Details, "Consider expanding the contributor base")
		}

		// Check for inactive contributors
		inactiveContributors := 0
		for _, author := range authorStats {
			if author.Percentage < 5 { // Less than 5% contribution
				inactiveContributors++
			}
		}

		if inactiveContributors > 0 && len(authorStats) > 3 {
			result.Details = append(result.Details, fmt.Sprintf("%d contributor(s) with minimal activity (<5%%)", inactiveContributors))
		}

		// Check for email consistency
		emailConsistency := c.checkEmailConsistency(authorStats)
		if !emailConsistency {
			result.Details = append(result.Details, "Inconsistent email addresses detected")
			result.Details = append(result.Details, "Consider configuring git config user.email consistently")
		} else {
			result.Details = append(result.Details, "Email addresses are consistent")
		}

	} else {
		result.Status = types.StatusWarning
		result.Score = 0
		result.Message = "No author information found"
		result.Details = append(result.Details, "Unable to extract author information from commits")
	}

	result.Score = score
	return result
}

// analyzeAuthors analyzes commit authors and returns statistics
func (c *CommitAuthorInsightsChecker) analyzeAuthors(commits []types.CommitInfo) []AuthorStats {
	authorMap := make(map[string]*AuthorStats)
	totalCommits := len(commits)

	// Count commits per author
	for _, commit := range commits {
		authorKey := fmt.Sprintf("%s <%s>", commit.Author, commit.AuthorEmail)
		
		if stats, exists := authorMap[authorKey]; exists {
			stats.Commits++
		} else {
			authorMap[authorKey] = &AuthorStats{
				Name:    commit.Author,
				Email:   commit.AuthorEmail,
				Commits: 1,
			}
		}
	}

	// Convert to slice and calculate percentages
	var authorStats []AuthorStats
	for _, stats := range authorMap {
		stats.Percentage = float64(stats.Commits) / float64(totalCommits) * 100
		authorStats = append(authorStats, *stats)
	}

	return authorStats
}

// checkEmailConsistency checks if authors use consistent email addresses
func (c *CommitAuthorInsightsChecker) checkEmailConsistency(authorStats []AuthorStats) bool {
	if len(authorStats) <= 1 {
		return true
	}

	// Check if all authors have proper email format
	for _, author := range authorStats {
		if !strings.Contains(author.Email, "@") {
			return false
		}
	}

	return true
}
