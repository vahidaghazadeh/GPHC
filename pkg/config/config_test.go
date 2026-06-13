package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigUsesFileAndEnvironment(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "gphc.yml")
	if err := os.WriteFile(configPath, []byte("max_commits_to_analyze: 12\nmax_commit_size_lines: 90\n"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GPHC_MAX_COMMIT_SIZE_LINES", "120")

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.MaxCommitsToAnalyze != 12 {
		t.Fatalf("MaxCommitsToAnalyze = %d, want 12", cfg.MaxCommitsToAnalyze)
	}
	if cfg.MaxCommitSizeLines != 120 {
		t.Fatalf("environment override = %d, want 120", cfg.MaxCommitSizeLines)
	}
}
