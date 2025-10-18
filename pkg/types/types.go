package types

import (
	"time"
)

// CheckResult represents the result of a single check
type CheckResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      Status    `json:"status"`
	Score       int       `json:"score"`
	Message     string    `json:"message"`
	Details     []string  `json:"details,omitempty"`
	Category    Category  `json:"category"`
	Timestamp   time.Time `json:"timestamp"`
}

// Status represents the status of a check
type Status int

const (
	StatusPass Status = iota
	StatusFail
	StatusWarning
)

func (s Status) String() string {
	switch s {
	case StatusPass:
		return "PASS"
	case StatusFail:
		return "FAIL"
	case StatusWarning:
		return "WARNING"
	default:
		return "UNKNOWN"
	}
}

// Category represents the category of a check
type Category int

const (
	CategoryDocs Category = iota
	CategoryCommits
	CategoryHygiene
)

func (c Category) String() string {
	switch c {
	case CategoryDocs:
		return "Documentation & Project Structure"
	case CategoryCommits:
		return "Commit History Quality"
	case CategoryHygiene:
		return "Git Cleanup & Hygiene"
	default:
		return "Unknown"
	}
}

// HealthReport represents the overall health report
type HealthReport struct {
	OverallScore int           `json:"overall_score"`
	Grade        string        `json:"grade"`
	Results      []CheckResult `json:"results"`
	Summary      ReportSummary `json:"summary"`
	Timestamp    time.Time    `json:"timestamp"`
}

// ReportSummary provides a summary of the health check
type ReportSummary struct {
	TotalChecks   int `json:"total_checks"`
	PassedChecks  int `json:"passed_checks"`
	FailedChecks  int `json:"failed_checks"`
	WarningChecks int `json:"warning_checks"`
}

// RepositoryData contains the analyzed repository data
type RepositoryData struct {
	Path           string
	Commits        []CommitInfo
	Branches        []BranchInfo
	Files           []string
	HasReadme       bool
	HasLicense      bool
	HasContributing bool
	HasCodeOfConduct bool
	HasGitignore    bool
	GitignoreContent string
}

// CommitInfo contains information about a commit
type CommitInfo struct {
	Hash        string
	Message     string
	Subject     string
	Body        string
	Author      string
	Date        time.Time
	LinesAdded  int
	LinesDeleted int
}

// BranchInfo contains information about a branch
type BranchInfo struct {
	Name         string
	IsMerged     bool
	LastCommit   time.Time
	CommitCount  int
	IsStale      bool
}
