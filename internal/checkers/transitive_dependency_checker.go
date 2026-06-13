package checkers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/vahidaghazadeh/gphc/pkg/types"
)

// TransitiveDependencyChecker checks for vulnerabilities in transitive dependencies
type TransitiveDependencyChecker struct {
	BaseChecker
}

// Dependency represents a single dependency
type Dependency struct {
	Name            string          `json:"name"`
	Version         string          `json:"version"`
	Direct          bool            `json:"direct"`
	Vulnerable      bool            `json:"vulnerable"`
	Severity        string          `json:"severity"`
	Description     string          `json:"description"`
	Path            []string        `json:"path"`
	Children        []*Dependency   `json:"children"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	ID          string    `json:"id"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	CVSS        float64   `json:"cvss"`
	Published   time.Time `json:"published"`
	Fixed       string    `json:"fixed"`
}

// DependencyTree represents the complete dependency tree
type DependencyTree struct {
	Root       *Dependency `json:"root"`
	Total      int         `json:"total"`
	Vulnerable int         `json:"vulnerable"`
	Critical   int         `json:"critical"`
	High       int         `json:"high"`
	Medium     int         `json:"medium"`
	Low        int         `json:"low"`
}

// NewTransitiveDependencyChecker creates a new TransitiveDependencyChecker
func NewTransitiveDependencyChecker() *TransitiveDependencyChecker {
	return &TransitiveDependencyChecker{
		BaseChecker: NewBaseChecker("Transitive Dependency Vetting", "TRANSITIVE-DEPS", types.CategorySecurity, 9),
	}
}

// Check performs transitive dependency vulnerability scanning
func (c *TransitiveDependencyChecker) Check(data *types.RepositoryData) *types.CheckResult {
	result, _ := c.CheckWithOptions(data, false, "deep")
	return result
}

// CheckWithOptions inventories dependencies using the requested depth.
func (c *TransitiveDependencyChecker) CheckWithOptions(data *types.RepositoryData, directOnly bool, depth string) (*types.CheckResult, *DependencyTree) {
	result := &types.CheckResult{
		ID:        c.ID(),
		Name:      c.Name(),
		Status:    types.StatusPass,
		Score:     100,
		Message:   "No transitive dependency vulnerabilities found",
		Details:   []string{},
		Category:  c.Category(),
		Timestamp: time.Now(),
	}

	// Detect project type and analyze dependencies
	projectType := c.detectProjectType(data.Path)
	if projectType == "" {
		result.Status = types.StatusPass
		result.Message = "No supported dependency manifest found"
		result.Score = 100
		return result, nil
	}

	// Build dependency tree
	tree, err := c.buildDependencyTree(data.Path, projectType)
	if err != nil {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Failed to build dependency tree: %v", err)
		result.Score = 0
		return result, tree
	}

	if directOnly || depth == "shallow" {
		for _, dependency := range tree.Root.Children {
			dependency.Children = nil
		}
	}
	tree.Total = countDependencies(tree.Root)

	result.Status = types.StatusWarning
	result.Score = 100
	result.Message = "Dependency inventory built; vulnerability database scan is unavailable"

	// Add detailed information
	result.Details = append(result.Details, fmt.Sprintf("Project Type: %s", projectType))
	result.Details = append(result.Details, fmt.Sprintf("Total Dependencies: %d", tree.Total))
	result.Details = append(result.Details, "Run an ecosystem scanner for authoritative results:")
	result.Details = append(result.Details, vulnerabilityScannerHint(projectType))

	return result, tree
}

func countDependencies(root *Dependency) int {
	total := 0
	for _, child := range root.Children {
		total += 1 + countDependencies(child)
	}
	return total
}

func vulnerabilityScannerHint(projectType string) string {
	switch projectType {
	case "go":
		return "govulncheck ./..."
	case "nodejs":
		return "npm audit"
	case "python":
		return "pip-audit"
	case "rust":
		return "cargo audit"
	case "java":
		return "Use OWASP Dependency-Check or a Maven security plugin"
	default:
		return "Use an ecosystem-specific vulnerability scanner"
	}
}

