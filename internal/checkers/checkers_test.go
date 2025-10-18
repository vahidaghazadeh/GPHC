package checkers

import (
	"testing"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

func TestDocChecker(t *testing.T) {
	checker := NewDocChecker()
	
	// Test basic properties
	if checker.Name() != "Documentation Checker" {
		t.Errorf("Expected name 'Documentation Checker', got '%s'", checker.Name())
	}
	
	if checker.ID() != "DOC" {
		t.Errorf("Expected ID 'DOC', got '%s'", checker.ID())
	}
	
	if checker.Category() != types.CategoryDocs {
		t.Errorf("Expected category CategoryDocs, got %v", checker.Category())
	}
	
	if checker.Weight() != 8 {
		t.Errorf("Expected weight 8, got %d", checker.Weight())
	}
}

func TestDocCheckerWithData(t *testing.T) {
	checker := NewDocChecker()
	
	// Test with repository data that has all files
	data := &types.RepositoryData{
		HasReadme:       true,
		HasLicense:      true,
		HasContributing: true,
		HasCodeOfConduct: true,
	}
	
	result := checker.Check(data)
	
	if result == nil {
		t.Error("Expected result, got nil")
		return
	}
	
	if result.ID != "DOC-101" {
		t.Errorf("Expected ID 'DOC-101', got '%s'", result.ID)
	}
	
	if result.Category != types.CategoryDocs {
		t.Errorf("Expected category CategoryDocs, got %v", result.Category)
	}
}

func TestIgnoreChecker(t *testing.T) {
	checker := NewIgnoreChecker()
	
	if checker.Name() != "Gitignore Checker" {
		t.Errorf("Expected name 'Gitignore Checker', got '%s'", checker.Name())
	}
	
	if checker.ID() != "IGNORE" {
		t.Errorf("Expected ID 'IGNORE', got '%s'", checker.ID())
	}
}

func TestConventionalCommitChecker(t *testing.T) {
	checker := NewConventionalCommitChecker()
	
	if checker.Name() != "Conventional Commit Checker" {
		t.Errorf("Expected name 'Conventional Commit Checker', got '%s'", checker.Name())
	}
	
	if checker.ID() != "CONV" {
		t.Errorf("Expected ID 'CONV', got '%s'", checker.ID())
	}
}

func TestConventionalCommitCheckerWithCommits(t *testing.T) {
	checker := NewConventionalCommitChecker()
	
	// Test with conventional commits
	data := &types.RepositoryData{
		Commits: []types.CommitInfo{
			{
				Hash:   "abc123",
				Subject: "feat: add new feature",
				Date:   time.Now(),
			},
			{
				Hash:   "def456",
				Subject: "fix: resolve bug",
				Date:   time.Now(),
			},
		},
	}
	
	result := checker.Check(data)
	
	if result == nil {
		t.Error("Expected result, got nil")
		return
	}
	
	if result.ID != "CHQ-301" {
		t.Errorf("Expected ID 'CHQ-301', got '%s'", result.ID)
	}
	
	// Should pass since both commits follow conventional format
	if result.Status != types.StatusPass {
		t.Errorf("Expected StatusPass, got %v", result.Status)
	}
}

func TestMsgLengthChecker(t *testing.T) {
	checker := NewMsgLengthChecker()
	
	if checker.Name() != "Message Length Checker" {
		t.Errorf("Expected name 'Message Length Checker', got '%s'", checker.Name())
	}
	
	if checker.ID() != "LENGTH" {
		t.Errorf("Expected ID 'LENGTH', got '%s'", checker.ID())
	}
}

func TestCommitSizeChecker(t *testing.T) {
	checker := NewCommitSizeChecker()
	
	if checker.Name() != "Commit Size Checker" {
		t.Errorf("Expected name 'Commit Size Checker', got '%s'", checker.Name())
	}
	
	if checker.ID() != "SIZE" {
		t.Errorf("Expected ID 'SIZE', got '%s'", checker.ID())
	}
}

func TestLocalBranchChecker(t *testing.T) {
	checker := NewLocalBranchChecker()
	
	if checker.Name() != "Local Branch Checker" {
		t.Errorf("Expected name 'Local Branch Checker', got '%s'", checker.Name())
	}
	
	if checker.ID() != "LOCAL" {
		t.Errorf("Expected ID 'LOCAL', got '%s'", checker.ID())
	}
}

func TestStaleBranchChecker(t *testing.T) {
	checker := NewStaleBranchChecker()
	
	if checker.Name() != "Stale Branch Checker" {
		t.Errorf("Expected name 'Stale Branch Checker', got '%s'", checker.Name())
	}
	
	if checker.ID() != "STALE" {
		t.Errorf("Expected ID 'STALE', got '%s'", checker.ID())
	}
}

func TestBareRepoChecker(t *testing.T) {
	checker := NewBareRepoChecker()
	
	if checker.Name() != "Bare Repository Checker" {
		t.Errorf("Expected name 'Bare Repository Checker', got '%s'", checker.Name())
	}
	
	if checker.ID() != "BARE" {
		t.Errorf("Expected ID 'BARE', got '%s'", checker.ID())
	}
}
