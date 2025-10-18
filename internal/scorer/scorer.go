package scorer

import (
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// Scorer calculates the overall health score from check results
type Scorer struct {
	results []types.CheckResult
}

// NewScorer creates a new scorer
func NewScorer() *Scorer {
	return &Scorer{
		results: make([]types.CheckResult, 0),
	}
}

// AddResult adds a check result to the scorer
func (s *Scorer) AddResult(result types.CheckResult) {
	s.results = append(s.results, result)
}

// CalculateHealthReport calculates the overall health report
func (s *Scorer) CalculateHealthReport() *types.HealthReport {
	if len(s.results) == 0 {
		return &types.HealthReport{
			OverallScore: 0,
			Grade:        "F",
			Results:      []types.CheckResult{},
			Summary: types.ReportSummary{
				TotalChecks:   0,
				PassedChecks:  0,
				FailedChecks:  0,
				WarningChecks: 0,
			},
			Timestamp: time.Now(),
		}
	}

	// Calculate weighted score
	totalWeightedScore := 0
	totalWeight := 0

	summary := types.ReportSummary{
		TotalChecks: len(s.results),
	}

	for _, result := range s.results {
		// Weight the score based on the checker's importance
		weight := getWeightForCategory(result.Category)
		totalWeight += weight
		totalWeightedScore += result.Score * weight

		// Count statuses
		switch result.Status {
		case types.StatusPass:
			summary.PassedChecks++
		case types.StatusFail:
			summary.FailedChecks++
		case types.StatusWarning:
			summary.WarningChecks++
		}
	}

	// Calculate overall score (0-100)
	var overallScore int
	if totalWeight > 0 {
		overallScore = totalWeightedScore / totalWeight
	}

	// Calculate grade
	grade := calculateGrade(overallScore)

	return &types.HealthReport{
		OverallScore: overallScore,
		Grade:        grade,
		Results:      s.results,
		Summary:      summary,
		Timestamp:    time.Now(),
	}
}

// getWeightForCategory returns the weight for a category
func getWeightForCategory(category types.Category) int {
	switch category {
	case types.CategoryDocs:
		return 3 // Documentation is important
	case types.CategoryCommits:
		return 4 // Commit quality is very important
	case types.CategoryHygiene:
		return 2 // Hygiene is important but less critical
	default:
		return 1
	}
}

// calculateGrade converts score to letter grade
func calculateGrade(score int) string {
	switch {
	case score >= 95:
		return "A+"
	case score >= 90:
		return "A"
	case score >= 85:
		return "A-"
	case score >= 80:
		return "B+"
	case score >= 75:
		return "B"
	case score >= 70:
		return "B-"
	case score >= 65:
		return "C+"
	case score >= 60:
		return "C"
	case score >= 55:
		return "C-"
	case score >= 50:
		return "D"
	default:
		return "F"
	}
}

// GetCategoryResults returns results grouped by category
func (s *Scorer) GetCategoryResults() map[types.Category][]types.CheckResult {
	categoryResults := make(map[types.Category][]types.CheckResult)

	for _, result := range s.results {
		categoryResults[result.Category] = append(categoryResults[result.Category], result)
	}

	return categoryResults
}

// GetFailedChecks returns all failed checks
func (s *Scorer) GetFailedChecks() []types.CheckResult {
	var failed []types.CheckResult
	for _, result := range s.results {
		if result.Status == types.StatusFail {
			failed = append(failed, result)
		}
	}
	return failed
}

// GetWarningChecks returns all warning checks
func (s *Scorer) GetWarningChecks() []types.CheckResult {
	var warnings []types.CheckResult
	for _, result := range s.results {
		if result.Status == types.StatusWarning {
			warnings = append(warnings, result)
		}
	}
	return warnings
}

// GetNextSteps generates actionable next steps based on failed checks
func (s *Scorer) GetNextSteps() []string {
	var steps []string

	for _, result := range s.results {
		if result.Status == types.StatusFail {
			switch result.ID {
			case "DOC-101":
				if !containsString(steps, "Add missing documentation files") {
					steps = append(steps, "Add missing documentation files (README.md, LICENSE, etc.)")
				}
			case "IG-201":
				if !containsString(steps, "Create or improve .gitignore") {
					steps = append(steps, "Create or improve .gitignore file")
				}
			case "CHQ-301":
				if !containsString(steps, "Follow conventional commit format") {
					steps = append(steps, "Follow conventional commit format (feat:, fix:, etc.)")
				}
			case "CHQ-302":
				if !containsString(steps, "Shorten commit messages") {
					steps = append(steps, "Keep commit messages under 72 characters")
				}
			case "CLEAN-401":
				if !containsString(steps, "Delete merged branches") {
					steps = append(steps, "Delete merged local branches")
				}
			case "CLEAN-402":
				if !containsString(steps, "Review stale branches") {
					steps = append(steps, "Review and delete stale branches")
				}
			}
		}
	}

	return steps
}

// containsString checks if a slice contains a string
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