// ScanVulnerabilities runs the installed authoritative scanner for the detected ecosystem.
func (c *TransitiveDependencyChecker) ScanVulnerabilities(data *types.RepositoryData, minSeverity string) *types.CheckResult {
	result := &types.CheckResult{
		ID:        c.ID(),
		Name:      c.Name(),
		Status:    types.StatusWarning,
		Score:     100,
		Category:  c.Category(),
		Timestamp: time.Now(),
	}

	projectType := c.detectProjectType(data.Path)
	command, args := vulnerabilityCommand(projectType, minSeverity)
	if command == "" {
		result.Message = "No supported vulnerability scanner is configured for this project"
		return result
	}
	if _, err := exec.LookPath(command); err != nil {
		result.Message = fmt.Sprintf("%s is not installed", command)
		result.Details = []string{vulnerabilityScannerHint(projectType)}
		return result
	}

	cmd := exec.Command(command, args...)
	cmd.Dir = data.Path
	output, err := cmd.CombinedOutput()
	result.Details = compactScannerOutput(string(output), 40)
	if err == nil {
		result.Status = types.StatusPass
		result.Score = 100
		result.Message = "No known vulnerabilities reported by " + command
		return result
	}

	if _, ok := err.(*exec.ExitError); ok {
		result.Status = types.StatusFail
		result.Score = 0
		result.Message = command + " reported dependency vulnerabilities or audit failures"
		return result
	}

	result.Message = fmt.Sprintf("Could not run %s: %v", command, err)
	return result
}

func vulnerabilityCommand(projectType, minSeverity string) (string, []string) {
	switch projectType {
	case "go":
		return "govulncheck", []string{"./..."}
	case "nodejs":
		return "npm", []string{"audit", "--audit-level=" + strings.ToLower(minSeverity)}
	case "python":
		return "pip-audit", nil
	case "rust":
		return "cargo", []string{"audit"}
	default:
		return "", nil
	}
}

func compactScannerOutput(output string, maxLines int) []string {
	var details []string
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		details = append(details, line)
		if len(details) == maxLines {
			details = append(details, "Scanner output truncated")
			break
		}
	}
	return details
}

// detectProjectType detects the type of project based on manifest files
func (c *TransitiveDependencyChecker) detectProjectType(repoPath string) string {
	manifestFiles := []struct {
		name        string
		projectType string
	}{
		{name: "go.mod", projectType: "go"},
		{name: "package.json", projectType: "nodejs"},
		{name: "package-lock.json", projectType: "nodejs"},
		{name: "yarn.lock", projectType: "nodejs"},
		{name: "requirements.txt", projectType: "python"},
		{name: "Pipfile", projectType: "python"},
		{name: "Pipfile.lock", projectType: "python"},
		{name: "Cargo.toml", projectType: "rust"},
		{name: "Cargo.lock", projectType: "rust"},
		{name: "pom.xml", projectType: "java"},
		{name: "build.gradle", projectType: "java"},
	}

	for _, manifest := range manifestFiles {
		if _, err := os.Stat(filepath.Join(repoPath, manifest.name)); err == nil {
			return manifest.projectType
		}
	}

	return ""
}

// buildDependencyTree builds the complete dependency tree for the project
func (c *TransitiveDependencyChecker) buildDependencyTree(repoPath string, projectType string) (*DependencyTree, error) {
	tree := &DependencyTree{
		Root: &Dependency{
			Name:     "project",
			Version:  "1.0.0",
			Direct:   true,
			Children: []*Dependency{},
		},
	}

	switch projectType {
	case "go":
		return c.buildGoDependencyTree(repoPath, tree)
	case "nodejs":
		return c.buildNodeJSDependencyTree(repoPath, tree)
	case "python":
		return c.buildPythonDependencyTree(repoPath, tree)
	case "rust":
		return c.buildRustDependencyTree(repoPath, tree)
	case "java":
		return c.buildJavaDependencyTree(repoPath, tree)
	default:
		return tree, fmt.Errorf("unsupported project type: %s", projectType)
	}
}

