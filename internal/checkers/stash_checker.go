package checkers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// StashChecker checks for Git stash entries
type StashChecker struct {
	BaseChecker
}

// NewStashChecker creates a new stash checker
func NewStashChecker() *StashChecker {
	return &StashChecker{
		BaseChecker: NewBaseChecker("Git stash management", "STASH-501", types.CategoryHygiene, 5),
	}
}

// Check analyzes Git stash entries
func (c *StashChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:       c.ID(),
		Name:     c.Name(),
		Category: c.Category(),
		Status:   types.StatusPass,
		Score:    100,
		Message:  "No stash entries found",
		Details:  []string{},
	}

	// Check if we're in a git repository
	if !isGitRepository(data.Path) {
		result.Status = types.StatusFail
		result.Score = -50
		result.Message = "Not a Git repository"
		return result
	}

	// Get stash list
	stashEntries, err := getStashEntries(data.Path)
	if err != nil {
		result.Status = types.StatusWarning
		result.Score = 50
		result.Message = fmt.Sprintf("Could not analyze stash: %v", err)
		result.Details = append(result.Details, fmt.Sprintf("Debug: Path=%s", data.Path))
		return result
	}

	if len(stashEntries) == 0 {
		return result
	}

	// Analyze stash entries
	oldStashes := 0
	totalStashes := len(stashEntries)

	for _, entry := range stashEntries {
		if entry.IsOld() {
			oldStashes++
		}
		result.Details = append(result.Details, entry.String())
	}

	// Determine status and score
	if oldStashes > 0 {
		result.Status = types.StatusWarning
		result.Score = 50
		result.Message = fmt.Sprintf("Found %d stash entries (%d old)", totalStashes, oldStashes)
	} else {
		result.Status = types.StatusPass
		result.Score = 80
		result.Message = fmt.Sprintf("Found %d recent stash entries", totalStashes)
	}

	return result
}

// StashEntry represents a Git stash entry
type StashEntry struct {
	Index     int
	Message   string
	Timestamp time.Time
	Branch    string
}

// IsOld checks if stash entry is older than 30 days
func (s *StashEntry) IsOld() bool {
	return time.Since(s.Timestamp) > 30*24*time.Hour
}

// String returns a formatted string representation
func (s *StashEntry) String() string {
	age := time.Since(s.Timestamp)
	ageStr := formatDuration(age)

	status := "✅ Recent"
	if s.IsOld() {
		status = "⚠️ Old"
	}

	return fmt.Sprintf("%s stash@{%d}: %s (%s ago) [%s]",
		status, s.Index, s.Message, ageStr, s.Branch)
}

// getStashEntries retrieves Git stash entries
func getStashEntries(repoPath string) ([]StashEntry, error) {
	cmd := exec.Command("git", "stash", "list", "--format=%gd|%gs|%ct")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return []StashEntry{}, nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	entries := make([]StashEntry, 0, len(lines))

	for i, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue
		}

		// Parse timestamp (Unix timestamp)
		timestampInt := int64(0)
		if _, err := fmt.Sscanf(parts[2], "%d", &timestampInt); err != nil {
			continue
		}
		timestamp := time.Unix(timestampInt, 0)

		entry := StashEntry{
			Index:     i,
			Message:   parts[1],
			Timestamp: timestamp,
			Branch:    "main", // Default branch
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// formatDuration formats duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.0fd", d.Hours()/24)
}

// isGitRepository checks if path is a Git repository
func isGitRepository(path string) bool {
	gitPath := fmt.Sprintf("%s/.git", path)
	_, err := os.Stat(gitPath)
	return err == nil
}
