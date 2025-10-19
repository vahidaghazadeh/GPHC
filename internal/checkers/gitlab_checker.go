package checkers

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/opsource/gphc/internal/integration"
	"github.com/opsource/gphc/pkg/types"
)

// GitLabIntegrationChecker checks GitLab-specific features
type GitLabIntegrationChecker struct {
	BaseChecker
}

// NewGitLabIntegrationChecker creates a new GitLab integration checker
func NewGitLabIntegrationChecker() *GitLabIntegrationChecker {
	return &GitLabIntegrationChecker{
		BaseChecker: BaseChecker{
			id:       "GL-602",
			name:     "GitLab Integration",
			category: types.CategoryHygiene,
		},
	}
}

// Check performs GitLab integration checks
func (c *GitLabIntegrationChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        c.ID(),
		Name:      c.Name(),
		Category:  c.Category(),
		Status:    types.StatusPass,
		Score:     0,
		Message:   "GitLab integration checks completed",
		Details:   []string{},
		Timestamp: time.Now(),
	}

	// Check if we have GitLab token
	gitlabClient := integration.NewGitLabClient()
	if !gitlabClient.IsAuthenticated() {
		result.Status = types.StatusWarning
		result.Score = 30
		result.Message = "GitLab token not provided - limited integration checks"
		result.Details = append(result.Details, "Set GPHC_TOKEN or GITLAB_TOKEN environment variable for full GitLab integration")
		result.Details = append(result.Details, "Without token, only basic checks are performed")
		return result
	}

	// Get repository remote URL
	remoteURL, err := c.getRemoteURL(data.Path)
	if err != nil {
		result.Status = types.StatusWarning
		result.Score = 40
		result.Message = "Could not determine GitLab repository"
		result.Details = append(result.Details, fmt.Sprintf("Error getting remote URL: %v", err))
		return result
	}

	// Extract project path from URL
	projectPath, err := integration.ExtractGitLabProjectInfo(remoteURL)
	if err != nil {
		result.Status = types.StatusWarning
		result.Score = 40
		result.Message = "Could not parse GitLab repository URL"
		result.Details = append(result.Details, fmt.Sprintf("Error parsing URL: %v", err))
		return result
	}

	result.Details = append(result.Details, fmt.Sprintf("Project: %s", projectPath))

	// Check project info
	projectInfo, err := gitlabClient.GetProjectInfo(projectPath)
	if err != nil {
		result.Status = types.StatusWarning
		result.Score = 50
		result.Message = "Could not fetch project information"
		result.Details = append(result.Details, fmt.Sprintf("GitLab API error: %v", err))
		return result
	}

	// Analyze project features
	score := 0
	totalChecks := 0

	// Check if project has issues enabled
	if projectInfo.IssuesEnabled {
		score += 10
		result.Details = append(result.Details, "✅ Issues are enabled")
	} else {
		result.Details = append(result.Details, "⚠️ Issues are disabled")
	}
	totalChecks++

	// Check if project has merge requests enabled
	if projectInfo.MergeRequestsEnabled {
		score += 10
		result.Details = append(result.Details, "✅ Merge requests are enabled")
	} else {
		result.Details = append(result.Details, "⚠️ Merge requests are disabled")
	}
	totalChecks++

	// Check if project has wiki enabled
	if projectInfo.WikiEnabled {
		score += 5
		result.Details = append(result.Details, "✅ Wiki is enabled")
	} else {
		result.Details = append(result.Details, "ℹ️ Wiki is disabled")
	}
	totalChecks++

	// Check if project has snippets enabled
	if projectInfo.SnippetsEnabled {
		score += 5
		result.Details = append(result.Details, "✅ Snippets are enabled")
	} else {
		result.Details = append(result.Details, "ℹ️ Snippets are disabled")
	}
	totalChecks++

	// Check branch protection
	protection, err := gitlabClient.GetBranchProtection(projectPath, projectInfo.DefaultBranch)
	if err != nil {
		result.Details = append(result.Details, fmt.Sprintf("⚠️ Could not check branch protection: %v", err))
	} else if protection != nil {
		score += 25
		result.Details = append(result.Details, "✅ Branch protection is enabled")
		
		if len(protection.PushAccessLevels) > 0 {
			score += 10
			result.Details = append(result.Details, "✅ Push access is restricted")
		}
		
		if len(protection.MergeAccessLevels) > 0 {
			score += 10
			result.Details = append(result.Details, "✅ Merge access is restricted")
		}
		
		if protection.CodeOwnerApprovalRequired {
			score += 10
			result.Details = append(result.Details, "✅ Code owner approval required")
		}
	} else {
		result.Details = append(result.Details, "❌ Branch protection is not enabled")
	}
	totalChecks++

	// Check GitLab CI/CD pipelines
	pipelines, err := gitlabClient.GetPipelines(projectPath)
	if err != nil {
		result.Details = append(result.Details, fmt.Sprintf("⚠️ Could not check pipelines: %v", err))
	} else if len(pipelines) > 0 {
		score += 20
		result.Details = append(result.Details, fmt.Sprintf("✅ Found %d pipeline(s)", len(pipelines)))
		
		successfulPipelines := 0
		for _, pipeline := range pipelines {
			if pipeline.Status == "success" {
				successfulPipelines++
			}
		}
		
		if successfulPipelines > 0 {
			score += 10
			result.Details = append(result.Details, fmt.Sprintf("✅ %d successful pipeline(s)", successfulPipelines))
		}
	} else {
		result.Details = append(result.Details, "⚠️ No GitLab CI/CD pipelines found")
	}
	totalChecks++

	// Check contributors activity
	contributors, err := gitlabClient.GetContributors(projectPath)
	if err != nil {
		result.Details = append(result.Details, fmt.Sprintf("⚠️ Could not check contributors: %v", err))
	} else {
		result.Details = append(result.Details, fmt.Sprintf("📊 Found %d contributor(s)", len(contributors)))
		
		if len(contributors) > 1 {
			score += 10
			result.Details = append(result.Details, "✅ Multiple contributors")
		} else {
			result.Details = append(result.Details, "⚠️ Single contributor project")
		}
	}
	totalChecks++

	// Check open merge requests
	mergeRequests, err := gitlabClient.GetMergeRequests(projectPath)
	if err != nil {
		result.Details = append(result.Details, fmt.Sprintf("⚠️ Could not check merge requests: %v", err))
	} else {
		result.Details = append(result.Details, fmt.Sprintf("📊 Found %d open merge request(s)", len(mergeRequests)))
		
		if len(mergeRequests) > 0 {
			score += 5
			result.Details = append(result.Details, "✅ Active development with open MRs")
		}
	}
	totalChecks++

	// Calculate final score
	if score >= 80 {
		result.Status = types.StatusPass
		result.Message = "Excellent GitLab integration and configuration"
	} else if score >= 60 {
		result.Status = types.StatusPass
		result.Message = "Good GitLab integration with room for improvement"
	} else if score >= 40 {
		result.Status = types.StatusWarning
		result.Message = "Basic GitLab integration - consider enabling more features"
	} else {
		result.Status = types.StatusFail
		result.Message = "Limited GitLab integration - significant improvements needed"
	}

	result.Score = score
	return result
}

// getRemoteURL gets the remote URL of the repository
func (c *GitLabIntegrationChecker) getRemoteURL(repoPath string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoPath
	
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(string(output)), nil
}

// isGitLabRepository checks if the repository is hosted on GitLab
func (c *GitLabIntegrationChecker) isGitLabRepository(repoPath string) bool {
	remoteURL, err := c.getRemoteURL(repoPath)
	if err != nil {
		return false
	}
	
	return strings.Contains(remoteURL, "gitlab.com") || strings.Contains(remoteURL, "gitlab")
}
