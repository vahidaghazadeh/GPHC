package main

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// Simple test to ensure the main package compiles
	// This is a placeholder test - in a real project you'd have more comprehensive tests

	if version == "" {
		t.Error("Version should not be empty")
	}

	if effectiveVersion() == "" {
		t.Error("Effective version should not be empty")
	}
}

func TestMainFunction(t *testing.T) {
	// Test that main function exists and can be called
	// This is a basic smoke test

	// In a real test, you might test command line arguments
	// or other functionality, but for now this ensures the package builds
}

func TestBuildCommitSuggestion(t *testing.T) {
	tests := []struct {
		name    string
		changes []stagedChange
		want    string
	}{
		{
			name:    "single source file",
			changes: []stagedChange{{status: "M", path: "internal/checkers/binary_file_checker.go"}},
			want:    "chore(checkers): improve binary file checker",
		},
		{
			name:    "documentation",
			changes: []stagedChange{{status: "M", path: "README.md"}},
			want:    "docs: update README",
		},
		{
			name:    "new command",
			changes: []stagedChange{{status: "A", path: "cmd/report/main.go"}},
			want:    "feat(report): add main",
		},
		{
			name:    "tests",
			changes: []stagedChange{{status: "M", path: "internal/checkers/checkers_test.go"}},
			want:    "test(checkers): update checkers test",
		},
		{
			name: "shared package",
			changes: []stagedChange{
				{status: "M", path: "internal/reporter/reporter.go"},
				{status: "M", path: "internal/reporter/reporter_test.go"},
			},
			want: "chore(reporter): update reporter implementation",
		},
		{
			name:    "deleted file",
			changes: []stagedChange{{status: "D", path: "internal/legacy/adapter.go"}},
			want:    "refactor(legacy): remove adapter",
		},
		{
			name: "kubernetes and nginx configuration",
			changes: []stagedChange{
				{status: "M", path: "Laravel/Offl/Prod/Admin/StatefulSet/Secret.yaml"},
				{status: "M", path: "Laravel/Offl/Prod/Admin/StatefulSet/env.yaml"},
				{status: "M", path: "Laravel/Offl/Prod/Core/StatefulSet/env.yaml"},
				{status: "M", path: "ReversProxyNginx/sites-available/iwcs-kibana.conf"},
			},
			want: "chore(laravel): update Laravel deployment and Kibana proxy config",
		},
		{
			name: "deployment secrets and environment",
			changes: []stagedChange{
				{status: "M", path: "services/api/k8s/Secret.yaml"},
				{status: "M", path: "services/api/k8s/env.yaml"},
			},
			want: "chore(api): update API deployment secrets and environment config",
		},
		{
			name:    "nginx configuration",
			changes: []stagedChange{{status: "M", path: "nginx/sites-available/api.conf"}},
			want:    "chore(api): update API proxy config",
		},
		{
			name:    "clickhouse operator bundle",
			changes: []stagedChange{{status: "M", path: "Kubernetes/ClickHouse/clickhouse-operator-install-bundle.yaml"}},
			want:    "chore(clickhouse): update clickhouse operator install bundle",
		},
		{
			name: "signoz stack manifests",
			changes: []stagedChange{
				{status: "A", path: "ArgoCD/applications/infra/signoz.yaml"},
				{status: "A", path: "Signoz/Deployment/clickhouse-ha.yaml"},
				{status: "A", path: "Signoz/Deployment/clickhouse-operator-install-bundle.yaml"},
				{status: "A", path: "Signoz/Deployment/namespaces.yaml"},
				{status: "A", path: "Signoz/Deployment/signoz-otel-collector.yaml"},
				{status: "A", path: "Signoz/Deployment/signoz-secrets.yaml"},
				{status: "A", path: "Signoz/Deployment/signoz-values.yaml"},
				{status: "A", path: "Signoz/Deployment/traefik-ingressroutes.yaml"},
				{status: "A", path: "Signoz/Deployment/zookeeper.yaml"},
			},
			want: "feat(signoz): add Signoz deployment and ArgoCD application manifests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildCommitSuggestion(tt.changes); got != tt.want {
				t.Fatalf("buildCommitSuggestion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShellQuote(t *testing.T) {
	got := shellQuote("fix(cli): handle user's input")
	want := "'fix(cli): handle user'\"'\"'s input'"
	if got != want {
		t.Fatalf("shellQuote() = %q, want %q", got, want)
	}
}

func TestIsAffirmativeAnswer(t *testing.T) {
	yesAnswers := []string{"y", "Y", "yes", "YES", " yes \n", "بله", "آره", "اره"}
	for _, answer := range yesAnswers {
		if !isAffirmativeAnswer(answer) {
			t.Fatalf("isAffirmativeAnswer(%q) = false, want true", answer)
		}
	}

	noAnswers := []string{"", "n", "no", "خیر", "random"}
	for _, answer := range noAnswers {
		if isAffirmativeAnswer(answer) {
			t.Fatalf("isAffirmativeAnswer(%q) = true, want false", answer)
		}
	}
}

func TestFilterRepositories(t *testing.T) {
	repos := []string{"/work/api", "/work/web", "/work/legacy-api"}
	got := filterRepositories(repos, []string{"*api"}, []string{"legacy"})
	if len(got) != 1 || got[0] != "/work/api" {
		t.Fatalf("filterRepositories() = %#v", got)
	}
}
