package git

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/opsource/gphc/pkg/types"
)

// RepositoryAnalyzer analyzes git repositories
type RepositoryAnalyzer struct {
	repo *git.Repository
	path string
}

// NewRepositoryAnalyzer creates a new repository analyzer
func NewRepositoryAnalyzer(path string) (*RepositoryAnalyzer, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	return &RepositoryAnalyzer{
		repo: repo,
		path: path,
	}, nil
}

// Analyze performs a comprehensive analysis of the repository
func (ra *RepositoryAnalyzer) Analyze() (*types.RepositoryData, error) {
	data := &types.RepositoryData{
		Path: ra.path,
	}

	// Analyze commits
	commits, err := ra.analyzeCommits()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze commits: %w", err)
	}
	data.Commits = commits

	// Analyze branches
	branches, err := ra.analyzeBranches()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze branches: %w", err)
	}
	data.Branches = branches

	// Analyze files
	files, err := ra.analyzeFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze files: %w", err)
	}
	data.Files = files

	// Check for essential files
	ra.checkEssentialFiles(data)

	// Check .gitignore
	ra.checkGitignore(data)

	return data, nil
}

// analyzeCommits analyzes the last 50 commits
func (ra *RepositoryAnalyzer) analyzeCommits() ([]types.CommitInfo, error) {
	ref, err := ra.repo.Head()
	if err != nil {
		// Repository might be empty (no commits)
		return []types.CommitInfo{}, nil
	}

	cIter, err := ra.repo.Log(&git.LogOptions{
		From: ref.Hash(),
	})
	if err != nil {
		return nil, err
	}

	var commits []types.CommitInfo
	count := 0
	maxCommits := 50

	err = cIter.ForEach(func(c *object.Commit) error {
		if count >= maxCommits {
			return nil
		}

		// Calculate diff stats
		var linesAdded, linesDeleted int
		if count > 0 && len(c.ParentHashes) > 0 {
			parent := c.ParentHashes[0]
			pCommit, err := ra.repo.CommitObject(parent)
			if err == nil {
				stats, err := getCommitStats(c, pCommit)
				if err == nil {
					linesAdded = stats.Additions
					linesDeleted = stats.Deletions
				}
			}
		}

		message := c.Message
		subject := strings.Split(message, "\n")[0]
		body := strings.Join(strings.Split(message, "\n")[1:], "\n")

		commits = append(commits, types.CommitInfo{
			Hash:         c.Hash.String()[:8],
			Message:      message,
			Subject:      subject,
			Body:         body,
			Author:       c.Author.Name,
			Date:         c.Author.When,
			LinesAdded:   linesAdded,
			LinesDeleted: linesDeleted,
		})

		count++
		return nil
	})

	return commits, err
}

// analyzeBranches analyzes local branches
func (ra *RepositoryAnalyzer) analyzeBranches() ([]types.BranchInfo, error) {
	branches, err := ra.repo.Branches()
	if err != nil {
		return nil, err
	}

	var branchInfos []types.BranchInfo
	mainBranch, err := ra.repo.Reference(plumbing.HEAD, false)
	if err != nil {
		// Repository might be empty (no commits)
		return branchInfos, nil
	}

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			branchName := ref.Name().Short()
			
			// Get branch commit
			commit, err := ra.repo.CommitObject(ref.Hash())
			if err != nil {
				return err
			}

			// Check if branch is merged into main
			isMerged := ra.isBranchMerged(ref.Hash(), mainBranch.Hash())

			// Count commits in branch
			commitCount := ra.countBranchCommits(ref.Hash())

			// Check if branch is stale (older than 60 days)
			isStale := time.Since(commit.Author.When) > 60*24*time.Hour

			branchInfos = append(branchInfos, types.BranchInfo{
				Name:        branchName,
				IsMerged:    isMerged,
				LastCommit:  commit.Author.When,
				CommitCount: commitCount,
				IsStale:     isStale,
			})
		}
		return nil
	})

	return branchInfos, err
}

// analyzeFiles gets a list of files in the repository
func (ra *RepositoryAnalyzer) analyzeFiles() ([]string, error) {
	var files []string
	
	err := filepath.Walk(ra.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		
		if !info.IsDir() {
			relPath, err := filepath.Rel(ra.path, path)
			if err == nil {
				files = append(files, relPath)
			}
		}
		
		return nil
	})
	
	return files, err
}

// checkEssentialFiles checks for essential project files
func (ra *RepositoryAnalyzer) checkEssentialFiles(data *types.RepositoryData) {
	essentialFiles := map[string]*bool{
		"README.md":           &data.HasReadme,
		"LICENSE":             &data.HasLicense,
		"CONTRIBUTING.md":     &data.HasContributing,
		"CODE_OF_CONDUCT.md": &data.HasCodeOfConduct,
	}

	for filename, hasFile := range essentialFiles {
		filePath := filepath.Join(ra.path, filename)
		if _, err := os.Stat(filePath); err == nil {
			*hasFile = true
		}
	}
}

// checkGitignore checks for .gitignore file
func (ra *RepositoryAnalyzer) checkGitignore(data *types.RepositoryData) {
	gitignorePath := filepath.Join(ra.path, ".gitignore")
	if content, err := os.ReadFile(gitignorePath); err == nil {
		data.HasGitignore = true
		data.GitignoreContent = string(content)
	}
}

// isBranchMerged checks if a branch is merged into main
func (ra *RepositoryAnalyzer) isBranchMerged(branchHash, mainHash plumbing.Hash) bool {
	// Simple implementation - check if branch hash is reachable from main
	mainCommit, err := ra.repo.CommitObject(mainHash)
	if err != nil {
		return false
	}

	// Check if branch commit is in main's history
	iter := object.NewCommitIterCTime(mainCommit, nil, nil)
	defer iter.Close()

	for {
		commit, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false
		}

		if commit.Hash == branchHash {
			return true
		}
	}

	return false
}

// countBranchCommits counts commits in a branch
func (ra *RepositoryAnalyzer) countBranchCommits(branchHash plumbing.Hash) int {
	commit, err := ra.repo.CommitObject(branchHash)
	if err != nil {
		return 0
	}

	count := 0
	iter := object.NewCommitIterCTime(commit, nil, nil)
	defer iter.Close()

	for {
		_, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		count++
	}

	return count
}

// CommitStats represents commit statistics
type CommitStats struct {
	Additions int
	Deletions int
}

// getCommitStats calculates diff statistics between two commits
func getCommitStats(commit, parent *object.Commit) (*CommitStats, error) {
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	parentTree, err := parent.Tree()
	if err != nil {
		return nil, err
	}

	changes, err := object.DiffTree(parentTree, tree)
	if err != nil {
		return nil, err
	}

	stats := &CommitStats{}
	for _, change := range changes {
		patch, err := change.Patch()
		if err != nil {
			continue
		}

		fileStats := patch.Stats()
		for _, fileStat := range fileStats {
			stats.Additions += fileStat.Addition
			stats.Deletions += fileStat.Deletion
		}
	}

	return stats, nil
}
