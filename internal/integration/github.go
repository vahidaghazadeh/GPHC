package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// GitHubClient handles GitHub API interactions
type GitHubClient struct {
	client  *http.Client
	token   string
	baseURL string
}

// GitHubBranchProtection represents branch protection rules
type GitHubBranchProtection struct {
	RequiredStatusChecks     *RequiredStatusChecks `json:"required_status_checks"`
	EnforceAdmins            bool                  `json:"enforce_admins"`
	RequiredPullRequestReviews *RequiredPullRequestReviews `json:"required_pull_request_reviews"`
	Restrictions             *Restrictions          `json:"restrictions"`
	AllowForcePushes         bool                  `json:"allow_force_pushes"`
	AllowDeletions           bool                  `json:"allow_deletions"`
}

type RequiredStatusChecks struct {
	Strict   bool     `json:"strict"`
	Contexts []string `json:"contexts"`
}

type RequiredPullRequestReviews struct {
	RequiredApprovingReviewCount int  `json:"required_approving_review_count"`
	DismissStaleReviews         bool `json:"dismiss_stale_reviews"`
	RequireCodeOwnerReviews     bool `json:"require_code_owner_reviews"`
}

type Restrictions struct {
	Users []string `json:"users"`
	Teams []string `json:"teams"`
}

// GitHubWorkflow represents a GitHub Actions workflow
type GitHubWorkflow struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GitHubContributor represents a repository contributor
type GitHubContributor struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	Contributions     int    `json:"contributions"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

// GitHubRepoInfo represents repository information
type GitHubRepoInfo struct {
	FullName        string `json:"full_name"`
	DefaultBranch   string `json:"default_branch"`
	IsPrivate       bool   `json:"private"`
	HasIssues       bool   `json:"has_issues"`
	HasProjects     bool   `json:"has_projects"`
	HasWiki         bool   `json:"has_wiki"`
	HasPages        bool   `json:"has_pages"`
	Archived        bool   `json:"archived"`
	Disabled        bool   `json:"disabled"`
	OpenIssuesCount int    `json:"open_issues_count"`
	ForksCount      int    `json:"forks_count"`
	StargazersCount int    `json:"stargazers_count"`
	WatchersCount   int    `json:"watchers_count"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	PushedAt        string `json:"pushed_at"`
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient() *GitHubClient {
	token := os.Getenv("GPHC_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	return &GitHubClient{
		client:  &http.Client{Timeout: 30 * time.Second},
		token:   token,
		baseURL: "https://api.github.com",
	}
}

// IsAuthenticated checks if the client has a valid token
func (c *GitHubClient) IsAuthenticated() bool {
	return c.token != ""
}

// GetRepositoryInfo fetches repository information
func (c *GitHubClient) GetRepositoryInfo(owner, repo string) (*GitHubRepoInfo, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("GitHub token not provided. Set GPHC_TOKEN or GITHUB_TOKEN environment variable")
	}

	url := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var repoInfo GitHubRepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&repoInfo); err != nil {
		return nil, err
	}

	return &repoInfo, nil
}

// GetBranchProtection fetches branch protection rules
func (c *GitHubClient) GetBranchProtection(owner, repo, branch string) (*GitHubBranchProtection, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("GitHub token not provided. Set GPHC_TOKEN or GITHUB_TOKEN environment variable")
	}

	url := fmt.Sprintf("%s/repos/%s/%s/branches/%s/protection", c.baseURL, owner, repo, branch)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // No protection rules
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var protection GitHubBranchProtection
	if err := json.NewDecoder(resp.Body).Decode(&protection); err != nil {
		return nil, err
	}

	return &protection, nil
}

// GetWorkflows fetches GitHub Actions workflows
func (c *GitHubClient) GetWorkflows(owner, repo string) ([]GitHubWorkflow, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("GitHub token not provided. Set GPHC_TOKEN or GITHUB_TOKEN environment variable")
	}

	url := fmt.Sprintf("%s/repos/%s/%s/actions/workflows", c.baseURL, owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var response struct {
		TotalCount int               `json:"total_count"`
		Workflows  []GitHubWorkflow  `json:"workflows"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Workflows, nil
}

// GetContributors fetches repository contributors
func (c *GitHubClient) GetContributors(owner, repo string) ([]GitHubContributor, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("GitHub token not provided. Set GPHC_TOKEN or GITHUB_TOKEN environment variable")
	}

	url := fmt.Sprintf("%s/repos/%s/%s/contributors", c.baseURL, owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var contributors []GitHubContributor
	if err := json.NewDecoder(resp.Body).Decode(&contributors); err != nil {
		return nil, err
	}

	return contributors, nil
}

// ExtractRepoInfo extracts owner and repo from a Git remote URL
func ExtractRepoInfo(remoteURL string) (owner, repo string, err error) {
	// Handle different URL formats
	// https://github.com/owner/repo.git
	// git@github.com:owner/repo.git
	// https://github.com/owner/repo

	// Remove .git suffix
	remoteURL = strings.TrimSuffix(remoteURL, ".git")

	// Handle SSH format
	if strings.HasPrefix(remoteURL, "git@github.com:") {
		parts := strings.Split(strings.TrimPrefix(remoteURL, "git@github.com:"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1], nil
		}
	}

	// Handle HTTPS format
	if strings.Contains(remoteURL, "github.com/") {
		parts := strings.Split(remoteURL, "github.com/")
		if len(parts) >= 2 {
			pathParts := strings.Split(parts[1], "/")
			if len(pathParts) >= 2 {
				return pathParts[0], pathParts[1], nil
			}
		}
	}

	return "", "", fmt.Errorf("unable to extract owner/repo from URL: %s", remoteURL)
}
