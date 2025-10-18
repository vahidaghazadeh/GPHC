package exporter

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/opsource/gphc/pkg/types"
	"gopkg.in/yaml.v2"
)

// ExportFormat represents the supported export formats
type ExportFormat string

const (
	FormatJSON     ExportFormat = "json"
	FormatYAML     ExportFormat = "yaml"
	FormatMarkdown ExportFormat = "markdown"
	FormatHTML     ExportFormat = "html"
)

// Exporter handles exporting health reports in different formats
type Exporter struct{}

// NewExporter creates a new exporter instance
func NewExporter() *Exporter {
	return &Exporter{}
}

// Export exports the health report in the specified format
func (e *Exporter) Export(report *types.HealthReport, format ExportFormat) (string, error) {
	switch format {
	case FormatJSON:
		return e.exportJSON(report)
	case FormatYAML:
		return e.exportYAML(report)
	case FormatMarkdown:
		return e.exportMarkdown(report)
	case FormatHTML:
		return e.exportHTML(report)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// exportJSON exports the report as JSON
func (e *Exporter) exportJSON(report *types.HealthReport) (string, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// exportYAML exports the report as YAML
func (e *Exporter) exportYAML(report *types.HealthReport) (string, error) {
	data, err := yaml.Marshal(report)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// exportMarkdown exports the report as Markdown
func (e *Exporter) exportMarkdown(report *types.HealthReport) (string, error) {
	var output strings.Builder

	// Header
	output.WriteString("# Repository Health Report\n\n")
	output.WriteString(fmt.Sprintf("**Overall Health Score:** %d/100 (%s)\n\n", report.OverallScore, report.Grade))

	// Summary
	output.WriteString("## Summary\n\n")
	output.WriteString(fmt.Sprintf("- **Total Checks:** %d\n", report.Summary.TotalChecks))
	output.WriteString(fmt.Sprintf("- **Passed:** %d\n", report.Summary.PassedChecks))
	output.WriteString(fmt.Sprintf("- **Failed:** %d\n", report.Summary.FailedChecks))
	output.WriteString(fmt.Sprintf("- **Warnings:** %d\n\n", report.Summary.WarningChecks))

	// Results by category
	categories := make(map[string][]types.CheckResult)
	for _, result := range report.Results {
		categoryName := result.Category.String()
		categories[categoryName] = append(categories[categoryName], result)
	}

	for category, results := range categories {
		output.WriteString(fmt.Sprintf("## %s\n\n", category))

		for _, result := range results {
			status := "✅"
			if result.Status == types.StatusFail {
				status = "❌"
			} else if result.Status == types.StatusWarning {
				status = "⚠️"
			}

			output.WriteString(fmt.Sprintf("### %s %s\n\n", status, result.ID))
			output.WriteString(fmt.Sprintf("**Message:** %s\n\n", result.Message))
			output.WriteString(fmt.Sprintf("**Score:** %+d\n\n", result.Score))

			if len(result.Details) > 0 {
				output.WriteString("**Details:**\n")
				for _, detail := range result.Details {
					output.WriteString(fmt.Sprintf("- %s\n", detail))
				}
				output.WriteString("\n")
			}
		}
	}

	// Next Steps (generate from failed checks)
	failedChecks := make([]string, 0)
	for _, result := range report.Results {
		if result.Status == types.StatusFail {
			failedChecks = append(failedChecks, fmt.Sprintf("Fix %s: %s", result.ID, result.Message))
		}
	}

	if len(failedChecks) > 0 {
		output.WriteString("## Next Steps\n\n")
		for i, step := range failedChecks {
			output.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
		}
	}

	return output.String(), nil
}

// exportHTML exports the report as HTML
func (e *Exporter) exportHTML(report *types.HealthReport) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Repository Health Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .score { font-size: 2em; font-weight: bold; color: #2c3e50; }
        .summary { display: flex; justify-content: space-around; margin: 20px 0; }
        .summary-item { text-align: center; padding: 15px; background-color: #ecf0f1; border-radius: 5px; }
        .category { margin: 20px 0; }
        .result { margin: 15px 0; padding: 15px; border-left: 4px solid #3498db; background-color: #f8f9fa; }
        .result.fail { border-left-color: #e74c3c; }
        .result.warning { border-left-color: #f39c12; }
        .result.pass { border-left-color: #27ae60; }
        .next-steps { background-color: #e8f4fd; padding: 20px; border-radius: 5px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Repository Health Report</h1>
            <div class="score">{{.OverallScore}}/100 ({{.Grade}})</div>
        </div>
        
        <div class="summary">
            <div class="summary-item">
                <h3>{{.Summary.TotalChecks}}</h3>
                <p>Total Checks</p>
            </div>
            <div class="summary-item">
                <h3>{{.Summary.PassedChecks}}</h3>
                <p>Passed</p>
            </div>
            <div class="summary-item">
                <h3>{{.Summary.FailedChecks}}</h3>
                <p>Failed</p>
            </div>
            <div class="summary-item">
                <h3>{{.Summary.WarningChecks}}</h3>
                <p>Warnings</p>
            </div>
        </div>
        
        {{range .Results}}
        <div class="result {{.Status}}">
            <h3>{{.ID}}: {{.Message}}</h3>
            <p><strong>Score:</strong> {{.Score}}</p>
            {{if .Details}}
            <ul>
                {{range .Details}}
                <li>{{.}}</li>
                {{end}}
            </ul>
            {{end}}
        </div>
        {{end}}
    </div>
</body>
</html>`

	t, err := template.New("html").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var output strings.Builder
	err = t.Execute(&output, report)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateBadgeURL generates a badge URL for the health score
func (e *Exporter) GenerateBadgeURL(score int) string {
	var color string
	switch {
	case score >= 90:
		color = "brightgreen"
	case score >= 80:
		color = "green"
	case score >= 70:
		color = "yellowgreen"
	case score >= 60:
		color = "yellow"
	case score >= 50:
		color = "orange"
	default:
		color = "red"
	}

	return fmt.Sprintf("https://img.shields.io/badge/Health_Score-%d%%2F100-%s?style=for-the-badge&logo=github", score, color)
}

// GenerateMarkdownBadge generates a markdown badge for the health score
func (e *Exporter) GenerateMarkdownBadge(score int) string {
	badgeURL := e.GenerateBadgeURL(score)
	return fmt.Sprintf("![Health Score](%s)", badgeURL)
}
