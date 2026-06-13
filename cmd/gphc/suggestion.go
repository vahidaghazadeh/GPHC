package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func analyzeStagedChanges(repoPath string, stagedFiles []string) string {
	changes, err := getStagedChanges(repoPath)
	if err != nil || len(changes) == 0 {
		changes = make([]stagedChange, 0, len(stagedFiles))
		for _, file := range stagedFiles {
			changes = append(changes, stagedChange{status: "M", path: file})
		}
	}

	return buildCommitSuggestion(changes)
}

type stagedChange struct {
	status string
	path   string
}

func getStagedChanges(repoPath string) ([]stagedChange, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-status")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var changes []stagedChange
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		status := parts[0][:1]
		path := parts[len(parts)-1]
		changes = append(changes, stagedChange{status: status, path: path})
	}
	return changes, nil
}

func buildCommitSuggestion(changes []stagedChange) string {
	if len(changes) == 0 {
		return "chore: update project files"
	}

	commitType := inferCommitType(changes)
	scope := inferCommitScope(changes)
	subject := inferCommitSubject(changes, commitType)

	if scope != "" {
		return fmt.Sprintf("%s(%s): %s", commitType, scope, subject)
	}
	return fmt.Sprintf("%s: %s", commitType, subject)
}

func inferCommitType(changes []stagedChange) string {
	allDocs, allTests, allCI, allDependencies, allDeployment := true, true, true, true, true
	allAdded, allDeleted := true, true

	for _, change := range changes {
		allDocs = allDocs && isDocumentationPath(change.path)
		allTests = allTests && isTestPath(change.path)
		allCI = allCI && isCIPath(change.path)
		allDependencies = allDependencies && isDependencyPath(change.path)
		allDeployment = allDeployment && isDeploymentPath(change.path)
		allAdded = allAdded && change.status == "A"
		allDeleted = allDeleted && change.status == "D"
	}

	switch {
	case allDocs:
		return "docs"
	case allTests:
		return "test"
	case allCI:
		return "ci"
	case allDependencies:
		return "build"
	case allDeployment:
		return "chore"
	case allAdded:
		return "feat"
	case allDeleted:
		return "refactor"
	default:
		return "chore"
	}
}

func inferCommitScope(changes []stagedChange) string {
	if allChangesMatch(changes, isDeploymentPath) {
		return "deploy"
	}

	scope := scopeForPath(changes[0].path)
	if scope == "" {
		return ""
	}

	for _, change := range changes[1:] {
		if scopeForPath(change.path) != scope {
			return ""
		}
	}
	return scope
}

func scopeForPath(path string) string {
	cleanPath := filepath.ToSlash(path)
	parts := strings.Split(cleanPath, "/")
	if len(parts) == 1 {
		return ""
	}

	switch parts[0] {
	case "cmd":
		if len(parts) > 1 {
			return sanitizeScope(parts[1])
		}
	case "internal", "pkg", "src", "lib":
		if len(parts) > 1 {
			return sanitizeScope(parts[1])
		}
	case ".github":
		return "github"
	case "docs":
		return ""
	default:
		return sanitizeScope(parts[0])
	}
	return ""
}

func sanitizeScope(value string) string {
	value = strings.TrimSuffix(strings.ToLower(value), filepath.Ext(value))
	value = strings.NewReplacer("_", "-", " ", "-").Replace(value)
	return strings.Trim(value, "-")
}

func inferCommitSubject(changes []stagedChange, commitType string) string {
	if allChangesMatch(changes, isDeploymentPath) {
		return describeDeploymentChanges(changes)
	}

	if len(changes) == 1 {
		change := changes[0]
		target := describePath(change.path)
		switch change.status {
		case "A":
			return "add " + target
		case "D":
			return "remove " + target
		default:
			switch commitType {
			case "docs":
				return "update " + target
			case "test":
				return "update " + target
			case "ci":
				return "update " + target
			case "build":
				return "update " + target
			default:
				return "improve " + target
			}
		}
	}

	switch commitType {
	case "docs":
		return "update project documentation"
	case "test":
		return "update test coverage"
	case "ci":
		return "update CI workflows"
	case "build":
		return "update project dependencies"
	case "feat":
		return "add staged functionality"
	case "refactor":
		return "remove obsolete project files"
	default:
		if scope := inferCommitScope(changes); scope != "" {
			return "update " + strings.ReplaceAll(scope, "-", " ") + " implementation"
		}
		return "update related project files"
	}
}

func allChangesMatch(changes []stagedChange, match func(string) bool) bool {
	if len(changes) == 0 {
		return false
	}
	for _, change := range changes {
		if !match(change.path) {
			return false
		}
	}
	return true
}

