# Web Dashboard Guide

This guide covers the web dashboard feature for team collaboration and health monitoring.

## Overview

The GPHC Web Dashboard provides a modern web interface for monitoring repository health, perfect for teams that need centralized health tracking across multiple projects.

## Getting Started

### Starting the Dashboard
```bash
# Start dashboard on default port (8080)
git hc serve

# Start with custom port
git hc serve --port 3000

# Start with custom host and port
git hc serve --host 0.0.0.0 --port 8080

# Start with authentication
git hc serve --auth --username admin --password secret

# Start with custom title
git hc serve --title "My Project Dashboard"
```

### Accessing the Dashboard
Once started, open your browser and navigate to:
- **Default**: http://localhost:8080
- **Custom Port**: http://localhost:YOUR_PORT
- **Custom Host**: http://YOUR_HOST:YOUR_PORT

## Dashboard Features

### Health Overview
The dashboard displays:
- **Overall Health Score**: Current repository health score
- **Grade**: Letter grade (A+, A, A-, B+, etc.)
- **Check Summary**: Total checks, passed, failed, warnings
- **Last Updated**: Timestamp of last health check
- **Repository Name**: Current repository being monitored

### Real-time Updates
- **Auto-refresh**: Dashboard updates every 30 seconds
- **Manual Refresh**: Click "Refresh" button for immediate update
- **Live Data**: Health data is fetched from the repository in real-time

### Export Options
- **JSON Export**: Download health report as JSON
- **PDF Export**: Download health report as PDF (placeholder)
- **API Access**: Direct API access for integration

## Server Configuration

### Command Line Options
```bash
# Server configuration
git hc serve --host localhost --port 8080

# Authentication
git hc serve --auth --username admin --password secret

# CORS settings
git hc serve --cors  # Enable CORS (default: true)
git hc serve --no-cors  # Disable CORS

# Dashboard customization
git hc serve --title "My Custom Dashboard"
```

### Configuration File
Create a `git-hc.yml` file for persistent configuration:

```yaml
# Server configuration
server:
  port: 8080
  host: "localhost"
  auth:
    enabled: true
    username: "admin"
    password: "secret"
  cors:
    enabled: true
    origins: ["http://localhost:3000"]

# Dashboard settings
dashboard:
  title: "Project Health Monitor"
  theme: "dark"  # dark, light, auto
  refresh_interval: "30s"
  auto_refresh: true
  show_timestamps: true
  max_projects: 50
  default_view: "overview"  # overview, trends, details
```

## API Endpoints

### Health Data
```bash
# Get current health data
GET /api/health

# Response format:
{
  "overall_score": 85,
  "grade": "B+",
  "summary": {
    "total_checks": 12,
    "passed_checks": 8,
    "failed_checks": 2,
    "warning_checks": 2
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "repository": "project-name"
}
```

### Export Endpoints
```bash
# Export to JSON
GET /api/export/json

# Export to PDF (placeholder)
GET /api/export/pdf

# Export trends (future)
GET /api/export/trends.csv
```

### CORS Support
The dashboard supports CORS for integration with other applications:
```bash
# CORS headers are automatically added when enabled
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

## Team Collaboration

### Shared Dashboard
Perfect for team environments:
```bash
# Start server accessible to team
git hc serve --host 0.0.0.0 --port 8080

# Team members can access:
# http://your-server:8080
```

### Authentication
Enable basic authentication for team access:
```bash
# Start with authentication
git hc serve --auth --username team --password secure123

# Access dashboard with credentials
# Username: team
# Password: secure123
```

### Integration with CI/CD
Use the API endpoints in your CI/CD pipeline:
```yaml
# GitHub Actions example
- name: Health Check
  run: |
    curl -f http://localhost:8080/api/health || exit 1
    
- name: Export Health Report
  run: |
    curl -o health-report.json http://localhost:8080/api/export/json
```

## Dashboard Interface

### Main Dashboard
```
┌─────────────────────────────────────────────────────────────┐
│ GPHC Dashboard - Project Health Monitor                     │
│ Last Updated: 2024-01-15 10:30:00                          │
├─────────────────────────────────────────────────────────────┤
│ Health Overview                                            │
│                                                             │
│ Overall Score: 85/100 (B+)                                 │
│ Status: PASS                                                │
│                                                             │
│ Total Checks: 12                                           │
│ Passed: 8                                                   │
│ Failed: 2                                                   │
│ Warnings: 2                                                 │
│                                                             │
│ [Refresh] [Export JSON] [Export PDF]                        │
└─────────────────────────────────────────────────────────────┘
```

### Responsive Design
The dashboard is fully responsive and works on:
- **Desktop**: Full-featured interface
- **Tablet**: Optimized layout
- **Mobile**: Touch-friendly interface

## Advanced Features

### Multi-Project Monitoring
```bash
# Scan multiple repositories
git hc scan ~/projects --recursive

# Start dashboard for multiple projects
git hc serve --multi-project ~/projects
```

### Custom Themes
```yaml
# git-hc.yml
dashboard:
  theme: "dark"  # dark, light, auto
  colors:
    score_excellent: "#00ff00"  # Green
    score_good: "#ffff00"       # Yellow
    score_poor: "#ff0000"      # Red
    status_pass: "#00ff00"
    status_fail: "#ff0000"
    status_warn: "#ffaa00"
```

### Notifications
```bash
# Enable Slack notifications
git hc serve --notifications slack --webhook-url https://hooks.slack.com/...

# Enable email notifications
git hc serve --notifications email --smtp-server smtp.company.com
```

## Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Use a different port
git hc serve --port 8081

# Check what's using the port
lsof -i :8080
```

#### Permission Denied
```bash
# Use a different host
git hc serve --host 127.0.0.1

# Check repository permissions
ls -la /path/to/repository
```

#### Dashboard Not Loading
```bash
# Check server status
curl http://localhost:8080/api/health

# Check server logs
git hc serve --verbose
```

### Performance Issues
```bash
# Large number of projects
git hc serve --max-projects 100 --cache-ttl 300s

# Slow network
git hc serve --compression --cache-ttl 600s
```

## Security Considerations

### Authentication
- Use strong passwords for authentication
- Consider using environment variables for credentials
- Enable HTTPS in production environments

### Network Security
- Use firewall rules to restrict access
- Consider VPN access for remote teams
- Monitor access logs for suspicious activity

### Data Privacy
- Health data is processed locally
- No data is sent to external services
- Consider data retention policies

## Best Practices

### For Teams
1. **Centralized Monitoring**: Use dashboard for team-wide health tracking
2. **Regular Reviews**: Schedule weekly health check reviews
3. **Integration**: Integrate with existing team tools
4. **Documentation**: Document dashboard usage for team members
5. **Maintenance**: Regular server maintenance and updates

### For Organizations
1. **Multi-Project View**: Monitor all repositories simultaneously
2. **Trend Analysis**: Track health improvements over time
3. **Reporting**: Generate reports for stakeholders
4. **Compliance**: Ensure health standards across projects
5. **Training**: Train teams on health monitoring practices

## Next Steps

- [Multi-Repository Scan](multi-repository-scan.md) - Batch repository analysis
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
- [Notifications](notifications.md) - Slack and webhook setup
- [Terminal UI](terminal-ui.md) - Interactive terminal interface
