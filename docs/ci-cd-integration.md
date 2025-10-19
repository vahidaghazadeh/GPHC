# CI/CD Integration Guide

This guide covers integrating GPHC with CI/CD pipelines for automated quality gates.

## Overview

CI/CD integration allows you to automatically check repository health in your build pipelines, ensuring quality standards are maintained before code is merged.

## GitHub Actions

### Basic Integration
```yaml
name: Health Check
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
      - name: Health Check
        run: gphc check --min-score 80
```

### Advanced Integration
```yaml
name: Comprehensive Health Check
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
      - name: Health Check
        run: gphc check --min-score 80 --format json
      - name: Upload Health Report
        uses: actions/upload-artifact@v3
        with:
          name: health-report
          path: health-report.json
      - name: Comment PR
        if: github.event_name == 'pull_request'
        run: |
          gphc check --format markdown > health-report.md
          gh pr comment ${{ github.event.pull_request.number }} --body-file health-report.md
```

## GitLab CI

### Basic Integration
```yaml
health_check:
  stage: test
  script:
    - go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
    - gphc check --min-score 80
  artifacts:
    reports:
      junit: health-report.xml
```

### Advanced Integration
```yaml
health_check:
  stage: test
  script:
    - go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
    - gphc check --min-score 80 --format json
    - gphc check --format markdown > health-report.md
  artifacts:
    reports:
      junit: health-report.xml
    paths:
      - health-report.md
    expire_in: 1 week
```

## Jenkins Pipeline

### Basic Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Health Check') {
            steps {
                sh 'go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest'
                sh 'gphc check --min-score 80'
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: 'health-report.json'
        }
    }
}
```

### Advanced Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Health Check') {
            steps {
                sh 'go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest'
                sh 'gphc check --min-score 80 --format json'
                sh 'gphc check --format markdown > health-report.md'
            }
        }
        stage('Quality Gate') {
            steps {
                script {
                    def healthScore = sh(
                        script: 'gphc check --format json | jq -r ".overall_score"',
                        returnStdout: true
                    ).trim()
                    
                    if (healthScore.toInteger() < 80) {
                        error "Health score ${healthScore} is below threshold of 80"
                    }
                }
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: 'health-report.*'
        }
    }
}
```

## Quality Gates

### Setting Thresholds
```bash
# Fail if score is below 80
gphc check --min-score 80

# Fail if score is below 70
gphc check --min-score 70

# Fail on warnings
gphc check --fail-on-warnings
```

### Exit Codes
- **0**: Health check passed
- **1**: Health check failed (score below threshold)
- **2**: Error occurred during check

## Configuration

### CI/CD Settings
```yaml
# gphc.yml
ci_cd:
  enabled: true
  min_score: 80
  fail_on_warnings: false
  
  # Quality gates
  quality_gates:
    - name: "Documentation"
      min_score: 90
      category: "documentation"
    
    - name: "Commit Quality"
      min_score: 85
      category: "commits"
    
    - name: "Git Hygiene"
      min_score: 80
      category: "hygiene"
```

## Best Practices

### For Teams
1. **Set Realistic Thresholds**: Start with achievable scores
2. **Gradual Improvement**: Increase thresholds over time
3. **Team Education**: Train team on health standards
4. **Regular Reviews**: Review health trends regularly
5. **Continuous Improvement**: Use results to improve processes

### For Organizations
1. **Standardized Thresholds**: Use consistent thresholds across projects
2. **Quality Metrics**: Track quality metrics over time
3. **Team Accountability**: Share results with teams
4. **Process Improvement**: Use data to improve processes
5. **Training Programs**: Implement quality training programs

## Troubleshooting

### Common Issues
- **Build Failures**: Check threshold settings
- **Performance Issues**: Optimize check frequency
- **False Positives**: Adjust threshold settings
- **Integration Issues**: Check CI/CD configuration

## Next Steps
- [Pre-commit Hooks](pre-commit-hooks.md) - Pre-commit integration guide
- [Historical Tracking](historical-tracking.md) - Health trend analysis
- [Multi-Repository Scan](multi-repository-scan.md) - Batch repository analysis
