package scorer

import (
	"strings"
	"testing"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

func TestNewScorer(t *testing.T) {
	scorer := NewScorer()

	if scorer == nil {
		t.Error("Expected scorer, got nil")
	}

	if len(scorer.results) != 0 {
		t.Errorf("Expected empty results, got %d results", len(scorer.results))
	}
}

func TestAddResult(t *testing.T) {
	scorer := NewScorer()

	result := types.CheckResult{
		ID:        "TEST-001",
		Name:      "Test Check",
		Status:    types.StatusPass,
		Score:     100,
		Message:   "Test passed",
		Category:  types.CategoryDocs,
		Timestamp: time.Now(),
	}

	scorer.AddResult(result)

	if len(scorer.results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(scorer.results))
	}

	if scorer.results[0].ID != "TEST-001" {
		t.Errorf("Expected ID 'TEST-001', got '%s'", scorer.results[0].ID)
	}
}

func TestCalculateHealthReport(t *testing.T) {
	scorer := NewScorer()

	// Add some test results
	scorer.AddResult(types.CheckResult{
		ID:        "DOC-101",
		Name:      "Documentation Check",
		Status:    types.StatusPass,
		Score:     100,
		Message:   "All docs present",
		Category:  types.CategoryDocs,
		Timestamp: time.Now(),
	})

	scorer.AddResult(types.CheckResult{
		ID:        "CHQ-301",
		Name:      "Commit Check",
		Status:    types.StatusPass,
		Score:     80,
		Message:   "Commits look good",
		Category:  types.CategoryCommits,
		Timestamp: time.Now(),
	})

	report := scorer.CalculateHealthReport()

	if report == nil {
		t.Error("Expected report, got nil")
		return
	}

	if report.OverallScore <= 0 {
		t.Errorf("Expected positive score, got %d", report.OverallScore)
	}

	if report.Grade == "" {
		t.Error("Expected grade, got empty string")
	}

	if len(report.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(report.Results))
	}

	if report.Summary.TotalChecks != 2 {
		t.Errorf("Expected 2 total checks, got %d", report.Summary.TotalChecks)
	}

	if report.Summary.PassedChecks != 2 {
		t.Errorf("Expected 2 passed checks, got %d", report.Summary.PassedChecks)
	}
}

func TestCalculateHealthReportEmpty(t *testing.T) {
	scorer := NewScorer()

	report := scorer.CalculateHealthReport()

	if report == nil {
		t.Error("Expected report, got nil")
		return
	}

	if report.OverallScore != 0 {
		t.Errorf("Expected score 0, got %d", report.OverallScore)
	}

	if report.Grade != "F" {
		t.Errorf("Expected grade 'F', got '%s'", report.Grade)
	}

	if len(report.Results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(report.Results))
	}
}

func TestGetCategoryResults(t *testing.T) {
	scorer := NewScorer()

	// Add results from different categories
	scorer.AddResult(types.CheckResult{
		ID:       "DOC-101",
		Category: types.CategoryDocs,
	})

	scorer.AddResult(types.CheckResult{
		ID:       "CHQ-301",
		Category: types.CategoryCommits,
	})

	scorer.AddResult(types.CheckResult{
		ID:       "DOC-102",
		Category: types.CategoryDocs,
	})

	categoryResults := scorer.GetCategoryResults()

	if len(categoryResults) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categoryResults))
	}

	if len(categoryResults[types.CategoryDocs]) != 2 {
		t.Errorf("Expected 2 docs results, got %d", len(categoryResults[types.CategoryDocs]))
	}

	if len(categoryResults[types.CategoryCommits]) != 1 {
		t.Errorf("Expected 1 commit result, got %d", len(categoryResults[types.CategoryCommits]))
	}
}

func TestGetFailedChecks(t *testing.T) {
	scorer := NewScorer()

	// Add mixed results
	scorer.AddResult(types.CheckResult{
		ID:     "PASS-001",
		Status: types.StatusPass,
	})

	scorer.AddResult(types.CheckResult{
		ID:     "FAIL-001",
		Status: types.StatusFail,
	})

	scorer.AddResult(types.CheckResult{
		ID:     "WARN-001",
		Status: types.StatusWarning,
	})

	failedChecks := scorer.GetFailedChecks()

	if len(failedChecks) != 1 {
		t.Errorf("Expected 1 failed check, got %d", len(failedChecks))
	}

	if failedChecks[0].ID != "FAIL-001" {
		t.Errorf("Expected ID 'FAIL-001', got '%s'", failedChecks[0].ID)
	}
}

func TestGetWarningChecks(t *testing.T) {
	scorer := NewScorer()

	// Add mixed results
	scorer.AddResult(types.CheckResult{
		ID:     "PASS-001",
		Status: types.StatusPass,
	})

	scorer.AddResult(types.CheckResult{
		ID:     "FAIL-001",
		Status: types.StatusFail,
	})

	scorer.AddResult(types.CheckResult{
		ID:     "WARN-001",
		Status: types.StatusWarning,
	})

	warningChecks := scorer.GetWarningChecks()

	if len(warningChecks) != 1 {
		t.Errorf("Expected 1 warning check, got %d", len(warningChecks))
	}

	if warningChecks[0].ID != "WARN-001" {
		t.Errorf("Expected ID 'WARN-001', got '%s'", warningChecks[0].ID)
	}
}

func TestGetNextSteps(t *testing.T) {
	scorer := NewScorer()

	// Add failed checks that should generate next steps
	scorer.AddResult(types.CheckResult{
		ID:     "DOC-101",
		Status: types.StatusFail,
	})

	scorer.AddResult(types.CheckResult{
		ID:     "IG-201",
		Status: types.StatusFail,
	})

	nextSteps := scorer.GetNextSteps()

	if len(nextSteps) == 0 {
		t.Error("Expected next steps, got empty")
	}

	// Check that we have the expected steps
	foundDocStep := false
	foundIgnoreStep := false

	for _, step := range nextSteps {
		if strings.Contains(strings.ToLower(step), "documentation") {
			foundDocStep = true
		}
		if strings.Contains(strings.ToLower(step), "gitignore") {
			foundIgnoreStep = true
		}
	}

	if !foundDocStep {
		t.Error("Expected documentation step not found")
	}

	if !foundIgnoreStep {
		t.Error("Expected gitignore step not found")
	}
}

func TestCalculateGrade(t *testing.T) {
	testCases := []struct {
		score int
		grade string
	}{
		{95, "A+"},
		{90, "A"},
		{85, "A-"},
		{80, "B+"},
		{75, "B"},
		{70, "B-"},
		{65, "C+"},
		{60, "C"},
		{55, "C-"},
		{50, "D"},
		{30, "F"},
	}

	for _, tc := range testCases {
		grade := calculateGrade(tc.score)
		if grade != tc.grade {
			t.Errorf("For score %d, expected grade '%s', got '%s'", tc.score, tc.grade, grade)
		}
	}
}

func TestGetWeightForCategory(t *testing.T) {
	testCases := []struct {
		category types.Category
		weight   int
	}{
		{types.CategoryDocs, 3},
		{types.CategoryCommits, 4},
		{types.CategoryHygiene, 2},
	}

	for _, tc := range testCases {
		weight := getWeightForCategory(tc.category)
		if weight != tc.weight {
			t.Errorf("For category %v, expected weight %d, got %d", tc.category, tc.weight, weight)
		}
	}
}