func describeDeploymentChanges(changes []stagedChange) string {
	var hasKubernetes, hasProxy, hasSecret, hasEnvironment bool
	var deploymentName, proxyName string
	for _, change := range changes {
		lower := strings.ToLower(filepath.ToSlash(change.path))
		base := strings.ToLower(filepath.Base(lower))

		isKubernetes :=
			strings.Contains(lower, "/statefulset/") ||
				strings.Contains(lower, "/deployment/") ||
				strings.Contains(lower, "/kubernetes/") ||
				strings.Contains(lower, "/k8s/") ||
				strings.Contains(lower, "/helm/")
		isProxy :=
			strings.Contains(lower, "nginx") ||
				strings.Contains(lower, "sites-available") ||
				strings.Contains(lower, "reverseproxy") ||
				strings.Contains(lower, "reversproxy")

		hasKubernetes = hasKubernetes || isKubernetes
		hasProxy = hasProxy || isProxy
		if isKubernetes && deploymentName == "" {
			deploymentName = topLevelComponent(change.path)
		}
		if isProxy && proxyName == "" {
			proxyName = proxyTarget(change.path)
		}
		hasSecret = hasSecret || strings.HasPrefix(base, "secret.")
		hasEnvironment = hasEnvironment ||
			strings.HasPrefix(base, "env.") ||
			strings.HasPrefix(base, ".env")
	}

	switch {
	case hasKubernetes && hasProxy:
		if deploymentName != "" && proxyName != "" {
			return fmt.Sprintf("update %s deployment and %s proxy config", deploymentName, proxyName)
		}
		return "update deployment and reverse proxy configs"
	case hasSecret && hasEnvironment:
		return "update deployment secrets and environment config"
	case hasSecret:
		return "update deployment secrets"
	case hasEnvironment:
		return "update deployment environment config"
	case hasKubernetes:
		return "update Kubernetes deployment config"
	case hasProxy:
		return "update reverse proxy config"
	default:
		return "update deployment config"
	}
}

func topLevelComponent(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) == 0 {
		return ""
	}
	return humanizeIdentifier(parts[0])
}

func proxyTarget(path string) string {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '-' || r == '_' || r == '.'
	})
	if len(parts) == 0 {
		return ""
	}

	// Vhost names commonly end in the proxied service, for example iwcs-kibana.
	return humanizeIdentifier(parts[len(parts)-1])
}

func humanizeIdentifier(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if value == strings.ToLower(value) || value == strings.ToUpper(value) {
		return strings.ToUpper(value[:1]) + strings.ToLower(value[1:])
	}
	return value
}

func describePath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	lowerBase := strings.ToLower(base)

	switch lowerBase {
	case "readme.md", "readme.rst", "readme.txt":
		return "README"
	case "go.mod", "go.sum":
		return "Go dependencies"
	case "package.json", "package-lock.json", "yarn.lock", "pnpm-lock.yaml":
		return "JavaScript dependencies"
	case "dockerfile":
		return "Docker configuration"
	}

	name = strings.NewReplacer("_", " ", "-", " ", ".", " ").Replace(name)
	name = strings.Join(strings.Fields(name), " ")
	if name == "" {
		return "project file"
	}

	if isTestPath(path) && !strings.Contains(strings.ToLower(name), "test") {
		return name + " tests"
	}
	return name
}

func isDocumentationPath(path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	ext := strings.ToLower(filepath.Ext(lower))
	base := filepath.Base(lower)
	return strings.HasPrefix(lower, "docs/") ||
		strings.HasPrefix(base, "readme") ||
		strings.HasPrefix(base, "changelog") ||
		ext == ".md" || ext == ".rst" || ext == ".adoc"
}

func isTestPath(path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	base := filepath.Base(lower)
	return strings.Contains(lower, "/test/") ||
		strings.Contains(lower, "/tests/") ||
		strings.HasSuffix(base, "_test.go") ||
		strings.Contains(base, ".test.") ||
		strings.Contains(base, ".spec.")
}

func isCIPath(path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	return strings.HasPrefix(lower, ".github/workflows/") ||
		strings.HasPrefix(lower, ".gitlab/") ||
		lower == ".gitlab-ci.yml" ||
		lower == "jenkinsfile"
}

func isDependencyPath(path string) bool {
	switch strings.ToLower(filepath.Base(path)) {
	case "go.mod", "go.sum", "package.json", "package-lock.json", "yarn.lock",
		"pnpm-lock.yaml", "cargo.toml", "cargo.lock", "requirements.txt",
		"poetry.lock", "composer.json", "composer.lock", "gemfile", "gemfile.lock":
		return true
	default:
		return false
	}
}

func isDeploymentPath(path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	base := strings.ToLower(filepath.Base(lower))

	if strings.Contains(lower, "/statefulset/") ||
		strings.Contains(lower, "/deployment/") ||
		strings.Contains(lower, "/kubernetes/") ||
		strings.Contains(lower, "/k8s/") ||
		strings.Contains(lower, "/helm/") ||
		strings.Contains(lower, "sites-available") ||
		strings.Contains(lower, "nginx") ||
		strings.Contains(lower, "reverseproxy") ||
		strings.Contains(lower, "reversproxy") {
		return true
	}

	switch base {
	case "dockerfile", "docker-compose.yml", "docker-compose.yaml",
		"compose.yml", "compose.yaml", "values.yml", "values.yaml":
		return true
	default:
		return strings.HasPrefix(base, "env.") || strings.HasPrefix(base, "secret.")
	}
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}
