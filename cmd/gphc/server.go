package main

import (
	"crypto/subtle"
	_ "embed"
	"encoding/json"
	"html"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/vahidaghazadeh/gphc/internal/checkers"
	"github.com/vahidaghazadeh/gphc/internal/exporter"
	"github.com/vahidaghazadeh/gphc/internal/git"
)

//go:embed dashboard.html
var dashboardTemplate string

// HTTP handlers
func isAuthorized(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}

	usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(serverUsername))
	passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(serverPassword))
	return usernameMatch&passwordMatch == 1
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	// Check authentication if enabled
	if serverAuth {
		if !isAuthorized(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="GPHC Dashboard"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
	}

	htmlContent := strings.ReplaceAll(dashboardTemplate, "{{TITLE}}", html.EscapeString(serverTitle))

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}

func handleHealthAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication if enabled
	if serverAuth {
		if !isAuthorized(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="GPHC Dashboard"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
	}

	repoPath := serverRepoPath

	healthReport, err := buildHealthReport(repoPath)
	if err != nil {
		http.Error(w, "Error analyzing repository", http.StatusInternalServerError)
		return
	}

	// Set CORS headers if enabled
	if serverCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"overall_score": healthReport.OverallScore,
		"grade":         healthReport.Grade,
		"summary":       healthReport.Summary,
		"timestamp":     healthReport.Timestamp.Format(time.RFC3339),
		"repository":    filepath.Base(repoPath),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func handleDiffAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication if enabled
	if serverAuth {
		if !isAuthorized(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="GPHC Dashboard"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	repoPath := serverRepoPath

	// Check if it's a git repository
	if !isGitRepository(repoPath) {
		http.Error(w, "Not a git repository", http.StatusBadRequest)
		return
	}

	// Get diff type and file from query parameters
	diffType := r.URL.Query().Get("type")
	if diffType == "" {
		diffType = "all"
	}
	file := r.URL.Query().Get("file")

	// If specific file requested, return single file diff
	if file != "" {
		handleSingleFileDiff(w, repoPath, diffType, file)
		return
	}

	// Get list of changed files first
	var cmd *exec.Cmd
	switch diffType {
	case "staged":
		cmd = exec.Command("git", "diff", "--cached", "--name-only")
	case "unstaged":
		cmd = exec.Command("git", "diff", "--name-only")
	case "all":
		cmd = exec.Command("git", "diff", "HEAD", "--name-only")
	default:
		cmd = exec.Command("git", "diff", "HEAD", "--name-only")
	}

	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Error running git diff", http.StatusInternalServerError)
		return
	}

	// Parse file list
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	var fileList []string
	for _, f := range files {
		if strings.TrimSpace(f) != "" {
			fileList = append(fileList, strings.TrimSpace(f))
		}
	}

	// If no files changed, return empty response
	if len(fileList) == 0 {
		response := map[string]interface{}{
			"status":    "success",
			"files":     []map[string]interface{}{},
			"timestamp": time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get diff for each file
	var fileDiffs []map[string]interface{}
	for _, fileName := range fileList {
		fileDiff := getFileDiff(repoPath, diffType, fileName)
		fileDiffs = append(fileDiffs, fileDiff)
	}

	response := map[string]interface{}{
		"status":    "success",
		"files":     fileDiffs,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleSingleFileDiff(w http.ResponseWriter, repoPath, diffType, file string) {
	fileDiff := getFileDiff(repoPath, diffType, file)

	response := map[string]interface{}{
		"status":    "success",
		"file":      fileDiff,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getFileDiff(repoPath, diffType, fileName string) map[string]interface{} {
	// Build git diff command for specific file
	var cmd *exec.Cmd
	switch diffType {
	case "staged":
		cmd = exec.Command("git", "diff", "--cached", fileName)
	case "unstaged":
		cmd = exec.Command("git", "diff", fileName)
	case "all":
		cmd = exec.Command("git", "diff", "HEAD", fileName)
	default:
		cmd = exec.Command("git", "diff", "HEAD", fileName)
	}

	cmd.Dir = repoPath
	output, err := cmd.Output()

	if err != nil {
		return map[string]interface{}{
			"name":  fileName,
			"error": err.Error(),
			"lines": []map[string]interface{}{},
		}
	}

	// Parse diff output
	lines := strings.Split(string(output), "\n")
	var diffLines []map[string]interface{}

	// Check if there's no output (no changes)
	if len(strings.TrimSpace(string(output))) == 0 {
		diffLines = []map[string]interface{}{}
	} else {
		for _, line := range lines {
			if line == "" {
				continue
			}

			lineType := "context"
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				lineType = "addition"
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				lineType = "deletion"
			} else if strings.HasPrefix(line, "@@") {
				lineType = "hunk"
			} else if strings.HasPrefix(line, "+++") {
				lineType = "file_header"
			} else if strings.HasPrefix(line, "index ") {
				lineType = "index"
			} else if strings.HasPrefix(line, "diff --git") {
				lineType = "file_name"
			}

			diffLines = append(diffLines, map[string]interface{}{
				"content": line,
				"type":    lineType,
			})
		}
	}

	// Calculate stats for this file
	additions := 0
	deletions := 0
	for _, line := range diffLines {
		if line["type"] == "addition" {
			additions++
		} else if line["type"] == "deletion" {
			deletions++
		}
	}

	return map[string]interface{}{
		"name":       fileName,
		"lines":      diffLines,
		"additions":  additions,
		"deletions":  deletions,
		"totalLines": len(diffLines),
	}
}

func handleTagsAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication if enabled
	if serverAuth {
		if !isAuthorized(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="GPHC Dashboard"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
	}

	repoPath := serverRepoPath

	// Analyze repository
	analyzer, err := git.NewRepositoryAnalyzer(repoPath)
	if err != nil {
		http.Error(w, "Error analyzing repository", http.StatusInternalServerError)
		return
	}

	data, err := analyzer.Analyze()
	if err != nil {
		http.Error(w, "Error analyzing repository data", http.StatusInternalServerError)
		return
	}

	// Run TagChecker
	tagChecker := checkers.NewTagChecker()
	result := tagChecker.Check(data)

	// Set CORS headers if enabled
	if serverCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":     result.Status.String(),
		"score":      result.Score,
		"message":    result.Message,
		"details":    result.Details,
		"timestamp":  result.Timestamp.Format(time.RFC3339),
		"repository": filepath.Base(repoPath),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func handleExportJSON(w http.ResponseWriter, r *http.Request) {
	// Check authentication if enabled
	if serverAuth {
		if !isAuthorized(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="GPHC Dashboard"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
	}

	repoPath := serverRepoPath

	healthReport, err := buildHealthReport(repoPath)
	if err != nil {
		http.Error(w, "Error analyzing repository", http.StatusInternalServerError)
		return
	}

	// Export to JSON
	exporter := exporter.NewExporter()
	jsonData, err := exporter.Export(healthReport, "json")
	if err != nil {
		http.Error(w, "Error exporting to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=health-report.json")
	w.Write([]byte(jsonData))
}
