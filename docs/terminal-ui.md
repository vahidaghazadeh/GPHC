# Terminal UI (TUI) Guide

This guide covers the interactive terminal user interface for health monitoring.

## Overview

The GPHC Terminal UI (TUI) provides a beautiful, interactive terminal-based user interface that makes health checking and monitoring an engaging experience for developers and technical teams.

## Getting Started

### Launching the TUI
```bash
# Start the TUI
git hc tui

# Start TUI with specific repository
git hc tui /path/to/repository

# Start TUI with auto-refresh
git hc tui --refresh 30s

# Start TUI in full-screen mode
git hc tui --fullscreen
```

### Interface Overview
The TUI provides multiple views and interactive features:

```
┌─────────────────────────────────────────────────────────────┐
│ GPHC - Git Project Health Checker                          │
│ Repository: /path/to/project                                │
│ Last Updated: 2024-01-15 10:30:00                          │
├─────────────────────────────────────────────────────────────┤
│ Overall Health Score: 85/100 (B+)                          │
│ Status: PASS                                                │
│                                                             │
│ ┌─ Health Overview ──────────────────────────────────────┐ │
│ │ Documentation & Project Structure: 90/100 (A-)         │ │
│ │ Commit History Quality: 85/100 (B+)                    │ │
│ │ Git Cleanup & Hygiene: 80/100 (B-)                    │ │
│ │ Codebase Structure: 75/100 (C+)                        │ │
│ └─────────────────────────────────────────────────────────┘ │
│                                                             │
│ ┌─ Quick Actions ────────────────────────────────────────┐ │
│ │ [F1] Help  [F2] Filter  [F3] Trends  [F4] Settings    │ │
│ │ [F5] Refresh  [F6] Export  [F7] Notify  [F8] Quit      │ │
│ └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Main Views

### 1. Health Overview Tab
Displays comprehensive health information:
- **Overall Score**: Current health score with grade
- **Category Breakdown**: Scores for each health category
- **Check Summary**: Total checks, passed, failed, warnings
- **Repository Info**: Repository name and last update time

### 2. Detailed Check Results Tab
Shows individual check results:
- **Check Status**: PASS, WARN, or FAIL for each check
- **Check Details**: Detailed messages and recommendations
- **Filtering Options**: Filter by status, category, or score
- **Interactive Navigation**: Use arrow keys to navigate

### 3. Trend Analysis Tab
Historical health tracking:
- **Score Trends**: Health score changes over time
- **Improvement Tracking**: Visual representation of progress
- **Historical Data**: Access to past health reports
- **Trend Analysis**: Identify patterns and improvements

## Interactive Features

### Keyboard Shortcuts
```
Navigation:
  ↑/↓/←/→    Move cursor
  Tab        Switch between panels
  Enter      Select/expand item
  Esc        Go back/close dialog

Actions:
  F1         Help and shortcuts
  F2         Open filter menu
  F3         View trend analysis
  F4         Settings and configuration
  F5         Refresh data
  F6         Export results
  F7         Send notifications
  F8         Quit application

Search and Filter:
  /          Search checks
  Ctrl+F     Advanced filtering
  Ctrl+R     Reset filters
  Space      Toggle selection
```

### Real-time Updates
```bash
# Auto-refresh every 30 seconds
git hc tui --refresh 30s

# Auto-refresh every 2 minutes
git hc tui --refresh 2m

