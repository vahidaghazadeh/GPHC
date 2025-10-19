# Export Formats Guide

This guide covers exporting health reports in various formats.

## Overview

GPHC supports multiple export formats for health reports, making it easy to integrate with other tools and share results with stakeholders.

## Supported Formats

### JSON Export
```bash
# Export to JSON
gphc check --format json

# Save to file
gphc check --format json --output health-report.json
```

### YAML Export
```bash
# Export to YAML
gphc check --format yaml

# Save to file
gphc check --format yaml --output health-report.yaml
```

### Markdown Export
```bash
# Export to Markdown
gphc check --format markdown

# Save to file
gphc check --format markdown --output health-report.md
```

### HTML Export
```bash
# Export to HTML
gphc check --format html

# Save to file
gphc check --format html --output health-report.html
```

## Format Examples

### JSON Format
```json
{
  "repository": "/path/to/project",
  "timestamp": "2024-01-15T10:30:00Z",
  "overall_score": 85,
  "grade": "B+",
  "status": "PASS",
  "summary": {
    "total_checks": 12,
    "passed_checks": 8,
    "failed_checks": 2,
    "warning_checks": 2
  },
  "categories": {
    "documentation": {
      "score": 90,
      "grade": "A-",
      "checks": [
        {
          "id": "DOC-101",
          "name": "README.md exists",
          "status": "PASS",
          "score": 5,
          "message": "README.md file exists"
        }
      ]
    }
  }
}
```

### Markdown Format
```markdown
# Health Report

**Repository**: /path/to/project  
**Date**: 2024-01-15 10:30:00  
**Overall Score**: 85/100 (B+)  
**Status**: PASS

## Summary

| Metric | Value |
|--------|-------|
| Total Checks | 12 |
| Passed | 8 |
| Failed | 2 |
| Warnings | 2 |

## Categories

### Documentation & Project Structure: 90/100 (A-)
- ✅ README.md exists
- ✅ LICENSE file exists
- ⚠️ CONTRIBUTING.md missing

### Commit History Quality: 85/100 (B+)
- ✅ Conventional commits format
- ✅ Commit message length
- ⚠️ Large commits detected

## Recommendations

1. Create CONTRIBUTING.md file
2. Break down large commits
3. Add more test coverage
```

### HTML Format
```html
<!DOCTYPE html>
<html>
<head>
    <title>Health Report</title>
    <style>
        body { font-family: Arial, sans-serif; }
        .score { font-size: 2em; font-weight: bold; }
        .pass { color: green; }
        .fail { color: red; }
        .warn { color: orange; }
    </style>
</head>
<body>
    <h1>Health Report</h1>
    <div class="score pass">85/100 (B+)</div>
    <p>Repository: /path/to/project</p>
    <p>Date: 2024-01-15 10:30:00</p>
    
    <h2>Summary</h2>
    <ul>
        <li>Total Checks: 12</li>
        <li>Passed: 8</li>
        <li>Failed: 2</li>
        <li>Warnings: 2</li>
    </ul>
</body>
</html>
```

## Integration Examples

### CI/CD Integration
```yaml
# GitHub Actions
- name: Health Check
  run: gphc check --format json --output health-report.json

- name: Upload Health Report
  uses: actions/upload-artifact@v3
  with:
    name: health-report
    path: health-report.json
```

### Documentation Integration
```bash
# Generate markdown for README
gphc check --format markdown >> README.md

# Generate HTML for project website
gphc check --format html --output docs/health-report.html
```

### Team Sharing
```bash
# Generate report for team meeting
gphc check --format markdown --output team-health-report.md

# Generate JSON for dashboard
gphc check --format json --output dashboard-data.json
```

## Configuration

### Export Settings
```yaml
# gphc.yml
export:
  default_format: "json"
  include_timestamps: true
  include_recommendations: true
  
  # Custom templates
  templates:
    markdown:
      include_summary: true
      include_categories: true
      include_recommendations: true
    
    html:
      theme: "dark"
      include_charts: true
      responsive: true
```

### Custom Templates
```yaml
# gphc.yml
export:
  custom_templates:
    - name: "team-report"
      format: "markdown"
      template: |
        # Team Health Report
        
        **Project**: {repository}
        **Score**: {score}/100 ({grade})
        **Status**: {status}
        
        ## Issues
        {issues}
        
        ## Recommendations
        {recommendations}
```

## Best Practices

### For Teams
1. **Consistent Format**: Use consistent export format across team
2. **Regular Reports**: Generate reports regularly for tracking
3. **Share Results**: Share reports with team members
4. **Documentation**: Include reports in project documentation
5. **Trend Analysis**: Track export trends over time

### For Organizations
1. **Standardized Reports**: Use standardized report formats
2. **Dashboard Integration**: Integrate with organizational dashboards
3. **Stakeholder Reports**: Generate reports for stakeholders
4. **Compliance**: Use reports for compliance documentation
5. **Quality Metrics**: Track quality metrics across projects

## Troubleshooting

### Common Issues
- **Format Not Supported**: Check supported formats
- **File Permission**: Check file write permissions
- **Template Errors**: Validate custom templates
- **Large Files**: Optimize for large repositories

### Debugging
```bash
# Test export functionality
gphc check --format json --dry-run

# Validate export format
gphc check --format json --validate

# Check export options
gphc check --help
```

## Next Steps
- [Web Dashboard](web-dashboard.md) - Web server and team collaboration
- [Multi-Repository Scan](multi-repository-scan.md) - Batch repository analysis
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
