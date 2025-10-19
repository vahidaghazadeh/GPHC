package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// GitLabClient handles GitLab API interactions
type GitLabClient struct {
	client  *http.Client
	token   string
	baseURL string
}

// GitLabProject represents a GitLab project
type GitLabProject struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Path                 string `json:"path"`
	PathWithNamespace    string `json:"path_with_namespace"`
	Description          string `json:"description"`
	DefaultBranch        string `json:"default_branch"`
	Visibility           string `json:"visibility"`
	IssuesEnabled        bool   `json:"issues_enabled"`
	MergeRequestsEnabled bool   `json:"merge_requests_enabled"`
	WikiEnabled          bool   `json:"wiki_enabled"`
	SnippetsEnabled      bool   `json:"snippets_enabled"`
	Archived             bool   `json:"archived"`
	CreatedAt            string `json:"created_at"`
	LastActivityAt       string `json:"last_activity_at"`
	StarCount            int    `json:"star_count"`
	ForksCount           int    `json:"forks_count"`
	OpenIssuesCount      int    `json:"open_issues_count"`
}

// GitLabBranchProtection represents branch protection rules
type GitLabBranchProtection struct {
	Name                      string              `json:"name"`
	PushAccessLevels          []GitLabAccessLevel `json:"push_access_levels"`
	MergeAccessLevels         []GitLabAccessLevel `json:"merge_access_levels"`
	UnprotectAccessLevels     []GitLabAccessLevel `json:"unprotect_access_levels"`
	CodeOwnerApprovalRequired bool                `json:"code_owner_approval_required"`
}

type GitLabAccessLevel struct {
	AccessLevel            int    `json:"access_level"`
	AccessLevelDescription string `json:"access_level_description"`
	UserID                 int    `json:"user_id,omitempty"`
	GroupID                int    `json:"group_id,omitempty"`
}

// GitLabPipeline represents a GitLab CI/CD pipeline
type GitLabPipeline struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	Ref       string `json:"ref"`
	Sha       string `json:"sha"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	WebURL    string `json:"web_url"`
}

// GitLabContributor represents a project contributor
type GitLabContributor struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

// GitLabMergeRequest represents a merge request
type GitLabMergeRequest struct {
	ID           int         `json:"id"`
	Title        string      `json:"title"`
	State        string      `json:"state"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
	Author       GitLabUser  `json:"author"`
	Assignee     *GitLabUser `json:"assignee"`
	SourceBranch string      `json:"source_branch"`
	TargetBranch string      `json:"target_branch"`
	WebURL       string      `json:"web_url"`
}

type GitLabUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

// NewGitLabClient creates a new GitLab client
func NewGitLabClient() *GitLabClient {
	token := os.Getenv("GPHC_TOKEN")
	if token == "" {
		token = os.Getenv("GITLAB_TOKEN")
	}

	baseURL := os.Getenv("GITLAB_URL")
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}

	return &GitLabClient{
		client:  &http.Client{Timeout: 30 * time.Second},
		token:   token,
		baseURL: baseURL,
	}
}

// IsAuthenticated checks if the client has a valid token
func (c *GitLabClient) IsAuthenticated() bool {
	return c.token != ""
}

// GetProjectInfo fetches project information
func (c *GitLabClient) GetProjectInfo(projectPath string) (*GitLabProject, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("GitLab token not provided. Set GPHC_TOKEN or GITLAB_TOKEN environment variable")
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s", c.baseURL, projectPath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitLab API error: %d", resp.StatusCode)
	}

	var project GitLabProject
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

// GetBranchProtection fetches branch protection rules
func (c *GitLabClient) GetBranchProtection(projectPath, branch string) (*GitLabBranchProtection, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("GitLab token not provided. Set GPHC_TOKEN or GITLAB_TOKEN environment variable")
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/protected_branches/%s", c.baseURL, projectPath, branch)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // No protection rules
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitLab API error: %d", resp.StatusCode)
	}

	var protection GitLabBranchProtection
	if err := json.NewDecoder(resp.Body).Decode(&protection); err != nil {
		return nil, err
	}

	return &protection, nil
}

// GetPipelines fetches GitLab CI/CD pipelines
func (c *GitLabClient) GetPipelines(projectPath string) ([]GitLabPipeline, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("GitLab token not provided. Set GPHC_TOKEN or GITLAB_TOKEN environment variable")
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/pipelines", c.baseURL, projectPath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitLab API error: %d", resp.StatusCode)
	}

	var pipelines []GitLabPipeline
	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return nil, err
	}

	return pipelines, nil
}

// GetContributors fetches project contributors
func (c *GitLabClient) GetContributors(projectPath string) ([]GitLabContributor, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("GitLab token not provided. Set GPHC_TOKEN or GITLAB_TOKEN environment variable")
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/contributors", c.baseURL, projectPath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitLab API error: %d", resp.StatusCode)
	}

	var contributors []GitLabContributor
	if err := json.NewDecoder(resp.Body).Decode(&contributors); err != nil {
		return nil, err
	}

	return contributors, nil
}

// GetMergeRequests fetches merge requests
func (c *GitLabClient) GetMergeRequests(projectPath string) ([]GitLabMergeRequest, error) {
	if !c.IsAuthenticated() {
		return nil, fmt.Errorf("GitLab token not provided. Set GPHC_TOKEN or GITLAB_TOKEN environment variable")
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests?state=opened", c.baseURL, projectPath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitLab API error: %d", resp.StatusCode)
	}

	var mergeRequests []GitLabMergeRequest
	if err := json.NewDecoder(resp.Body).Decode(&mergeRequests); err != nil {
		return nil, err
	}

	return mergeRequests, nil
}

// ExtractGitLabProjectInfo extracts project path from a Git remote URL
func ExtractGitLabProjectInfo(remoteURL string) (string, error) {
	// Handle different URL formats
	// https://gitlab.com/owner/repo.git
	// git@gitlab.com:owner/repo.git
	// https://gitlab.com/owner/repo

	// Remove .git suffix
	remoteURL = strings.TrimSuffix(remoteURL, ".git")

	// Handle SSH format
	if strings.HasPrefix(remoteURL, "git@gitlab.com:") {
		projectPath := strings.TrimPrefix(remoteURL, "git@gitlab.com:")
		return projectPath, nil
	}

	// Handle HTTPS format
	if strings.Contains(remoteURL, "gitlab.com/") {
		parts := strings.Split(remoteURL, "gitlab.com/")
		if len(parts) >= 2 {
			return parts[1], nil
		}
	}

	// Handle custom GitLab instances
	if strings.Contains(remoteURL, "/") && !strings.Contains(remoteURL, "github.com") {
		// Extract the path part after the domain
		parts := strings.Split(remoteURL, "://")
		if len(parts) >= 2 {
			pathParts := strings.Split(parts[1], "/")
			if len(pathParts) >= 3 {
				return strings.Join(pathParts[1:], "/"), nil
			}
		}
	}

	return "", fmt.Errorf("unable to extract project path from URL: %s", remoteURL)
}
