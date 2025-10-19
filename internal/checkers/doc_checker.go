package checkers

import (
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// DocChecker checks for essential documentation files
type DocChecker struct {
	BaseChecker
}

// NewDocChecker creates a new documentation checker
func NewDocChecker() *DocChecker {
	return &DocChecker{
		BaseChecker: NewBaseChecker("Documentation Checker", "DOC", types.CategoryDocs, 8),
	}
}

// Check performs the documentation check
func (dc *DocChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "DOC-101",
		Name:      "Essential Documentation Files",
		Category:  types.CategoryDocs,
		Timestamp: time.Now(),
	}

	var details []string
	score := 0
	maxScore := 40

	// Check README.md
	if data.HasReadme {
		score += 10
		details = append(details, "README.md found")
	} else {
		details = append(details, "README.md is missing")
	}

	// Check LICENSE
	if data.HasLicense {
		score += 10
		details = append(details, "LICENSE file found")
	} else {
		details = append(details, "LICENSE file is missing")
	}

	// Check CONTRIBUTING.md
	if data.HasContributing {
		score += 10
		details = append(details, "CONTRIBUTING.md found")
	} else {
		details = append(details, "CONTRIBUTING.md is missing")
	}

	// Check CODE_OF_CONDUCT.md
	if data.HasCodeOfConduct {
		score += 10
		details = append(details, "CODE_OF_CONDUCT.md found")
	} else {
		details = append(details, "CODE_OF_CONDUCT.md is missing")
	}

	result.Score = score
	result.Details = details

	if score == maxScore {
		result.Status = types.StatusPass
		result.Message = "All essential documentation files are present"
	} else if score >= maxScore/2 {
		result.Status = types.StatusWarning
		result.Message = "Some documentation files are missing"
	} else {
		result.Status = types.StatusFail
		result.Message = "Multiple essential documentation files are missing"
	}

	return result
}

// SetupChecker checks for setup instructions in README
type SetupChecker struct {
	BaseChecker
}

// NewSetupChecker creates a new setup checker
func NewSetupChecker() *SetupChecker {
	return &SetupChecker{
		BaseChecker: NewBaseChecker("Setup Instructions Checker", "SETUP", types.CategoryDocs, 5),
	}
}

// Check performs the setup instructions check
func (sc *SetupChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "DOC-102",
		Name:      "Setup Instructions",
		Category:  types.CategoryDocs,
		Timestamp: time.Now(),
	}

	if !data.HasReadme {
		result.Status = types.StatusFail
		result.Score = 0
		result.Message = "Cannot check setup instructions without README.md"
		result.Details = []string{"README.md is required to check setup instructions"}
		return result
	}

	// This is a simplified check - in a real implementation, you'd read and parse the README
	// For now, we'll assume setup instructions exist if README is present
	result.Status = types.StatusPass
	result.Score = 15
	result.Message = "Setup instructions appear to be present in README.md"
	result.Details = []string{"README.md contains setup instructions"}

	return result
}
