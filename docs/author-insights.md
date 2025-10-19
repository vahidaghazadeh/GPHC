# Author Insights Guide

This guide covers analyzing commit authors and contributor patterns for team insights.

## Overview

Author insights help teams understand contributor patterns, identify "bus factor" risks, and ensure healthy team collaboration.

## Basic Usage

### Analyzing Contributors
```bash
# Analyze commit authors
gphc authors

# Example output:
Author Analysis
===============

Total Contributors: 4
Total Commits: 156

Top Contributors:
1. vahidaghazadeh: 118 commits (76%)
2. john.doe: 25 commits (16%)
3. jane.smith: 10 commits (6%)
4. mike.wilson: 3 commits (2%)

⚠️ Single Author Dominance Detected (>70%)
Recommendation: Encourage more team participation
```

### Detailed Analysis
```bash
# Get detailed author insights
gphc authors --detailed

# Example output:
Detailed Author Analysis
========================

Contributors: 4
Total Commits: 156
Analysis Period: Last 90 days

Author Breakdown:
┌─────────────────┬─────────┬─────────┬─────────┬─────────┐
│ Author          │ Commits │ Percent │ Email   │ Status  │
├─────────────────┼─────────┼─────────┼─────────┼─────────┤
│ vahidaghazadeh  │ 118     │ 76%     │ v@...   │ Active  │
│ john.doe        │ 25      │ 16%     │ j@...   │ Active  │
│ jane.smith      │ 10      │ 6%      │ j@...   │ Recent  │
│ mike.wilson     │ 3       │ 2%      │ m@...   │ Inactive│
└─────────────────┴─────────┴─────────┴─────────┴─────────┘

Risk Analysis:
⚠️ Single Author Dominance: vahidaghazadeh (76%)
⚠️ Inactive Contributor: mike.wilson (last commit: 45 days ago)
✅ Email Consistency: All contributors use consistent email format

Recommendations:
1. Encourage more team participation
2. Review inactive contributor status
3. Consider pair programming sessions
4. Implement code review requirements
```

## Configuration

### Author Analysis Settings
```yaml
# gphc.yml
author_insights:
  enabled: true
  
  # Analysis settings
  analysis_period: 90  # days
  min_commits: 1
  max_authors: 50
  
  # Risk thresholds
  thresholds:
    single_author_dominance: 70  # percentage
    inactive_threshold: 30  # days
    email_consistency_check: true
  
  # Output settings
  output:
    show_email: false
    show_percentages: true
    show_recommendations: true
```

### Custom Rules
```yaml
# gphc.yml
author_insights:
  custom_rules:
    - id: "AUTH-CUSTOM-001"
      name: "Team Size Check"
      min_contributors: 3
      score: 5
      
    - id: "AUTH-CUSTOM-002"
      name: "Contributor Diversity"
      max_dominance: 60
      score: 3
      
    - id: "AUTH-CUSTOM-003"
      name: "Active Contributors"
      min_active_contributors: 2
      score: 4
```

## Use Cases

### Team Health Monitoring
```bash
# Monitor team health
gphc authors --period 30

# Check for bus factor risk
gphc authors --check-bus-factor

# Analyze contributor trends
gphc authors --trends
```

### Project Management
```bash
# Generate contributor report
gphc authors --format json --output contributors.json

# Check team balance
gphc authors --check-balance

# Identify inactive contributors
gphc authors --check-inactive
```

### Code Review Analysis
```bash
# Analyze code review patterns
gphc authors --review-analysis

# Check reviewer distribution
gphc authors --reviewer-distribution

# Identify review bottlenecks
gphc authors --review-bottlenecks
```

## Integration Examples

### CI/CD Integration
```yaml
# GitHub Actions
- name: Author Analysis
  run: gphc authors --format json --output authors.json

- name: Check Bus Factor
  run: gphc authors --check-bus-factor --min-contributors 3
```

### Team Reporting
```bash
# Generate weekly team report
gphc authors --period 7 --format markdown --output team-report.md

# Generate monthly contributor analysis
gphc authors --period 30 --format json --output monthly-analysis.json
```

## Best Practices

### For Teams
1. **Regular Analysis**: Analyze contributors regularly
2. **Team Balance**: Maintain balanced contribution levels
3. **Knowledge Sharing**: Encourage knowledge sharing
4. **Code Reviews**: Implement code review requirements
5. **Documentation**: Document team processes

### For Organizations
1. **Team Standards**: Establish team contribution standards
2. **Risk Management**: Monitor bus factor risks
3. **Training Programs**: Implement team training programs
4. **Mentoring**: Establish mentoring programs
5. **Succession Planning**: Plan for team member transitions

## Troubleshooting

### Common Issues
- **No Contributors**: Check repository access
- **Incomplete Data**: Verify commit history
- **Email Inconsistency**: Check email format
- **Analysis Errors**: Validate configuration

### Debugging
```bash
# Test author analysis
gphc authors --test

# Verbose output
gphc authors --verbose

# Check specific period
gphc authors --period 30 --debug
```

## Next Steps
- [Health Checks](health-checks.md) - Understanding health check categories
- [Multi-Repository Scan](multi-repository-scan.md) - Batch repository analysis
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