// buildGoDependencyTree builds dependency tree for Go projects
func (c *TransitiveDependencyChecker) buildGoDependencyTree(repoPath string, tree *DependencyTree) (*DependencyTree, error) {
	// Run go list -m all to get all dependencies
	cmd := exec.Command("go", "list", "-m", "all")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return tree, fmt.Errorf("failed to run go list: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	dependencies := make(map[string]*Dependency)

	for _, line := range lines {
		if line == "" || strings.Contains(line, "go:") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		version := parts[1]

		dep := &Dependency{
			Name:     name,
			Version:  version,
			Direct:   false, // Will be updated based on go.mod
			Children: []*Dependency{},
		}

		dependencies[name] = dep
		tree.Total++
	}

	// Parse go.mod to identify direct dependencies
	goModPath := filepath.Join(repoPath, "go.mod")
	if content, err := os.ReadFile(goModPath); err == nil {
		c.parseGoModDirectDeps(string(content), dependencies)
	}

	// Add all dependencies as children of root
	for _, dep := range dependencies {
		tree.Root.Children = append(tree.Root.Children, dep)
	}

	return tree, nil
}

// parseGoModDirectDeps parses go.mod to identify direct dependencies
func (c *TransitiveDependencyChecker) parseGoModDirectDeps(content string, dependencies map[string]*Dependency) {
	lines := strings.Split(content, "\n")
	inRequire := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "require (" {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}
		if strings.HasPrefix(line, "require ") && !inRequire {
			// Single line require
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[1]
				if dep, exists := dependencies[name]; exists {
					dep.Direct = true
				}
			}
		}
		if inRequire && !strings.HasPrefix(line, "//") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[0]
				if dep, exists := dependencies[name]; exists {
					dep.Direct = true
				}
			}
		}
	}
}

// buildNodeJSDependencyTree builds dependency tree for Node.js projects
func (c *TransitiveDependencyChecker) buildNodeJSDependencyTree(repoPath string, tree *DependencyTree) (*DependencyTree, error) {
	// Try npm ls first
	cmd := exec.Command("npm", "ls", "--json")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		// Fallback to parsing package-lock.json
		return c.parsePackageLockJson(repoPath, tree)
	}

	var npmTree map[string]interface{}
	if err := json.Unmarshal(output, &npmTree); err != nil {
		return c.parsePackageLockJson(repoPath, tree)
	}

	c.parseNpmTree(npmTree, tree.Root, []string{})
	return tree, nil
}

// parseNpmTree recursively parses npm dependency tree
func (c *TransitiveDependencyChecker) parseNpmTree(node map[string]interface{}, parent *Dependency, path []string) {
	if dependencies, ok := node["dependencies"].(map[string]interface{}); ok {
		for name, depInfo := range dependencies {
			if depMap, ok := depInfo.(map[string]interface{}); ok {
				version := "unknown"
				if v, exists := depMap["version"].(string); exists {
					version = v
				}

				dep := &Dependency{
					Name:     name,
					Version:  version,
					Direct:   len(path) == 0,
					Path:     append(path, name),
					Children: []*Dependency{},
				}

				parent.Children = append(parent.Children, dep)
				// Note: tree.Total is updated in the calling function

				// Recursively parse children
				c.parseNpmTree(depMap, dep, append(path, name))
			}
		}
	}
}

// parsePackageLockJson parses package-lock.json for dependency information
func (c *TransitiveDependencyChecker) parsePackageLockJson(repoPath string, tree *DependencyTree) (*DependencyTree, error) {
	lockPath := filepath.Join(repoPath, "package-lock.json")
	content, err := os.ReadFile(lockPath)
	if err != nil {
		return tree, fmt.Errorf("package-lock.json not found: %v", err)
	}

	var lockFile map[string]interface{}
	if err := json.Unmarshal(content, &lockFile); err != nil {
		return tree, fmt.Errorf("failed to parse package-lock.json: %v", err)
	}

	// Parse dependencies
	if deps, ok := lockFile["dependencies"].(map[string]interface{}); ok {
		c.parsePackageLockDeps(deps, tree.Root, []string{})
	}

	return tree, nil
}

// parsePackageLockDeps recursively parses package-lock.json dependencies
func (c *TransitiveDependencyChecker) parsePackageLockDeps(deps map[string]interface{}, parent *Dependency, path []string) {
	for name, depInfo := range deps {
		if depMap, ok := depInfo.(map[string]interface{}); ok {
			version := "unknown"
			if v, exists := depMap["version"].(string); exists {
				version = v
			}

			dep := &Dependency{
				Name:     name,
				Version:  version,
				Direct:   len(path) == 0,
				Path:     append(path, name),
				Children: []*Dependency{},
			}

			parent.Children = append(parent.Children, dep)
			// Note: tree.Total is updated in the calling function

			// Parse nested dependencies
			if nestedDeps, ok := depMap["dependencies"].(map[string]interface{}); ok {
				c.parsePackageLockDeps(nestedDeps, dep, append(path, name))
			}
		}
	}
}

