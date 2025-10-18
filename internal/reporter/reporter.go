package reporter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opsource/gphc/pkg/types"
)

// Reporter handles the colorful terminal output
type Reporter struct {
	style Style
}

// Style contains all the styling definitions
type Style struct {
	Header      lipgloss.Style
	Title       lipgloss.Style
	Score       lipgloss.Style
	Grade       lipgloss.Style
	Category    lipgloss.Style
	Pass        lipgloss.Style
	Fail        lipgloss.Style
	Warning     lipgloss.Style
	Detail      lipgloss.Style
	Separator   lipgloss.Style
	NextSteps   lipgloss.Style
}

// NewReporter creates a new reporter with default styling
func NewReporter() *Reporter {
	return &Reporter{
		style: createDefaultStyle(),
	}
}

// createDefaultStyle creates the default styling
func createDefaultStyle() Style {
	return Style{
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Margin(1, 0),
		
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Margin(0, 0, 1, 0),
		
		Score: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			Margin(0, 0, 1, 0),
		
		Grade: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			Margin(0, 0, 1, 0),
		
		Category: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEEB")).
			Bold(true).
			Margin(1, 0, 0, 0),
		
		Pass: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),
		
		Fail: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),
		
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true),
		
		Detail: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Margin(0, 0, 0, 3),
		
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Margin(0, 0, 1, 0),
		
		NextSteps: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			Margin(1, 0, 0, 0),
	}
}

// Report generates the full health report output
func (r *Reporter) Report(report *types.HealthReport) string {
	var output strings.Builder

	// Header
	output.WriteString(r.style.Header.Render("✅ Repository Health Check (GPHC v1.0.0)"))
	output.WriteString("\n\n")

	// Overall Score
	scoreText := fmt.Sprintf("🌟 Overall Health Score: %d/100 (%s)", report.OverallScore, report.Grade)
	output.WriteString(r.style.Score.Render(scoreText))
	output.WriteString("\n\n")

	// Separator
	output.WriteString(r.style.Separator.Render(strings.Repeat("-", 50)))
	output.WriteString("\n")

	// Category results
	categoryResults := r.groupResultsByCategory(report.Results)
	for category, results := range categoryResults {
		output.WriteString(r.renderCategory(category, results))
		output.WriteString("\n")
	}

	// Next Steps
	if report.Summary.FailedChecks > 0 || report.Summary.WarningChecks > 0 {
		output.WriteString(r.renderNextSteps(report))
	}

	return output.String()
}

// groupResultsByCategory groups results by category
func (r *Reporter) groupResultsByCategory(results []types.CheckResult) map[types.Category][]types.CheckResult {
	categories := make(map[types.Category][]types.CheckResult)
	
	for _, result := range results {
		categories[result.Category] = append(categories[result.Category], result)
	}
	
	return categories
}

// renderCategory renders a category section
func (r *Reporter) renderCategory(category types.Category, results []types.CheckResult) string {
	var output strings.Builder
	
	// Category header
	passed := 0
	for _, result := range results {
		if result.Status == types.StatusPass {
			passed++
		}
	}
	
	categoryHeader := fmt.Sprintf("[%s] %s (Passed: %d/%d)",
		getCategoryLetter(category),
		category.String(),
		passed,
		len(results))
	
	output.WriteString(r.style.Category.Render(categoryHeader))
	output.WriteString("\n")
	
	// Separator
	output.WriteString(r.style.Separator.Render(strings.Repeat("-", 50)))
	output.WriteString("\n")
	
	// Results
	for _, result := range results {
		output.WriteString(r.renderResult(result))
		output.WriteString("\n")
	}
	
	return output.String()
}

// renderResult renders a single check result
func (r *Reporter) renderResult(result types.CheckResult) string {
	var output strings.Builder
	
	// Status icon and ID
	var statusIcon string
	var statusStyle lipgloss.Style
	
	switch result.Status {
	case types.StatusPass:
		statusIcon = "✅"
		statusStyle = r.style.Pass
	case types.StatusFail:
		statusIcon = "❌"
		statusStyle = r.style.Fail
	case types.StatusWarning:
		statusIcon = "⚠️"
		statusStyle = r.style.Warning
	}
	
	statusText := fmt.Sprintf("%s %s: %s (Score: %+d)",
		statusIcon,
		result.ID,
		result.Message,
		result.Score)
	
	output.WriteString(statusStyle.Render(statusText))
	output.WriteString("\n")
	
	// Details
	for _, detail := range result.Details {
		output.WriteString(r.style.Detail.Render(detail))
		output.WriteString("\n")
	}
	
	return output.String()
}

// renderNextSteps renders the next steps section
func (r *Reporter) renderNextSteps(report *types.HealthReport) string {
	var output strings.Builder
	
	output.WriteString(r.style.NextSteps.Render("💡 Next Steps:"))
	output.WriteString("\n")
	
	// Generate next steps based on failed checks
	nextSteps := r.generateNextSteps(report.Results)
	
	for i, step := range nextSteps {
		output.WriteString(fmt.Sprintf("   %d. %s\n", i+1, step))
	}
	
	return output.String()
}

// generateNextSteps generates actionable next steps
func (r *Reporter) generateNextSteps(results []types.CheckResult) []string {
	var steps []string
	
	for _, result := range results {
		if result.Status == types.StatusFail {
			switch result.ID {
			case "DOC-101":
				if !containsString(steps, "Create missing documentation files") {
					steps = append(steps, "Create missing documentation files (README.md, LICENSE, etc.)")
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

// getCategoryLetter returns the letter for a category
func getCategoryLetter(category types.Category) string {
	switch category {
	case types.CategoryDocs:
		return "A"
	case types.CategoryCommits:
		return "B"
	case types.CategoryHygiene:
		return "C"
	default:
		return "?"
	}
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
