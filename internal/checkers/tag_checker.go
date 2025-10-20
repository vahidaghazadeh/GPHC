package checkers

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// TagChecker validates tags, freshness, and unreleased commits
type TagChecker struct {
	BaseChecker

	// Configurable thresholds (can be wired to config later)
	maxDaysSinceLastTag  int
	maxUnreleasedCommits int
	requireAnnotatedTags bool
}

func NewTagChecker() *TagChecker {
	return &TagChecker{
		BaseChecker:          NewBaseChecker("Tag & Release Checker", "TAGS", types.CategoryCommits, 6),
		maxDaysSinceLastTag:  45,
		maxUnreleasedCommits: 3,
		requireAnnotatedTags: true,
	}
}

func (tc *TagChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        "TAGS-901",
		Name:      "Tag & Release Health",
		Category:  types.CategoryCommits,
		Timestamp: time.Now(),
	}

	var details []string
	score := 0
	// maxScore reserved for future granular scoring

	// 1) Collect tags
	tags, err := gitTags()
	if err != nil {
		result.Status = types.StatusWarning
		result.Score = 10
		result.Message = "Could not read git tags"
		result.Details = []string{"git tags not available: " + err.Error()}
		return result
	}

	if len(tags) == 0 {
		result.Status = types.StatusWarning
		result.Score = 20
		result.Message = "No tags found in repository"
		result.Details = []string{"Consider creating your first release tag (vX.Y.Z)"}
		return result
	}

	// 2) Semantic version validation
	semverOK, invalid := validateSemanticTags(tags)
	if semverOK {
		details = append(details, "Semantic Versioning: OK")
		score += 30
	} else {
		details = append(details, "Invalid tags (non-semver): "+strings.Join(invalid, ", "))
	}

	// 3) Get latest tag and date
	latestTag, latestDate, err := latestTagAndDate()
	if err == nil {
		days := int(time.Since(latestDate).Hours() / 24)
		details = append(details, fmt.Sprintf("Last tag: %s (%d days ago)", latestTag, days))
		if days <= tc.maxDaysSinceLastTag {
			score += 25
		} else {
			details = append(details, fmt.Sprintf("Last tag older than %d days", tc.maxDaysSinceLastTag))
		}
	} else {
		details = append(details, "Could not determine last tag date: "+err.Error())
	}

	// 4) Unreleased commits
	unreleased, err := unreleasedCommitCount()
	if err == nil {
		details = append(details, fmt.Sprintf("Unreleased commits since last tag: %d", unreleased))
		if unreleased <= tc.maxUnreleasedCommits {
			score += 25
		} else {
			details = append(details, fmt.Sprintf("Too many unreleased commits (>%d)", tc.maxUnreleasedCommits))
		}
	} else {
		details = append(details, "Could not count unreleased commits: "+err.Error())
	}

	// 5) Annotated vs Lightweight tags ratio
	annotatedOK, annotatedPct, err := annotatedTagStats()
	if err == nil {
		details = append(details, fmt.Sprintf("Annotated tags: %d%%", annotatedPct))
		if annotatedOK || !tc.requireAnnotatedTags {
			score += 20
		} else {
			details = append(details, "Some release tags are lightweight; annotate release tags")
		}
	} else {
		details = append(details, "Could not determine tag types: "+err.Error())
	}

	// Finalize result
	result.Score = score
	result.Details = details

	switch {
	case score >= 80:
		result.Status = types.StatusPass
		result.Message = "Tags and releases look healthy"
	case score >= 50:
		result.Status = types.StatusWarning
		result.Message = "Some tag/release improvements suggested"
	default:
		result.Status = types.StatusFail
		result.Message = "Tagging strategy needs attention"
	}

	return result
}

// gitTags returns tag names sorted by version/date (as per git order)
func gitTags() ([]string, error) {
	out, err := exec.Command("git", "tag").CombinedOutput()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var tags []string
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags, nil
}

var semverRe = regexp.MustCompile(`^v\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$`)

func validateSemanticTags(tags []string) (bool, []string) {
	invalid := make([]string, 0)
	for _, t := range tags {
		if !semverRe.MatchString(t) {
			invalid = append(invalid, t)
		}
	}
	return len(invalid) == 0, invalid
}

func latestTagAndDate() (string, time.Time, error) {
	// Use git for-each-ref to get tag and committerdate
	out, err := exec.Command("bash", "-lc", `git for-each-ref --sort=-creatordate --format='%(refname:short)|%(creatordate:iso8601)' refs/tags | head -n1`).CombinedOutput()
	if err != nil {
		return "", time.Time{}, err
	}
	s := strings.TrimSpace(string(out))
	if s == "" || !strings.Contains(s, "|") {
		return "", time.Time{}, fmt.Errorf("no tags")
	}
	parts := strings.SplitN(s, "|", 2)
	tag := parts[0]
	when, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[1]))
	if err != nil {
		// Try a more flexible parse
		when, err = time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(parts[1]))
		if err != nil {
			return tag, time.Time{}, err
		}
	}
	return tag, when, nil
}

