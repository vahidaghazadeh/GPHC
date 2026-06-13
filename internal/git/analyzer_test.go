package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAnalyzerIncludesLatestCommitStatsAndHonorsLimit(t *testing.T) {
	repo := t.TempDir()
	runGitCommand(t, repo, "init", "-q")
	runGitCommand(t, repo, "config", "user.email", "test@example.com")
	runGitCommand(t, repo, "config", "user.name", "Test")

	file := filepath.Join(repo, "data.txt")
	if err := os.WriteFile(file, []byte("one\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGitCommand(t, repo, "add", ".")
	runGitCommand(t, repo, "commit", "-qm", "chore: initial commit")

	if err := os.WriteFile(file, []byte("one\ntwo\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGitCommand(t, repo, "add", ".")
	runGitCommand(t, repo, "commit", "-qm", "feat: add second line")

	analyzer, err := NewRepositoryAnalyzerWithOptions(repo, 1, 30)
	if err != nil {
		t.Fatal(err)
	}
	data, err := analyzer.Analyze()
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Commits) != 1 {
		t.Fatalf("got %d commits, want 1", len(data.Commits))
	}
	if data.Commits[0].LinesAdded == 0 {
		t.Fatalf("latest commit stats were not calculated: %+v", data.Commits[0])
	}
}

func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, output)
	}
}
