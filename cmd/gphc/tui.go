package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/vahidaghazadeh/gphc/pkg/types"
)

func runTUI(cmd *cobra.Command, args []string) {
	var repoPath string
	if len(args) > 0 {
		repoPath = args[0]
	} else {
		var err error
		repoPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check if it's a Git repository
	if !isGitRepository(repoPath) {
		fmt.Printf("Error: %s is not a Git repository\n", repoPath)
		os.Exit(1)
	}

	// Create TUI model
	model := NewTUIModel(repoPath)

	// Run the TUI
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

// TUIModel represents the TUI application state
type TUIModel struct {
	repoPath     string
	healthReport *types.HealthReport
	loading      bool
	err          error
	selectedTab  int
	tabs         []string
}

// NewTUIModel creates a new TUI model
func NewTUIModel(repoPath string) *TUIModel {
	return &TUIModel{
		repoPath:    repoPath,
		loading:     true,
		tabs:        []string{"Overview", "Details"},
		selectedTab: 0,
	}
}

// Init implements the tea.Model interface
func (m *TUIModel) Init() tea.Cmd {
	return m.loadHealthData()
}

// Update implements the tea.Model interface
func (m *TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.selectedTab = (m.selectedTab + 1) % len(m.tabs)
		case "shift+tab":
			m.selectedTab = (m.selectedTab - 1 + len(m.tabs)) % len(m.tabs)
		case "r":
			m.loading = true
			return m, m.loadHealthData()
		}
	case healthDataMsg:
		m.healthReport = msg.report
		m.loading = false
		m.err = msg.err
	}
	return m, nil
}

// View implements the tea.Model interface
func (m *TUIModel) View() string {
	if m.loading {
		return "Loading health data...\n\nPress 'q' to quit, 'r' to refresh"
	}

	if m.err != nil {
		return fmt.Sprintf("Error loading health data: %v\n\nPress 'q' to quit, 'r' to refresh", m.err)
	}

	var content string
	switch m.selectedTab {
	case 0:
		content = m.renderOverview()
	case 1:
		content = m.renderDetails()
	}

	return fmt.Sprintf("GPHC TUI - %s\n\n%s\n\nPress 'q' to quit, 'tab' to switch tabs, 'r' to refresh",
		filepath.Base(m.repoPath), content)
}

// healthDataMsg is used to pass health data to the model
type healthDataMsg struct {
	report *types.HealthReport
	err    error
}

// loadHealthData loads the health data for the repository
func (m *TUIModel) loadHealthData() tea.Cmd {
	return func() tea.Msg {
		healthReport, err := buildHealthReport(m.repoPath)
		if err != nil {
			return healthDataMsg{err: err}
		}
		return healthDataMsg{report: healthReport}
	}
}

// renderOverview renders the overview tab
func (m *TUIModel) renderOverview() string {
	if m.healthReport == nil {
		return "No data available"
	}

	return fmt.Sprintf(`Health Score: %d/100 (%s)

Total Checks: %d
Passed: %d
Failed: %d
Warnings: %d

Repository: %s
Last Updated: %s`,
		m.healthReport.OverallScore,
		m.healthReport.Grade,
		m.healthReport.Summary.TotalChecks,
		m.healthReport.Summary.PassedChecks,
		m.healthReport.Summary.FailedChecks,
		m.healthReport.Summary.WarningChecks,
		filepath.Base(m.repoPath),
		m.healthReport.Timestamp.Format("2006-01-02 15:04:05"))
}

// renderDetails renders the details tab
func (m *TUIModel) renderDetails() string {
	if m.healthReport == nil {
		return "No data available"
	}

	var details strings.Builder
	details.WriteString("Check Results:\n\n")

	for _, result := range m.healthReport.Results {
		status := "PASS"
		if result.Status == types.StatusFail {
			status = "FAIL"
		} else if result.Status == types.StatusWarning {
			status = "WARN"
		}

		details.WriteString(fmt.Sprintf("%s [%s] %s\n", status, result.ID, result.Name))
		if result.Message != "" {
			details.WriteString(fmt.Sprintf("  %s\n", result.Message))
		}
		details.WriteString("\n")
	}

	return details.String()
}
