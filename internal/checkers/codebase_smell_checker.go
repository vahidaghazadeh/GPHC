package checkers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// CodebaseSmellChecker performs lightweight codebase structure analysis
type CodebaseSmellChecker struct {
	BaseChecker
}

// NewCodebaseSmellChecker creates a new codebase smell checker
func NewCodebaseSmellChecker() *CodebaseSmellChecker {
	return &CodebaseSmellChecker{
		BaseChecker: BaseChecker{
			id:       "CBS-801",
			name:     "Codebase Smell Check",
			category: types.CategoryStructure,
		},
	}
}

// Check performs lightweight codebase structure analysis
func (c *CodebaseSmellChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result := &types.CheckResult{
		ID:        c.ID(),
		Name:      c.Name(),
		Category:  c.Category(),
		Status:    types.StatusPass,
		Score:     0,
		Message:   "Codebase structure analysis completed",
		Details:   []string{},
		Timestamp: time.Now(),
	}

	// Analyze codebase structure
	smells := c.analyzeCodebaseStructure(data.Path)

	score := 100
	totalChecks := 0
	warnings := 0
	issues := 0

	// Check for missing test directories
	if smells.MissingTestDir {
		score -= 30
		warnings++
		result.Details = append(result.Details, "No test directory found")
		result.Details = append(result.Details, "Consider adding tests/ or test/ directory")
	} else {
		result.Details = append(result.Details, "Test directory found")
	}
	totalChecks++

	// Check for oversized directories
	if smells.OversizedDirs > 0 {
		score -= 20
		issues++
		result.Details = append(result.Details, fmt.Sprintf("%d oversized directory(ies) found (>1000 files)", smells.OversizedDirs))
		result.Details = append(result.Details, "Consider splitting large directories into smaller modules")
	} else {
		result.Details = append(result.Details, "No oversized directories")
	}
	totalChecks++

	// Check code to test ratio
	if smells.CodeToTestRatio < 0.1 {
		score -= 25
		warnings++
		result.Details = append(result.Details, fmt.Sprintf("Low test coverage ratio (%.1f%%)", smells.CodeToTestRatio*100))
		result.Details = append(result.Details, "Consider adding more test files")
	} else if smells.CodeToTestRatio < 0.3 {
		score -= 10
		result.Details = append(result.Details, fmt.Sprintf("Moderate test coverage ratio (%.1f%%)", smells.CodeToTestRatio*100))
		result.Details = append(result.Details, "Consider improving test coverage")
	} else {
		result.Details = append(result.Details, fmt.Sprintf("Good test coverage ratio (%.1f%%)", smells.CodeToTestRatio*100))
	}
	totalChecks++

	// Check for empty directories
	if smells.EmptyDirs > 0 {
		score -= 5
		result.Details = append(result.Details, fmt.Sprintf("%d empty directory(ies) found", smells.EmptyDirs))
		result.Details = append(result.Details, "Consider removing empty directories")
	} else {
		result.Details = append(result.Details, "No empty directories")
	}
	totalChecks++

	// Check for deep nesting
	if smells.MaxDepth > 6 {
		score -= 15
		warnings++
		result.Details = append(result.Details, fmt.Sprintf("Deep directory nesting detected (depth: %d)", smells.MaxDepth))
		result.Details = append(result.Details, "Consider flattening directory structure")
	} else {
		result.Details = append(result.Details, fmt.Sprintf("Reasonable directory depth (%d)", smells.MaxDepth))
	}
	totalChecks++

	// Check for too many files in root
	if smells.RootFiles > 20 {
		score -= 10
		result.Details = append(result.Details, fmt.Sprintf("Too many files in root directory (%d)", smells.RootFiles))
		result.Details = append(result.Details, "Consider organizing files into subdirectories")
	} else {
		result.Details = append(result.Details, fmt.Sprintf("Reasonable number of root files (%d)", smells.RootFiles))
	}
	totalChecks++

	// Check for missing common directories
	if !smells.HasSrcDir && !smells.HasLibDir && !smells.HasAppDir {
		score -= 10
		result.Details = append(result.Details, "No standard source directories (src/, lib/, app/)")
		result.Details = append(result.Details, "Consider organizing code into standard directories")
	} else {
		result.Details = append(result.Details, "Standard source directories found")
	}
	totalChecks++

	// Check for documentation files
	if smells.DocFiles == 0 {
		score -= 15
		warnings++
		result.Details = append(result.Details, "No documentation files found")
		result.Details = append(result.Details, "Consider adding documentation")
	} else {
		result.Details = append(result.Details, fmt.Sprintf("Found %d documentation file(s)", smells.DocFiles))
	}
	totalChecks++

	// Summary
	result.Details = append(result.Details, fmt.Sprintf("\nCodebase Statistics:"))
	result.Details = append(result.Details, fmt.Sprintf("  Total directories: %d", smells.TotalDirs))
	result.Details = append(result.Details, fmt.Sprintf("  Total files: %d", smells.TotalFiles))
	result.Details = append(result.Details, fmt.Sprintf("  Test files: %d", smells.TestFiles))
	result.Details = append(result.Details, fmt.Sprintf("  Documentation files: %d", smells.DocFiles))
	result.Details = append(result.Details, fmt.Sprintf("  Max directory depth: %d", smells.MaxDepth))

	// Determine final status
	if issues > 0 {
		result.Status = types.StatusFail
		result.Message = "Codebase structure issues detected"
	} else if warnings > 0 {
		result.Status = types.StatusWarning
		result.Message = "Codebase structure warnings detected"
	} else {
		result.Status = types.StatusPass
		result.Message = "Codebase structure looks good"
	}

	result.Score = score
	return result
}

