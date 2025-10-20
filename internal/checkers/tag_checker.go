package checkers

import (
    "fmt"
    "os/exec"
    "regexp"
    "strings"
    "time"

    "github.com/opsource/gphc/pkg/types"
)

// TagChecker validates git tags, freshness, and policies
type TagChecker struct {
    BaseChecker
}

func NewTagChecker() *TagChecker {
    return &TagChecker{
        BaseChecker: NewBaseChecker("Tag & Release Checker", "TAGS", types.CategoryHygiene, 6),
    }
}

func (tc *TagChecker) Check(data *types.RepositoryData) *types.CheckResult {
    result := &types.CheckResult{
        ID:       tc.ID(),
        Name:     tc.Name(),
        Category: tc.Category(),
        Status:   types.StatusPass,
        Score:    0,
        Details:  []string{},
    }

    // 1) Collect tags
    tags, err := gitTagList(data.Path)
    if err != nil {
        result.Status = types.StatusFail
        result.Details = append(result.Details, fmt.Sprintf("Error reading tags: %v", err))
        return result
    }

    if len(tags) == 0 {
        result.Status = types.StatusFail
        result.Details = append(result.Details, "No tags found in repository")
        return result
    }

    // 2) Semantic validation
    semverRe := regexp.MustCompile(`^v\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$`)
    invalid := 0
    for _, t := range tags {
        if !semverRe.MatchString(t.Name) {
            invalid++
        }
    }
    if invalid == 0 {
        result.Score += 30
        result.Details = append(result.Details, "Semantic Versioning: All tags valid")
    } else {
        result.Status = types.StatusFail
        result.Details = append(result.Details, fmt.Sprintf("%d tag(s) with invalid semantic version", invalid))
    }

    // 3) Freshness: days since latest tag
    latest := tags[0]
    for _, t := range tags {
        if t.Date.After(latest.Date) {
            latest = t
        }
    }
    days := int(time.Since(latest.Date).Hours() / 24)
    result.Details = append(result.Details, fmt.Sprintf("Last tag: %s (%d days ago)", latest.Name, days))
    if days <= 45 {
        result.Score += 20
    } else {
        result.Status = types.StatusWarning
    }

    // 4) Unreleased commits since latest tag
    unreleased, err := gitUnreleasedCount(data.Path, latest.Name)
    if err == nil {
        if unreleased <= 3 {
            result.Score += 20
        } else {
            result.Status = types.StatusWarning
        }
        result.Details = append(result.Details, fmt.Sprintf("Unreleased commits since %s: %d", latest.Name, unreleased))
    } else {
        result.Details = append(result.Details, fmt.Sprintf("Cannot compute unreleased commits: %v", err))
    }

    // 5) Annotated vs lightweight ratio
    annotatedPct, err := gitAnnotatedPercentage(data.Path)
    if err == nil {
        if annotatedPct == 100 {
            result.Score += 20
            result.Details = append(result.Details, "Annotated Tags: 100%")
        } else {
            result.Details = append(result.Details, fmt.Sprintf("Annotated Tags: %d%%", annotatedPct))
        }
    }

    // Cap score to 100
    if result.Score > 100 {
        result.Score = 100
    }

    return result
}

type repoTag struct {
    Name string
    Date time.Time
}

func gitTagList(repoPath string) ([]repoTag, error) {
    cmd := exec.Command("git", "-C", repoPath, "for-each-ref", "--format=%(refname:short)|%(taggerdate:iso8601)", "refs/tags")
    out, err := cmd.Output()
    if err != nil {
        // Fallback to simple tag list (no dates)
        cmd = exec.Command("git", "-C", repoPath, "tag", "--list")
        out, err = cmd.Output()
        if err != nil {
            return nil, err
        }
        lines := strings.Split(strings.TrimSpace(string(out)), "\n")
        tags := make([]repoTag, 0, len(lines))
        for _, l := range lines {
            l = strings.TrimSpace(l)
            if l == "" {
                continue
            }
            tags = append(tags, repoTag{Name: l, Date: time.Time{}})
        }
        return tags, nil
    }

    lines := strings.Split(strings.TrimSpace(string(out)), "\n")
    tags := make([]repoTag, 0, len(lines))
    for _, l := range lines {
        if l == "" {
            continue
        }
        parts := strings.SplitN(l, "|", 2)
        name := parts[0]
        var dt time.Time
        if len(parts) == 2 {
            if t, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[1])); err == nil {
                dt = t
            }
        }
        tags = append(tags, repoTag{Name: name, Date: dt})
    }
    return tags, nil
}

func gitUnreleasedCount(repoPath, latestTag string) (int, error) {
    if latestTag == "" {
        return 0, nil
    }
    rangeExpr := fmt.Sprintf("%s..HEAD", latestTag)
    cmd := exec.Command("git", "-C", repoPath, "log", rangeExpr, "--oneline")
    out, err := cmd.Output()
    if err != nil {
        return 0, err
    }
    lines := strings.Split(strings.TrimSpace(string(out)), "\n")
    if len(lines) == 1 && lines[0] == "" {
        return 0, nil
    }
    return len(lines), nil
}

func gitAnnotatedPercentage(repoPath string) (int, error) {
    cmd := exec.Command("git", "-C", repoPath, "for-each-ref", "--format=%(refname:short) %(objecttype)", "refs/tags")
    out, err := cmd.Output()
    if err != nil {
        return 0, err
    }
    lines := strings.Split(strings.TrimSpace(string(out)), "\n")
    if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
        return 0, nil
    }
    total := 0
    annotated := 0
    for _, l := range lines {
        if l == "" {
            continue
        }
        total++
        if strings.Contains(l, "tag") { // annotated tags are objects of type 'tag'
            annotated++
        }
    }
    return int(float64(annotated) / float64(total) * 100.0), nil
}


