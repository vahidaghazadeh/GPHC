# Notifications Guide

This guide covers setting up Slack and webhook notifications for health monitoring.

## Overview

Notifications allow you to send health reports directly to team channels, keeping everyone informed about repository health status.

## Slack Integration

### Basic Setup
```bash
# Send health report to Slack
git hc check --notify slack --webhook-url https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

### Slack Configuration
```yaml
# gphc.yml
notifications:
  slack:
    enabled: true
    webhook_url: "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
    channel: "#dev-team"
    username: "GPHC Bot"
    icon_emoji: ":robot_face:"
```

### Slack Message Format
```
🧭 GPHC Health Report: 83/100 (B+)
Repository: project-name
Status: PASS

📊 Health Breakdown:
• Documentation: 90/100 (A-)
• Commits: 85/100 (B+)
• Hygiene: 80/100 (B-)
• Structure: 75/100 (C+)

⚠️ Issues Found:
• Missing CONTRIBUTING.md
• 2 stale branches found

🔗 View Details: http://localhost:8080
```

## Discord Integration

### Basic Setup
```bash
# Send health report to Discord
git hc check --notify discord --webhook-url https://discord.com/api/webhooks/YOUR/WEBHOOK/URL
```

### Discord Configuration
```yaml
# gphc.yml
notifications:
  discord:
    enabled: true
    webhook_url: "https://discord.com/api/webhooks/YOUR/WEBHOOK/URL"
    username: "GPHC Bot"
    avatar_url: "https://example.com/gphc-avatar.png"
```

## Custom Webhooks

### Generic Webhook
```bash
# Send to custom webhook
git hc check --notify webhook --webhook-url https://your-service.com/webhook
```

### Webhook Payload
```json
{
  "repository": "project-name",
  "timestamp": "2024-01-15T10:30:00Z",
  "overall_score": 83,
  "grade": "B+",
  "status": "PASS",
  "categories": {
    "documentation": 90,
    "commits": 85,
    "hygiene": 80,
    "structure": 75
  },
  "issues": [
    {
      "id": "DOC-101",
      "name": "Missing CONTRIBUTING.md",
      "status": "FAIL",
      "message": "CONTRIBUTING.md file is missing"
    }
  ]
}
```

## Configuration

### Notification Settings
```yaml
# gphc.yml
notifications:
  enabled: true
  
  # Slack configuration
  slack:
    enabled: true
    webhook_url: "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
    channel: "#dev-team"
    username: "GPHC Bot"
    icon_emoji: ":robot_face:"
    
  # Discord configuration
  discord:
    enabled: true
    webhook_url: "https://discord.com/api/webhooks/YOUR/WEBHOOK/URL"
    username: "GPHC Bot"
    
  # Custom webhook
  webhook:
    enabled: true
    url: "https://your-service.com/webhook"
    headers:
      Authorization: "Bearer YOUR_TOKEN"
      Content-Type: "application/json"
```

### Conditional Notifications
```yaml
# gphc.yml
notifications:
  conditions:
    - trigger: "score_below"
      threshold: 70
      message: "Health score dropped below 70!"
    
    - trigger: "new_issues"
      message: "New health issues detected"
    
    - trigger: "score_improved"
      threshold: 10
      message: "Health score improved by 10+ points!"
```

## Integration Examples

### GitHub Actions
```yaml
name: Health Check with Notifications
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  health-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - name: Install GPHC
        run: go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
      - name: Setup Git HC
        run: ./setup-git-hc.sh
      - name: Health Check with Notifications
        run: git hc check --notify slack --webhook-url ${{ secrets.SLACK_WEBHOOK_URL }}
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
```

### GitLab CI
```yaml
health_check:
  stage: test
  script:
    - go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
    - ./setup-git-hc.sh
    - git hc check --notify slack --webhook-url $SLACK_WEBHOOK_URL
  variables:
    SLACK_WEBHOOK_URL: $SLACK_WEBHOOK_URL
```

## Advanced Features

### Message Templates
```yaml
# gphc.yml
notifications:
  templates:
    slack:
      success: |
        ✅ Health Check Passed: {score}/100 ({grade})
        Repository: {repository}
        Status: {status}
      
      failure: |
        ❌ Health Check Failed: {score}/100 ({grade})
        Repository: {repository}
        Issues: {issue_count}
      
      warning: |
        ⚠️ Health Check Warning: {score}/100 ({grade})
        Repository: {repository}
        Warnings: {warning_count}
```

### Scheduled Notifications
```bash
# Send daily health report
git hc check --schedule daily --notify slack

# Send weekly health report
git hc check --schedule weekly --notify slack

# Send monthly health report
git hc check --schedule monthly --notify slack
```

## Troubleshooting

### Common Issues
- **Webhook URL Invalid**: Check webhook URL format
- **Authentication Failed**: Verify webhook credentials
- **Message Not Sent**: Check network connectivity
- **Rate Limiting**: Implement rate limiting for frequent notifications

### Debugging
```bash
# Test webhook connection
git hc check --test-webhook --webhook-url https://hooks.slack.com/services/YOUR/WEBHOOK/URL

# Verbose notification output
git hc check --notify slack --verbose

# Dry run notifications
git hc check --notify slack --dry-run
```

## Best Practices

### For Teams
1. **Appropriate Channels**: Use dedicated channels for health notifications
2. **Clear Messages**: Keep messages concise and actionable
3. **Regular Updates**: Send notifications at appropriate intervals
4. **Team Training**: Train team on notification meanings
5. **Feedback Loop**: Gather feedback on notification usefulness

### For Organizations
1. **Standardized Messages**: Use consistent message formats
2. **Escalation Procedures**: Define escalation for critical issues
3. **Notification Policies**: Establish notification policies
4. **Monitoring**: Monitor notification effectiveness
5. **Continuous Improvement**: Improve notification systems

## Next Steps
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
- [Web Dashboard](web-dashboard.md) - Web server and team collaboration
- [Multi-Repository Scan](multi-repository-scan.md) - Batch repository analysis
