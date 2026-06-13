package checkers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vahidaghazadeh/gphc/pkg/types"
)

func TestTagCheckerUsesRepositoryPath(t *testing.T) {
	repo := createGitRepository(t)
	runGit(t, repo, "tag", "-a", "v1.0.0", "-m", "release")

	other := createGitRepository(t)
	runGit(t, other, "tag", "-a", "v9.0.0", "-m", "other release")

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(other); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })

	result := NewTagChecker().Check(&types.RepositoryData{Path: repo})
	details := strings.Join(result.Details, "\n")
	if !strings.Contains(details, "v1.0.0") || strings.Contains(details, "v9.0.0") {
		t.Fatalf("tag details came from the wrong repository: %s", details)
	}

	next, err := SuggestNextTag(repo)
	if err != nil {
		t.Fatal(err)
	}
	if next != "v1.0.1" {
		t.Fatalf("SuggestNextTag() = %q, want v1.0.1", next)
	}
}

func TestSuggestNextTagRecognizesScopedFeatures(t *testing.T) {
	repo := createGitRepository(t)
	runGit(t, repo, "tag", "-a", "v1.0.0", "-m", "release")
	if err := os.WriteFile(filepath.Join(repo, "feature.txt"), []byte("feature\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "feat(cli): add command")

	next, err := SuggestNextTag(repo)
	if err != nil {
		t.Fatal(err)
	}
	if next != "v1.1.0" {
		t.Fatalf("SuggestNextTag() = %q, want v1.1.0", next)
	}
}

func createGitRepository(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runGit(t, repo, "init", "-q")
	runGit(t, repo, "config", "user.email", "test@example.com")
	runGit(t, repo, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("# Test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-qm", "chore: initial commit")
	return repo
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, output)
	}
}
