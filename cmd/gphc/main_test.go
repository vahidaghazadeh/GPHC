package main

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// Simple test to ensure the main package compiles
	// This is a placeholder test - in a real project you'd have more comprehensive tests

	if version == "" {
		t.Error("Version should not be empty")
	}

	expectedVersion := "1.0.0"
	if version != expectedVersion {
		t.Errorf("Expected version %s, got %s", expectedVersion, version)
	}
}

func TestMainFunction(t *testing.T) {
	// Test that main function exists and can be called
	// This is a basic smoke test

	// In a real test, you might test command line arguments
	// or other functionality, but for now this ensures the package builds
}

func TestBuildCommitSuggestion(t *testing.T) {
	tests := []struct {
		name    string
		changes []stagedChange
		want    string
	}{
		{
			name:    "single source file",
			changes: []stagedChange{{status: "M", path: "internal/checkers/binary_file_checker.go"}},
			want:    "chore(checkers): improve binary file checker",
		},
		{
			name:    "documentation",
			changes: []stagedChange{{status: "M", path: "README.md"}},
			want:    "docs: update README",
		},
		{
			name:    "new command",
			changes: []stagedChange{{status: "A", path: "cmd/report/main.go"}},
			want:    "feat(report): add main",
		},
		{
			name:    "tests",
			changes: []stagedChange{{status: "M", path: "internal/checkers/checkers_test.go"}},
			want:    "test(checkers): update checkers test",
		},
		{
			name: "shared package",
			changes: []stagedChange{
				{status: "M", path: "internal/reporter/reporter.go"},
				{status: "M", path: "internal/reporter/reporter_test.go"},
			},
			want: "chore(reporter): update reporter implementation",
		},
		{
			name:    "deleted file",
			changes: []stagedChange{{status: "D", path: "internal/legacy/adapter.go"}},
			want:    "refactor(legacy): remove adapter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildCommitSuggestion(tt.changes); got != tt.want {
				t.Fatalf("buildCommitSuggestion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShellQuote(t *testing.T) {
	got := shellQuote("fix(cli): handle user's input")
	want := "'fix(cli): handle user'\"'\"'s input'"
	if got != want {
		t.Fatalf("shellQuote() = %q, want %q", got, want)
	}
}
