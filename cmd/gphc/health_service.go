package main

import (
	"fmt"

	"github.com/vahidaghazadeh/gphc/internal/checkers"
	"github.com/vahidaghazadeh/gphc/internal/git"
	"github.com/vahidaghazadeh/gphc/internal/scorer"
	"github.com/vahidaghazadeh/gphc/pkg/types"
)

func buildHealthReport(repoPath string) (*types.HealthReport, error) {
	repositoryConfig, err := loadRepositoryConfig(repoPath)
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	analyzer, err := git.NewRepositoryAnalyzerWithOptions(
		repoPath,
		repositoryConfig.MaxCommitsToAnalyze,
		repositoryConfig.StaleBranchThresholdDays,
	)
	if err != nil {
		return nil, fmt.Errorf("initialize repository analyzer: %w", err)
	}

	data, err := analyzer.Analyze()
	if err != nil {
		return nil, fmt.Errorf("analyze repository: %w", err)
	}

	allCheckers := []checkers.Checker{
		checkers.NewDocChecker(),
		checkers.NewSetupChecker(),
		checkers.NewIgnoreChecker(),
		checkers.NewConventionalCommitChecker(),
		checkers.NewMsgLengthCheckerWithLimit(repositoryConfig.MaxCommitMessageLength),
		checkers.NewCommitSizeCheckerWithLimit(repositoryConfig.MaxCommitSizeLines),
		checkers.NewCommitAuthorInsightsChecker(),
		checkers.NewCodebaseSmellChecker(),
		checkers.NewLocalBranchChecker(),
		checkers.NewStaleBranchChecker(),
		checkers.NewBareRepoChecker(),
		checkers.NewStashChecker(),
		checkers.NewGitHubIntegrationChecker(),
		checkers.NewGitLabIntegrationChecker(),
		checkers.NewTagChecker(),
		checkers.NewSecretChecker(),
		checkers.NewTransitiveDependencyChecker(),
		checkers.NewGitPolicyChecker(),
		checkers.NewBinaryFileChecker(),
	}

	healthScorer := scorer.NewScorerWithWeights(map[types.Category]int{
		types.CategoryDocs:      repositoryConfig.Weights.Documentation,
		types.CategoryCommits:   repositoryConfig.Weights.Commits,
		types.CategoryHygiene:   repositoryConfig.Weights.Hygiene,
		types.CategoryStructure: repositoryConfig.Weights.Structure,
		types.CategorySecurity:  repositoryConfig.Weights.Security,
	})
	for _, checker := range allCheckers {
		healthScorer.AddResult(*checker.Check(data))
	}

	return healthScorer.CalculateHealthReport(), nil
}
