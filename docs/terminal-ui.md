# Terminal UI (TUI) Guide

The GPHC terminal interface provides a focused interactive view of a
repository health report.

## Launching

```bash
# Analyze the current Git repository
git hc tui

# Analyze another Git repository
git hc tui /path/to/repository
```

The supplied path must be a Git repository.

## Views

The interface has two tabs:

- **Overview**: overall score, grade, result counts, repository, and timestamp
- **Details**: status, identifier, name, and message for every check

The report uses the same analyzer, configuration, and scoring path as
`git hc check`, so CLI, TUI, and dashboard results remain consistent.

## Keyboard Controls

| Key | Action |
| --- | --- |
| `Tab` | Select the next tab |
| `Shift+Tab` | Select the previous tab |
| `r` | Refresh repository health data |
| `q` or `Ctrl+C` | Quit |

The TUI uses the terminal alternate screen and returns to the original
terminal content after exit.

## Troubleshooting

If the interface reports that the path is not a Git repository, verify it
with:

```bash
git -C /path/to/repository rev-parse --git-dir
```

For non-interactive output or file export, use:

```bash
git hc check /path/to/repository --format json
git hc check /path/to/repository --format markdown --output report.md
```

## Related Guides

- [Web Dashboard](web-dashboard.md)
- [Multi-Repository Scan](multi-repository-scan.md)
- [Health Checks](health-checks.md)
