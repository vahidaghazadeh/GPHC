package main

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiffAPIUsesConfiguredRepository(t *testing.T) {
	repo := createServerTestRepository(t)
	if err := os.WriteFile(filepath.Join(repo, "tracked.txt"), []byte("changed\n"), 0644); err != nil {
		t.Fatal(err)
	}
	serverRepoPath = repo
	serverAuth = false

	request := httptest.NewRequest("GET", "/api/diff?type=unstaged", nil)
	response := httptest.NewRecorder()
	handleDiffAPI(response, request)

	if response.Code != 200 {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	var payload struct {
		Files []map[string]interface{} `json:"files"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Files) != 1 || payload.Files[0]["name"] != "tracked.txt" {
		t.Fatalf("unexpected diff response: %#v", payload.Files)
	}
}

func TestDashboardEscapesTitle(t *testing.T) {
	serverAuth = false
	serverTitle = `<script>alert("x")</script>`
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	handleDashboard(response, request)

	body := response.Body.String()
	if strings.Contains(body, serverTitle) || !strings.Contains(body, "&lt;script&gt;") {
		t.Fatalf("dashboard title was not escaped")
	}
}

func createServerTestRepository(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runServerGit(t, repo, "init", "-q")
	runServerGit(t, repo, "config", "user.email", "test@example.com")
	runServerGit(t, repo, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(repo, "tracked.txt"), []byte("initial\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runServerGit(t, repo, "add", ".")
	runServerGit(t, repo, "commit", "-qm", "chore: initial commit")
	return repo
}

func runServerGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, output)
	}
}