# Watch mode - refresh on file changes
git hc tui --watch
```

### Rule Explanations
When you select a check, view detailed explanations:

```
┌─ Rule Details: DOC-101 ────────────────────────────────────┐
│ Name: README.md exists                                     │
│ ID: DOC-101                                               │
│ Category: Documentation & Project Structure               │
│ Score: 5 points                                           │
│                                                           │
│ Description:                                              │
│ Checks if the project has a README.md file in the root   │
│ directory. This file is essential for project             │
│ documentation and helps new contributors understand        │
│ the project.                                              │
│                                                           │
│ Requirements:                                             │
│ • README.md file must exist in repository root           │
│ • File should not be empty                               │
│ • Should contain project description                      │
│                                                           │
│ Benefits:                                                 │
│ • Improves project discoverability                       │
│ • Helps new contributors get started                     │
│ • Provides essential project information                  │
│                                                           │
│ [Press Enter to go back]                                 │
└───────────────────────────────────────────────────────────┘
```

## Advanced Features

### Multi-Repository View
```bash
# Compare multiple repositories
git hc tui --multi-repo ~/projects/*

# Scan and display all repositories
git hc tui --scan-recursive ~/projects
```

### Custom Themes
```yaml
# gphc.yml
tui:
  theme: "dark"  # dark, light, auto
  colors:
    score_excellent: "#00ff00"  # Green
    score_good: "#ffff00"       # Yellow
    score_poor: "#ff0000"      # Red
    status_pass: "#00ff00"
    status_fail: "#ff0000"
    status_warn: "#ffaa00"
  
  layout:
    show_timestamps: true
    show_categories: true
    show_recommendations: true
    compact_mode: false
```

### Filtering and Search
- **Status Filter**: Show only PASS, WARN, or FAIL checks
- **Category Filter**: Filter by health check category
- **Score Range**: Filter by score range (e.g., 0-50, 50-80, 80-100)
- **Text Search**: Search for specific checks or messages

### Export Options
```bash
# Export current view
F6 -> Export -> JSON/YAML/Markdown/HTML

# Export with filters applied
F6 -> Export Filtered -> JSON

# Export trends
F3 -> Trends -> Export -> CSV
```

## Configuration

### TUI Settings
```yaml
# gphc.yml
tui:
  refresh_interval: "30s"
  theme: "dark"
  fullscreen: false
  auto_refresh: true
  show_timestamps: true
  show_categories: true
  show_recommendations: true
  compact_mode: false
  
  # Keyboard shortcuts
  shortcuts:
    help: "F1"
    filter: "F2"
    trends: "F3"
    settings: "F4"
    refresh: "F5"
    export: "F6"
    notify: "F7"
    quit: "F8"
```

### Display Options
```bash
# Compact mode for smaller terminals
git hc tui --compact-mode

# Full-screen mode
git hc tui --fullscreen

# No color mode
git hc tui --no-colors

# Custom refresh interval
git hc tui --refresh 2m
```

## Use Cases

### Daily Development
```bash
# Quick health check during development
git hc tui --refresh 2m

# Check specific repository
git hc tui /path/to/current/project
```

### Code Reviews
```bash
# Review health before code review
git hc tui --filter failed

# Check trends over time
git hc tui --trends
```

### Team Standups
```bash
# Display health for team standup
git hc tui --multi-repo ~/team/projects/*

# Show only critical issues
git hc tui --filter failed --min-score 70
```

### Project Management
```bash
# Monitor project health trends
git hc tui --trends --refresh 5m

# Export health report for stakeholders
git hc tui --export json
```

## Troubleshooting

### Common Issues

#### Terminal Compatibility
```bash
# Use compatibility mode for older terminals
git hc tui --compatibility-mode

# Check terminal capabilities
echo $TERM
```

#### Performance Issues
```bash
# Use fast mode for large repositories
git hc tui --fast-mode

# Disable auto-refresh for better performance
git hc tui --no-refresh
```

#### Display Issues
```bash
# Use no color mode for problematic terminals
git hc tui --no-colors

# Use compact mode for small terminals
git hc tui --compact-mode
```

### Configuration Issues
```bash
# Reset configuration
rm ~/.gphc.yml
git hc tui

# Check configuration
git hc tui --config-check
```

## Best Practices

### For Developers
1. **Regular Monitoring**: Use TUI for daily health checks
2. **Quick Access**: Keep TUI running during development
3. **Filter Usage**: Use filters to focus on specific issues
4. **Export Reports**: Export health reports for documentation
5. **Trend Tracking**: Monitor health trends over time

### For Teams
1. **Shared Standards**: Use TUI to maintain team standards
2. **Code Reviews**: Check health before code reviews
3. **Standup Integration**: Use TUI for team standup meetings
4. **Training**: Train team members on TUI usage
5. **Documentation**: Document TUI usage for team members

## Integration

### With Other Tools
```bash
# Integrate with Git hooks
git hc tui --pre-commit

# Integrate with CI/CD
git hc tui --ci-mode

# Integrate with IDE
git hc tui --ide-integration
```

### With Team Workflows
```bash
# Team health monitoring
git hc tui --team-mode

# Project health tracking
git hc tui --project-tracking

# Quality assurance
git hc tui --qa-mode
```

## Next Steps

- [Web Dashboard](web-dashboard.md) - Web server and team collaboration
- [Multi-Repository Scan](multi-repository-scan.md) - Batch repository analysis
- [Historical Tracking](historical-tracking.md) - Health trend analysis
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