// buildPythonDependencyTree builds dependency tree for Python projects
func (c *TransitiveDependencyChecker) buildPythonDependencyTree(repoPath string, tree *DependencyTree) (*DependencyTree, error) {
	// Try pipdeptree first
	cmd := exec.Command("pipdeptree", "--json")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		// Fallback to parsing requirements.txt
		return c.parseRequirementsTxt(repoPath, tree)
	}

	var pipTree []map[string]interface{}
	if err := json.Unmarshal(output, &pipTree); err != nil {
		return c.parseRequirementsTxt(repoPath, tree)
	}

	for _, dep := range pipTree {
		if name, ok := dep["package_name"].(string); ok {
			version := "unknown"
			if v, exists := dep["installed_version"].(string); exists {
				version = v
			}

			dependency := &Dependency{
				Name:     name,
				Version:  version,
				Direct:   true, // pipdeptree shows direct dependencies
				Children: []*Dependency{},
			}

			tree.Root.Children = append(tree.Root.Children, dependency)
			tree.Total++
		}
	}

	return tree, nil
}

// parseRequirementsTxt parses requirements.txt for dependencies
func (c *TransitiveDependencyChecker) parseRequirementsTxt(repoPath string, tree *DependencyTree) (*DependencyTree, error) {
	reqPath := filepath.Join(repoPath, "requirements.txt")
	content, err := os.ReadFile(reqPath)
	if err != nil {
		return tree, fmt.Errorf("requirements.txt not found: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse package==version format
		parts := strings.Split(line, "==")
		if len(parts) == 2 {
			name := strings.TrimSpace(parts[0])
			version := strings.TrimSpace(parts[1])

			dep := &Dependency{
				Name:     name,
				Version:  version,
				Direct:   true,
				Children: []*Dependency{},
			}

			tree.Root.Children = append(tree.Root.Children, dep)
			tree.Total++
		}
	}

	return tree, nil
}

// buildRustDependencyTree builds dependency tree for Rust projects
func (c *TransitiveDependencyChecker) buildRustDependencyTree(repoPath string, tree *DependencyTree) (*DependencyTree, error) {
	cmd := exec.Command("cargo", "metadata", "--format-version", "1")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return tree, fmt.Errorf("failed to run cargo metadata: %v", err)
	}

	var metadata struct {
		Packages []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"packages"`
		WorkspaceMembers []string `json:"workspace_members"`
	}
	if err := json.Unmarshal(output, &metadata); err != nil {
		return tree, fmt.Errorf("failed to parse cargo metadata output: %v", err)
	}

	workspace := make(map[string]bool, len(metadata.WorkspaceMembers))
	for _, id := range metadata.WorkspaceMembers {
		workspace[id] = true
	}
	for _, pkg := range metadata.Packages {
		if workspace[pkg.ID] {
			continue
		}
		tree.Root.Children = append(tree.Root.Children, &Dependency{
			Name:     pkg.Name,
			Version:  pkg.Version,
			Direct:   false,
			Children: []*Dependency{},
		})
		tree.Total++
	}

	return tree, nil
}

// buildJavaDependencyTree builds dependency tree for Java projects
func (c *TransitiveDependencyChecker) buildJavaDependencyTree(repoPath string, tree *DependencyTree) (*DependencyTree, error) {
	// Try Maven dependency tree
	cmd := exec.Command("mvn", "dependency:tree", "-DoutputType=json")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		// Fallback to parsing pom.xml
		return c.parsePomXml(repoPath, tree)
	}

	var mavenTree map[string]interface{}
	if err := json.Unmarshal(output, &mavenTree); err != nil {
		return c.parsePomXml(repoPath, tree)
	}

	// Parse Maven dependency tree
	c.parseMavenTree(mavenTree, tree.Root, []string{})
	return tree, nil
}

// parseMavenTree recursively parses Maven dependency tree
func (c *TransitiveDependencyChecker) parseMavenTree(node map[string]interface{}, parent *Dependency, path []string) {
	if dependencies, ok := node["dependencies"].([]interface{}); ok {
		for _, depInfo := range dependencies {
			if depMap, ok := depInfo.(map[string]interface{}); ok {
				name := ""
				version := "unknown"

				if n, exists := depMap["groupId"].(string); exists {
					name = n
				}
				if a, exists := depMap["artifactId"].(string); exists {
					if name != "" {
						name += ":" + a
					} else {
						name = a
					}
				}
				if v, exists := depMap["version"].(string); exists {
					version = v
				}

				if name != "" {
					dep := &Dependency{
						Name:     name,
						Version:  version,
						Direct:   len(path) == 0,
						Path:     append(path, name),
						Children: []*Dependency{},
					}

					parent.Children = append(parent.Children, dep)
					// Note: tree.Total is updated in the calling function

					// Recursively parse children
					c.parseMavenTree(depMap, dep, append(path, name))
				}
			}
		}
	}
}

// parsePomXml parses pom.xml for dependencies
func (c *TransitiveDependencyChecker) parsePomXml(repoPath string, tree *DependencyTree) (*DependencyTree, error) {
	pomPath := filepath.Join(repoPath, "pom.xml")
	content, err := os.ReadFile(pomPath)
	if err != nil {
		return tree, fmt.Errorf("pom.xml not found: %v", err)
	}

	// Simple XML parsing for dependencies
	lines := strings.Split(string(content), "\n")
	inDependencies := false

	for i, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "<dependencies>") {
			inDependencies = true
			continue
		}
		if strings.Contains(line, "</dependencies>") {
			inDependencies = false
			continue
		}

		if inDependencies && strings.Contains(line, "<dependency>") {
			// Parse dependency block
			dep := c.parseDependencyBlock(lines, &i)
			if dep != nil {
				tree.Root.Children = append(tree.Root.Children, dep)
				tree.Total++
			}
		}
	}

	return tree, nil
}

// parseDependencyBlock parses a single dependency block from pom.xml
func (c *TransitiveDependencyChecker) parseDependencyBlock(lines []string, currentLine *int) *Dependency {
	dep := &Dependency{
		Children: []*Dependency{},
		Direct:   true,
	}

	for *currentLine < len(lines) {
		line := strings.TrimSpace(lines[*currentLine])
		*currentLine++

		if strings.Contains(line, "</dependency>") {
			break
		}

		if strings.Contains(line, "<groupId>") {
			groupId := c.extractXmlValue(line)
			if dep.Name == "" {
				dep.Name = groupId
			} else {
				dep.Name = groupId + ":" + dep.Name
			}
		}
		if strings.Contains(line, "<artifactId>") {
			artifactId := c.extractXmlValue(line)
			if dep.Name == "" {
				dep.Name = artifactId
			} else {
				dep.Name = dep.Name + ":" + artifactId
			}
		}
		if strings.Contains(line, "<version>") {
			dep.Version = c.extractXmlValue(line)
		}
	}

	if dep.Name != "" {
		return dep
	}
	return nil
}

// extractXmlValue extracts value from XML tag
func (c *TransitiveDependencyChecker) extractXmlValue(line string) string {
	start := strings.Index(line, ">")
	end := strings.Index(line, "</")
	if start != -1 && end != -1 && start < end {
		return line[start+1 : end]
	}
	return ""
}

// getSeverityLevel returns numeric severity level
func (c *TransitiveDependencyChecker) getSeverityLevel(severity string) int {
	switch strings.ToLower(severity) {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}

// updateTreeCounts updates vulnerability counts in the tree
func (c *TransitiveDependencyChecker) updateTreeCounts(tree *DependencyTree) {
	c.updateCounts(tree.Root, tree)
}

// updateCounts recursively updates counts
func (c *TransitiveDependencyChecker) updateCounts(dep *Dependency, tree *DependencyTree) {
	if dep.Vulnerable {
		tree.Vulnerable++
		switch dep.Severity {
		case "critical":
			tree.Critical++
		case "high":
			tree.High++
		case "medium":
			tree.Medium++
		case "low":
			tree.Low++
		}
	}

	for _, child := range dep.Children {
		c.updateCounts(child, tree)
	}
}

// calculateScore calculates security score based on vulnerabilities
func (c *TransitiveDependencyChecker) calculateScore(tree *DependencyTree) int {
	if tree.Total == 0 {
		return 100
	}

	// Base score
	score := 100

	// Deduct points based on vulnerability severity
	score -= tree.Critical * 20 // -20 points per critical
	score -= tree.High * 10     // -10 points per high
	score -= tree.Medium * 5    // -5 points per medium
	score -= tree.Low * 2       // -2 points per low

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

// GetCategory returns the category of this checker
func (c *TransitiveDependencyChecker) GetCategory() string {
	return "Security"
}
