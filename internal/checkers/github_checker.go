package checkers

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/opsource/gphc/internal/integration"
	"github.com/opsource/gphc/pkg/types"
)

// GitHubIntegrationChecker checks GitHub-specific features
type GitHubIntegrationChecker struct {
	BaseChecker
}

// NewGitHubIntegrationChecker creates a new GitHub integration checker
func NewGitHubIntegrationChecker() *GitHubIntegrationChecker {
	return &GitHubIntegrationChecker{
		BaseChecker: BaseChecker{
			id:       "GH-601",
			name:     "GitHub Integration",
			category: types.CategoryHygiene,
		},
	}
}

// Check performs GitHub integration checks
func (c *GitHubIntegrationChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        c.ID(),
		Name:      c.Name(),
		Category:  c.Category(),
		Status:    types.StatusPass,
		Score:     0,
		Message:   "GitHub integration checks completed",
		Details:   []string{},
		Timestamp: time.Now(),
	}

	// Check if we have GitHub token
	githubClient := integration.NewGitHubClient()
	if !githubClient.IsAuthenticated() {
		result.Status = types.StatusWarning
		result.Score = 30
		result.Message = "GitHub token not provided - limited integration checks"
		result.Details = append(result.Details, "Set GPHC_TOKEN or GITHUB_TOKEN environment variable for full GitHub integration")
		result.Details = append(result.Details, "Without token, only basic checks are performed")
		return result
	}

	// Get repository remote URL
	remoteURL, err := c.getRemoteURL(data.Path)
	if err != nil {
		result.Status = types.StatusWarning
		result.Score = 40
		result.Message = "Could not determine GitHub repository"
		result.Details = append(result.Details, fmt.Sprintf("Error getting remote URL: %v", err))
		return result
	}

	// Extract owner and repo from URL
	owner, repo, err := integration.ExtractRepoInfo(remoteURL)
	if err != nil {
		result.Status = types.StatusWarning
		result.Score = 40
		result.Message = "Could not parse GitHub repository URL"
		result.Details = append(result.Details, fmt.Sprintf("Error parsing URL: %v", err))
		return result
	}

	result.Details = append(result.Details, fmt.Sprintf("Repository: %s/%s", owner, repo))

	// Check repository info
	repoInfo, err := githubClient.GetRepositoryInfo(owner, repo)
	if err != nil {
		result.Status = types.StatusWarning
		result.Score = 50
		result.Message = "Could not fetch repository information"
		result.Details = append(result.Details, fmt.Sprintf("GitHub API error: %v", err))
		return result
	}

	// Analyze repository features
	score := 0
	totalChecks := 0

	// Check if repository has issues enabled
	if repoInfo.HasIssues {
		score += 10
		result.Details = append(result.Details, "✅ Issues are enabled")
	} else {
		result.Details = append(result.Details, "⚠️ Issues are disabled")
	}
	totalChecks++

	// Check if repository has projects enabled
	if repoInfo.HasProjects {
		score += 10
		result.Details = append(result.Details, "✅ Projects are enabled")
	} else {
		result.Details = append(result.Details, "⚠️ Projects are disabled")
	}
	totalChecks++

	// Check if repository has wiki enabled
	if repoInfo.HasWiki {
		score += 5
		result.Details = append(result.Details, "✅ Wiki is enabled")
	} else {
		result.Details = append(result.Details, "ℹ️ Wiki is disabled")
	}
	totalChecks++

	// Check branch protection
	protection, err := githubClient.GetBranchProtection(owner, repo, repoInfo.DefaultBranch)
	if err != nil {
		result.Details = append(result.Details, fmt.Sprintf("⚠️ Could not check branch protection: %v", err))
	} else if protection != nil {
		score += 25
		result.Details = append(result.Details, "✅ Branch protection is enabled")
		
		if protection.RequiredPullRequestReviews != nil {
			if protection.RequiredPullRequestReviews.RequiredApprovingReviewCount > 0 {
				score += 15
				result.Details = append(result.Details, fmt.Sprintf("✅ Required %d reviewer(s)", protection.RequiredPullRequestReviews.RequiredApprovingReviewCount))
			}
			if protection.RequiredPullRequestReviews.RequireCodeOwnerReviews {
				score += 10
				result.Details = append(result.Details, "✅ Code owner reviews required")
			}
		}
		
		if protection.RequiredStatusChecks != nil && len(protection.RequiredStatusChecks.Contexts) > 0 {
			score += 10
			result.Details = append(result.Details, fmt.Sprintf("✅ Required status checks: %s", strings.Join(protection.RequiredStatusChecks.Contexts, ", ")))
		}
	} else {
		result.Details = append(result.Details, "❌ Branch protection is not enabled")
	}
	totalChecks++

	// Check GitHub Actions workflows
	workflows, err := githubClient.GetWorkflows(owner, repo)
	if err != nil {
		result.Details = append(result.Details, fmt.Sprintf("⚠️ Could not check workflows: %v", err))
	} else if len(workflows) > 0 {
		score += 20
		result.Details = append(result.Details, fmt.Sprintf("✅ Found %d workflow(s)", len(workflows)))
		
		activeWorkflows := 0
		for _, workflow := range workflows {
			if workflow.State == "active" {
				activeWorkflows++
			}
		}
		
		if activeWorkflows > 0 {
			score += 10
			result.Details = append(result.Details, fmt.Sprintf("✅ %d active workflow(s)", activeWorkflows))
		}
	} else {
		result.Details = append(result.Details, "⚠️ No GitHub Actions workflows found")
	}
	totalChecks++

	// Check contributors activity
	contributors, err := githubClient.GetContributors(owner, repo)
	if err != nil {
		result.Details = append(result.Details, fmt.Sprintf("⚠️ Could not check contributors: %v", err))
	} else {
		result.Details = append(result.Details, fmt.Sprintf("📊 Found %d contributor(s)", len(contributors)))
		
		if len(contributors) > 1 {
			score += 10
			result.Details = append(result.Details, "✅ Multiple contributors")
		} else {
			result.Details = append(result.Details, "⚠️ Single contributor repository")
		}
	}
	totalChecks++

	// Calculate final score
	if score >= 80 {
		result.Status = types.StatusPass
		result.Message = "Excellent GitHub integration and configuration"
	} else if score >= 60 {
		result.Status = types.StatusPass
		result.Message = "Good GitHub integration with room for improvement"
	} else if score >= 40 {
		result.Status = types.StatusWarning
		result.Message = "Basic GitHub integration - consider enabling more features"
	} else {
		result.Status = types.StatusFail
		result.Message = "Limited GitHub integration - significant improvements needed"
	}

	result.Score = score
	return result
}

// getRemoteURL gets the remote URL of the repository
func (c *GitHubIntegrationChecker) getRemoteURL(repoPath string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoPath
	
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(string(output)), nil
}

// isGitHubRepository checks if the repository is hosted on GitHub
func (c *GitHubIntegrationChecker) isGitHubRepository(repoPath string) bool {
	remoteURL, err := c.getRemoteURL(repoPath)
	if err != nil {
		return false
	}
	
	return strings.Contains(remoteURL, "github.com")
}