// CodebaseSmells represents detected codebase structure issues
type CodebaseSmells struct {
	MissingTestDir  bool
	OversizedDirs   int
	CodeToTestRatio float64
	EmptyDirs       int
	MaxDepth        int
	RootFiles       int
	HasSrcDir       bool
	HasLibDir       bool
	HasAppDir       bool
	DocFiles        int
	TotalFiles      int
	TotalDirs       int
	TestFiles       int
}

// analyzeCodebaseStructure analyzes the codebase structure
func (c *CodebaseSmellChecker) analyzeCodebaseStructure(repoPath string) *CodebaseSmells {
	smells := &CodebaseSmells{}

	// Track directories and files
	var dirs []string
	var files []string
	var testFiles []string
	var docFiles []string

	maxDepth := 0
	rootFiles := 0

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		// Skip .git directory
		if strings.Contains(path, ".git") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Calculate depth
		relPath, _ := filepath.Rel(repoPath, path)
		depth := strings.Count(relPath, string(filepath.Separator))
		if depth > maxDepth {
			maxDepth = depth
		}

		// Count root files
		if depth == 0 && !info.IsDir() {
			rootFiles++
		}

		if info.IsDir() {
			dirs = append(dirs, path)

			// Check for empty directories
			if isEmptyDir(path) {
				smells.EmptyDirs++
			}

			// Check for standard directories
			dirName := filepath.Base(path)
			switch dirName {
			case "src":
				smells.HasSrcDir = true
			case "lib":
				smells.HasLibDir = true
			case "app":
				smells.HasAppDir = true
			case "test", "tests":
				smells.MissingTestDir = false
			}

			// Check for oversized directories
			if fileCount := countFilesInDir(path); fileCount > 1000 {
				smells.OversizedDirs++
			}
		} else {
			files = append(files, path)

			// Check file types
			fileName := strings.ToLower(filepath.Base(path))
			ext := strings.ToLower(filepath.Ext(path))

			// Test files
			if strings.Contains(fileName, "test") ||
				strings.Contains(fileName, "_test") ||
				ext == ".test" ||
				ext == ".spec" {
				testFiles = append(testFiles, path)
			}

			// Documentation files
			if ext == ".md" || ext == ".rst" || ext == ".txt" ||
				fileName == "readme" || fileName == "changelog" ||
				fileName == "license" || fileName == "copying" {
				docFiles = append(docFiles, path)
			}
		}

		return nil
	})

	if err != nil {
		// If we can't walk the directory, return basic info
		smells.MissingTestDir = true
		return smells
	}

	// Set missing test dir flag
	smells.MissingTestDir = len(testFiles) == 0

	// Calculate ratios
	smells.TotalFiles = len(files)
	smells.TotalDirs = len(dirs)
	smells.TestFiles = len(testFiles)
	smells.DocFiles = len(docFiles)
	smells.MaxDepth = maxDepth
	smells.RootFiles = rootFiles

	if len(files) > 0 {
		smells.CodeToTestRatio = float64(len(testFiles)) / float64(len(files))
	}

	return smells
}

// isEmptyDir checks if a directory is empty
func isEmptyDir(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) == 0
}

// countFilesInDir counts files in a directory (non-recursive)
func countFilesInDir(path string) int {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			count++
		}
	}
	return count
}
