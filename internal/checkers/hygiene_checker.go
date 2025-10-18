package checkers

import (
	"fmt"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// LocalBranchChecker checks for stale local branches
type LocalBranchChecker struct {
	BaseChecker
}

// NewLocalBranchChecker creates a new local branch checker
func NewLocalBranchChecker() *LocalBranchChecker {
	return &LocalBranchChecker{
		BaseChecker: NewBaseChecker("Local Branch Checker", "LOCAL", types.CategoryHygiene, 6),
	}
}

// Check performs the local branch check
func (lbc *LocalBranchChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "CLEAN-401",
		Name:      "Local Branch Cleanup",
		Category:  types.CategoryHygiene,
		Timestamp: time.Now(),
	}

	if len(data.Branches) == 0 {
		result.Status = types.StatusWarning
		result.Score = 0
		result.Message = "No branches found to analyze"
		result.Details = []string{"⚠️ No branches available for analysis"}
		return result
	}

	mergedBranches := 0
	var mergedBranchNames []string

	for _, branch := range data.Branches {
		if branch.IsMerged && branch.Name != "main" && branch.Name != "master" {
			mergedBranches++
			mergedBranchNames = append(mergedBranchNames, branch.Name)
		}
	}

	var details []string
	details = append(details, fmt.Sprintf("📊 Found %d local branches", len(data.Branches)))
	details = append(details, fmt.Sprintf("✅ %d branches are active", len(data.Branches)-mergedBranches))

	if mergedBranches > 0 {
		details = append(details, fmt.Sprintf("⚠️ %d branches are merged but not deleted:", mergedBranches))
		for i, branchName := range mergedBranchNames {
			if i < 5 { // Show only first 5 merged branches
				details = append(details, "   - "+branchName)
			}
		}
		if len(mergedBranchNames) > 5 {
			details = append(details, fmt.Sprintf("   ... and %d more", len(mergedBranchNames)-5))
		}
	}

	// Calculate score based on percentage of merged branches
	if len(data.Branches) == 0 {
		score := 100
		result.Score = score
		result.Status = types.StatusPass
		result.Message = "No local branches to clean up"
	} else {
		percentage := float64(len(data.Branches)-mergedBranches) / float64(len(data.Branches)) * 100
		score := int(percentage)
		result.Score = score
		result.Details = details

		if percentage >= 80 {
			result.Status = types.StatusPass
			result.Message = "Most local branches are active"
		} else if percentage >= 60 {
			result.Status = types.StatusWarning
			result.Message = "Some merged branches should be cleaned up"
		} else {
			result.Status = types.StatusFail
			result.Message = "Many merged branches need cleanup"
		}
	}

	return result
}

// StaleBranchChecker checks for untouched branches
type StaleBranchChecker struct {
	BaseChecker
}

// NewStaleBranchChecker creates a new stale branch checker
func NewStaleBranchChecker() *StaleBranchChecker {
	return &StaleBranchChecker{
		BaseChecker: NewBaseChecker("Stale Branch Checker", "STALE", types.CategoryHygiene, 5),
	}
}

// Check performs the stale branch check
func (sbc *StaleBranchChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "CLEAN-402",
		Name:      "Stale Branch Detection",
		Category:  types.CategoryHygiene,
		Timestamp: time.Now(),
	}

	if len(data.Branches) == 0 {
		result.Status = types.StatusWarning
		result.Score = 0
		result.Message = "No branches found to analyze"
		result.Details = []string{"⚠️ No branches available for analysis"}
		return result
	}

	staleBranches := 0
	var staleBranchNames []string
	staleThreshold := 60 // days

	for _, branch := range data.Branches {
		if branch.IsStale && branch.Name != "main" && branch.Name != "master" {
			staleBranches++
			daysSinceLastCommit := int(time.Since(branch.LastCommit).Hours() / 24)
			staleBranchNames = append(staleBranchNames,
				fmt.Sprintf("%s (%d days ago)", branch.Name, daysSinceLastCommit))
		}
	}

	var details []string
	details = append(details, fmt.Sprintf("📊 Found %d local branches", len(data.Branches)))
	details = append(details, fmt.Sprintf("✅ %d branches are active", len(data.Branches)-staleBranches))

	if staleBranches > 0 {
		details = append(details, fmt.Sprintf("⚠️ %d branches are stale (older than %d days):",
			staleBranches, staleThreshold))
		for i, branchInfo := range staleBranchNames {
			if i < 5 { // Show only first 5 stale branches
				details = append(details, "   - "+branchInfo)
			}
		}
		if len(staleBranchNames) > 5 {
			details = append(details, fmt.Sprintf("   ... and %d more", len(staleBranchNames)-5))
		}
	}

	// Calculate score based on percentage of stale branches
	if len(data.Branches) == 0 {
		score := 100
		result.Score = score
		result.Status = types.StatusPass
		result.Message = "No branches to analyze"
	} else {
		percentage := float64(len(data.Branches)-staleBranches) / float64(len(data.Branches)) * 100
		score := int(percentage)
		result.Score = score
		result.Details = details

		if percentage >= 80 {
			result.Status = types.StatusPass
			result.Message = "Most branches are active"
		} else if percentage >= 60 {
			result.Status = types.StatusWarning
			result.Message = "Some branches are stale"
		} else {
			result.Status = types.StatusFail
			result.Message = "Many branches are stale"
		}
	}

	return result
}

// BareRepoChecker checks for direct commits to main branch
type BareRepoChecker struct {
	BaseChecker
}

// NewBareRepoChecker creates a new bare repo checker
func NewBareRepoChecker() *BareRepoChecker {
	return &BareRepoChecker{
		BaseChecker: NewBaseChecker("Bare Repository Checker", "BARE", types.CategoryHygiene, 4),
	}
}

// Check performs the bare repository check
func (brc *BareRepoChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "CLEAN-403",
		Name:      "Main Branch Protection",
		Category:  types.CategoryHygiene,
		Timestamp: time.Now(),
	}

	// This is a simplified check - in a real implementation, you'd check for:
	// - GitHub branch protection rules (requires API)
	// - Pre-commit hooks
	// - CI/CD integration

	var details []string
	details = append(details, "ℹ️ This check requires GitHub API integration")
	details = append(details, "ℹ️ Consider enabling branch protection rules")
	details = append(details, "ℹ️ Set up pre-commit hooks for additional safety")

	result.Status = types.StatusWarning
	result.Score = 50
	result.Message = "Branch protection status unknown (requires GitHub API)"
	result.Details = details

	return result
}
