# Historical Tracking Guide

This guide covers tracking repository health over time with trend analysis.

## Overview

Historical tracking allows you to monitor how your repository health improves over time, providing valuable insights for teams and project managers.

## Basic Usage

### Viewing Trends
```bash
# View health trends
gphc trend

# View trends for specific time period
gphc trend --days 30

# View trends with detailed analysis
gphc trend --detailed
```

### Example Output
```
Health Trend Analysis
====================

Repository: /path/to/project
Period: Last 30 days

Score Progression:
  Jan 01: 72/100 (C+)
  Jan 08: 78/100 (C+)
  Jan 15: 85/100 (B+)
  Jan 22: 88/100 (B+)
  Jan 29: 92/100 (A-)

Trend: Improving (+20 points)
Average: 83.0/100
Best Score: 92/100 (Jan 29)
Worst Score: 72/100 (Jan 01)
```

## Configuration

### Historical Settings
```yaml
# gphc.yml
historical:
  enabled: true
  save_interval: "daily"
  retention_days: 365
  
  # Trend analysis
  trend_analysis:
    enabled: true
    min_data_points: 7
    alert_threshold: -10  # Alert if score drops by 10 points
```

## Data Storage

### History File
Health data is stored in `.gphc-history.json`:
```json
{
  "repository": "/path/to/project",
  "history": [
    {
      "timestamp": "2024-01-15T10:30:00Z",
      "score": 85,
      "grade": "B+",
      "categories": {
        "documentation": 90,
        "commits": 85,
        "hygiene": 80,
        "structure": 75
      }
    }
  ]
}
```

## Integration

### CI/CD Integration
```yaml
# GitHub Actions
- name: Health Check
  run: gphc check

- name: Save History
  run: gphc trend --save
```

## Next Steps
- [Multi-Repository Scan](multi-repository-scan.md) - Batch repository analysis
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