func unreleasedCommitCount() (int, error) {
	// Determine latest tag
	tagOut, _ := exec.Command("bash", "-lc", "git describe --tags --abbrev=0").CombinedOutput()
	latest := strings.TrimSpace(string(tagOut))
	if latest == "" {
		// no tags → consider all commits unreleased? Return 0 to avoid noise
		return 0, nil
	}
	out, err := exec.Command("bash", "-lc", fmt.Sprintf("git rev-list %s..HEAD --count", latest)).CombinedOutput()
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(out))
	// simple atoi
	var n int
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			continue
		}
		n = n*10 + int(ch-'0')
	}
	return n, nil
}

func annotatedTagStats() (bool, int, error) {
	out, err := exec.Command("bash", "-lc", `git for-each-ref --format='%(refname:short) %(objecttype)' refs/tags`).CombinedOutput()
	if err != nil {
		return false, 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && strings.TrimSpace(lines[0]) == "") {
		return false, 0, nil
	}
	total := 0
	annotated := 0
	for _, l := range lines {
		f := strings.Fields(l)
		if len(f) < 2 {
			continue
		}
		total++
		if f[1] == "tag" { // annotated tags have objecttype "tag"
			annotated++
		}
	}
	pct := int(0)
	if total > 0 {
		pct = int((float64(annotated) / float64(total)) * 100)
	}
	return annotated == total, pct, nil
}

// SuggestNextTag suggests the next semantic version based on commit messages since last tag
func SuggestNextTag() (string, error) {
	// Get latest tag
	tagOut, _ := exec.Command("bash", "-lc", "git describe --tags --abbrev=0").CombinedOutput()
	latest := strings.TrimSpace(string(tagOut))
	if latest == "" {
		return "v0.1.0", nil
	}
	commitsOut, _ := exec.Command("bash", "-lc", fmt.Sprintf("git log %s..HEAD --pretty=format:%s", latest, `"%s"`)).CombinedOutput()
	messages := strings.Split(strings.TrimSpace(string(commitsOut)), "\n")

	bumpMajor := false
	bumpMinor := false
	bumpPatch := false
	for _, m := range messages {
		mm := strings.ToLower(m)
		if strings.Contains(mm, "breaking change") || strings.HasPrefix(mm, "feat!:") || strings.HasPrefix(mm, "fix!:") {
			bumpMajor = true
			break
		}
		if strings.HasPrefix(mm, "feat:") {
			bumpMinor = true
		} else if strings.HasPrefix(mm, "fix:") {
			bumpPatch = true
		}
	}

	// Parse latest semver
	base := strings.TrimPrefix(latest, "v")
	parts := strings.Split(base, ".")
	if len(parts) != 3 {
		return "v0.1.0", nil
	}
	toInt := func(s string) int {
		n := 0
		for _, ch := range s {
			if ch < '0' || ch > '9' {
				continue
			}
			n = n*10 + int(ch-'0')
		}
		return n
	}
	major, minor, patch := toInt(parts[0]), toInt(parts[1]), toInt(parts[2])

	switch {
	case bumpMajor:
		major++
		minor = 0
		patch = 0
	case bumpMinor:
		minor++
		patch = 0
	case bumpPatch:
		patch++
	default:
		patch++
	}
	return fmt.Sprintf("v%d.%d.%d", major, minor, patch), nil
}

// GenerateChangelog creates a simple CHANGELOG.md between last tag and HEAD using conventional commit groups
func GenerateChangelog(outputPath string) (string, error) {
	// Determine range
	tagOut, _ := exec.Command("bash", "-lc", "git describe --tags --abbrev=0").CombinedOutput()
	latest := strings.TrimSpace(string(tagOut))
	rangeSpec := ""
	if latest != "" {
		rangeSpec = latest + "..HEAD"
	} else {
		rangeSpec = "HEAD"
	}

	out, err := exec.Command("bash", "-lc", fmt.Sprintf("git log %s --pretty=format:%s", rangeSpec, `"%s"`)).CombinedOutput()
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	groups := map[string][]string{"Features": {}, "Fixes": {}, "Docs": {}, "Refactors": {}, "Others": {}}
	order := []string{"Features", "Fixes", "Docs", "Refactors", "Others"}
	for _, l := range lines {
		ll := strings.TrimSpace(l)
		if strings.HasPrefix(strings.ToLower(ll), "feat:") {
			groups["Features"] = append(groups["Features"], ll)
		} else if strings.HasPrefix(strings.ToLower(ll), "fix:") {
			groups["Fixes"] = append(groups["Fixes"], ll)
		} else if strings.HasPrefix(strings.ToLower(ll), "docs:") {
			groups["Docs"] = append(groups["Docs"], ll)
		} else if strings.HasPrefix(strings.ToLower(ll), "refactor:") {
			groups["Refactors"] = append(groups["Refactors"], ll)
		} else if ll != "" {
			groups["Others"] = append(groups["Others"], ll)
		}
	}

	var b strings.Builder
	b.WriteString("# Changelog\n\n")
	if latest != "" {
		b.WriteString(fmt.Sprintf("Changes since %s\n\n", latest))
	}
	for _, g := range order {
		if len(groups[g]) == 0 {
			continue
		}
		b.WriteString("## " + g + "\n")
		sort.Strings(groups[g])
		for _, item := range groups[g] {
			b.WriteString("- " + item + "\n")
		}
		b.WriteString("\n")
	}

	content := b.String()
	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
			return "", err
		}
	}
	return content, nil
}
