package checkers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/opsource/gphc/pkg/types"
)

// BinaryFileChecker audits executable and large files in repository
type BinaryFileChecker struct {
	BaseChecker
}

// BinaryFile represents a detected binary or large file
type BinaryFile struct {
	Path        string  `json:"path"`
	Size        int64   `json:"size"`
	SizeMB      float64 `json:"size_mb"`
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	InGitignore bool    `json:"in_gitignore"`
	InHistory   bool    `json:"in_history"`
	Extension   string  `json:"extension"`
}

// BinaryAuditReport represents the complete binary file audit report
type BinaryAuditReport struct {
	ExecutableFiles []BinaryFile `json:"executable_files"`
	LargeFiles      []BinaryFile `json:"large_files"`
	SuspiciousFiles []BinaryFile `json:"suspicious_files"`
	TotalSize       int64        `json:"total_size"`
	TotalSizeMB     float64      `json:"total_size_mb"`
	FileCount       int          `json:"file_count"`
	Score           int          `json:"score"`
}

// NewBinaryFileChecker creates a new BinaryFileChecker
func NewBinaryFileChecker() *BinaryFileChecker {
	return &BinaryFileChecker{
		BaseChecker: BaseChecker{
			id:   "BINARY-AUDIT",
			name: "Executable & Large File Audit",
		},
	}
}

// Check performs binary and large file audit
func (c *BinaryFileChecker) Check(data *types.RepositoryData) *types.CheckResult {
	return c.CheckWithOptions(data, true, true, true, true, 10.0)
}

// CheckWithOptions performs binary and large file audit with specific options
func (c *BinaryFileChecker) CheckWithOptions(data *types.RepositoryData, checkExecutables, checkLarge, checkSuspicious, checkHistory bool, maxSizeMB float64) *types.CheckResult {
	result := &types.CheckResult{
		ID:        c.ID(),
		Name:      c.Name(),
		Status:    types.StatusPass,
		Score:     100,
		Message:   "No suspicious binary or large files found",
		Details:   []string{},
		Category:  c.Category(),
		Timestamp: time.Now(),
	}

	// Initialize audit report
	report := &BinaryAuditReport{
		ExecutableFiles: []BinaryFile{},
		LargeFiles:      []BinaryFile{},
		SuspiciousFiles: []BinaryFile{},
		TotalSize:       0,
		TotalSizeMB:     0,
		FileCount:       0,
	}

	// Check current working directory
	c.scanDirectoryForBinaryFiles(data.Path, report, checkExecutables, checkLarge, checkSuspicious, maxSizeMB)

	// Check Git history for binary files
	if checkHistory {
		c.checkGitHistoryForBinaryFiles(data.Path, report, checkExecutables, checkSuspicious)
	}

	// Calculate score based on findings
	score := c.calculateScore(report)
	result.Score = score

	// Update result based on findings
	totalIssues := len(report.ExecutableFiles) + len(report.LargeFiles) + len(report.SuspiciousFiles)
	if totalIssues > 0 {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Found %d binary/large file issues", totalIssues)
	} else {
		result.Status = types.StatusPass
		result.Message = "No suspicious binary or large files found"
	}

	// Add detailed information
	result.Details = append(result.Details, fmt.Sprintf("Executable Files: %d", len(report.ExecutableFiles)))
	result.Details = append(result.Details, fmt.Sprintf("Large Files: %d", len(report.LargeFiles)))
	result.Details = append(result.Details, fmt.Sprintf("Suspicious Files: %d", len(report.SuspiciousFiles)))
	result.Details = append(result.Details, fmt.Sprintf("Total Size: %.1f MB", report.TotalSizeMB))
	result.Details = append(result.Details, fmt.Sprintf("File Count: %d", report.FileCount))

	// Add summary table
	if len(report.ExecutableFiles) > 0 || len(report.LargeFiles) > 0 || len(report.SuspiciousFiles) > 0 {
		result.Details = append(result.Details, "")
		result.Details = append(result.Details, "📋 File Summary Table:")
		result.Details = append(result.Details, "┌─────────────────┬──────────┬──────────┬──────────┬──────────┐")
		result.Details = append(result.Details, "│ File Type       │ Count    │ Size     │ Severity │ Status   │")
		result.Details = append(result.Details, "├─────────────────┼──────────┼──────────┼──────────┼──────────┤")
		
		// Executable files summary
		if len(report.ExecutableFiles) > 0 {
			execSize := 0.0
			execCount := 0
			execSeverity := "low"
			for _, file := range report.ExecutableFiles {
				execSize += file.SizeMB
				execCount++
				if file.Severity == "critical" || file.Severity == "high" {
					execSeverity = file.Severity
				}
			}
			result.Details = append(result.Details, fmt.Sprintf("│ Executable      │ %-8d │ %-8.1f │ %-8s │ Active   │", execCount, execSize, execSeverity))
		}
		
		// Large files summary
		if len(report.LargeFiles) > 0 {
			largeSize := 0.0
			largeCount := 0
			largeSeverity := "low"
			for _, file := range report.LargeFiles {
				largeSize += file.SizeMB
				largeCount++
				if file.Severity == "critical" || file.Severity == "high" {
					largeSeverity = file.Severity
				}
			}
			result.Details = append(result.Details, fmt.Sprintf("│ Large Files     │ %-8d │ %-8.1f │ %-8s │ Active   │", largeCount, largeSize, largeSeverity))
		}
		
		// Suspicious files summary
		if len(report.SuspiciousFiles) > 0 {
			suspSize := 0.0
			suspCount := 0
			suspSeverity := "low"
			for _, file := range report.SuspiciousFiles {
				suspSize += file.SizeMB
				suspCount++
				if file.Severity == "critical" || file.Severity == "high" {
					suspSeverity = file.Severity
				}
			}
			result.Details = append(result.Details, fmt.Sprintf("│ Suspicious      │ %-8d │ %-8.1f │ %-8s │ Active   │", suspCount, suspSize, suspSeverity))
		}
		
		result.Details = append(result.Details, "└─────────────────┴──────────┴──────────┴──────────┴──────────┘")
		result.Details = append(result.Details, "")
		result.Details = append(result.Details, "💡 Use --format json for detailed file listings")
	}

	return result
}

// scanDirectoryForBinaryFiles scans current directory for binary and large files
func (c *BinaryFileChecker) scanDirectoryForBinaryFiles(repoPath string, report *BinaryAuditReport, checkExecutables, checkLarge, checkSuspicious bool, maxSizeMB float64) {
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip .git directory and other common directories
		if strings.Contains(path, ".git") || strings.Contains(path, "node_modules") || strings.Contains(path, "vendor") || strings.Contains(path, ".vscode") {
			return nil
		}

		// Only check files, not directories
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(repoPath, path)
		fileName := info.Name()
		fileSize := info.Size()
		fileExt := strings.ToLower(filepath.Ext(fileName))

		// Check for executable files
		if checkExecutables && c.isExecutableFile(fileName, fileExt) {
			binaryFile := BinaryFile{
				Path:        relPath,
				Size:        fileSize,
				SizeMB:      float64(fileSize) / (1024 * 1024),
				Type:        "executable",
				Severity:    c.getExecutableSeverity(fileExt),
				Description: c.getExecutableDescription(fileExt),
				InGitignore: c.isInGitignore(repoPath, relPath),
				InHistory:   false,
				Extension:   fileExt,
			}
			report.ExecutableFiles = append(report.ExecutableFiles, binaryFile)
			report.TotalSize += fileSize
			report.FileCount++
		}

		// Check for large files
		if checkLarge && fileSize > int64(maxSizeMB*1024*1024) {
			binaryFile := BinaryFile{
				Path:        relPath,
				Size:        fileSize,
				SizeMB:      float64(fileSize) / (1024 * 1024),
				Type:        "large",
				Severity:    c.getLargeFileSeverity(fileSize),
				Description: fmt.Sprintf("Large file: %.1f MB", float64(fileSize)/(1024*1024)),
				InGitignore: c.isInGitignore(repoPath, relPath),
				InHistory:   false,
				Extension:   fileExt,
			}
			report.LargeFiles = append(report.LargeFiles, binaryFile)
			report.TotalSize += fileSize
			report.FileCount++
		}

		// Check for suspicious files
		if checkSuspicious && c.isSuspiciousFile(fileName, fileExt) {
			binaryFile := BinaryFile{
				Path:        relPath,
				Size:        fileSize,
				SizeMB:      float64(fileSize) / (1024 * 1024),
				Type:        "suspicious",
				Severity:    "high",
				Description: fmt.Sprintf("Suspicious file type: %s", fileExt),
				InGitignore: c.isInGitignore(repoPath, relPath),
				InHistory:   false,
				Extension:   fileExt,
			}
			report.SuspiciousFiles = append(report.SuspiciousFiles, binaryFile)
			report.TotalSize += fileSize
			report.FileCount++
		}

		return nil
	})

	if err != nil {
		// Handle error silently
	}

	report.TotalSizeMB = float64(report.TotalSize) / (1024 * 1024)
}

// checkGitHistoryForBinaryFiles checks Git history for binary files
func (c *BinaryFileChecker) checkGitHistoryForBinaryFiles(repoPath string, report *BinaryAuditReport, checkExecutables, checkSuspicious bool) {
	// Get list of all files that have ever been tracked
	cmd := exec.Command("git", "log", "--name-only", "--pretty=format:", "--all")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return
	}

	files := strings.Split(string(output), "\n")
	seenFiles := make(map[string]bool)

	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" || seenFiles[file] {
			continue
		}
		seenFiles[file] = true

		fileName := filepath.Base(file)
		fileExt := strings.ToLower(filepath.Ext(fileName))

		// Check for executable files in history
		if checkExecutables && c.isExecutableFile(fileName, fileExt) {
			binaryFile := BinaryFile{
				Path:        file,
				Size:        0, // Size not available from git log
				SizeMB:      0,
				Type:        "executable",
				Severity:    c.getExecutableSeverity(fileExt),
				Description: c.getExecutableDescription(fileExt),
				InGitignore: c.isInGitignore(repoPath, file),
				InHistory:   true,
				Extension:   fileExt,
			}
			report.ExecutableFiles = append(report.ExecutableFiles, binaryFile)
			report.FileCount++
		}

		// Check for suspicious files in history
		if checkSuspicious && c.isSuspiciousFile(fileName, fileExt) {
			binaryFile := BinaryFile{
				Path:        file,
				Size:        0, // Size not available from git log
				SizeMB:      0,
				Type:        "suspicious",
				Severity:    "high",
				Description: fmt.Sprintf("Suspicious file type: %s", fileExt),
				InGitignore: c.isInGitignore(repoPath, file),
				InHistory:   true,
				Extension:   fileExt,
			}
			report.SuspiciousFiles = append(report.SuspiciousFiles, binaryFile)
			report.FileCount++
		}
	}
}

// isExecutableFile checks if a file is executable based on extension
func (c *BinaryFileChecker) isExecutableFile(fileName, fileExt string) bool {
	executableExtensions := map[string]bool{
		".exe":   true,
		".dll":   true,
		".so":    true,
		".dylib": true,
		".bin":   true,
		".jar":   true,
		".war":   true,
		".ear":   true,
		".app":   true,
		".deb":   true,
		".rpm":   true,
		".msi":   true,
		".pkg":   true,
		".dmg":   true,
		".iso":   true,
		".img":   true,
		".raw":   true,
		".vmdk":  true,
		".vdi":   true,
		".qcow2": true,
		".ova":   true,
		".ovf":   true,
	}

	return executableExtensions[fileExt]
}

// isSuspiciousFile checks if a file has suspicious characteristics
func (c *BinaryFileChecker) isSuspiciousFile(fileName, fileExt string) bool {
	// Check for suspicious extensions
	suspiciousExtensions := map[string]bool{
		".scr":     true,
		".bat":     true,
		".cmd":     true,
		".com":     true,
		".pif":     true,
		".vbs":     true,
		".js":      true,
		".jse":     true,
		".wsf":     true,
		".wsh":     true,
		".ps1":     true,
		".psm1":    true,
		".psd1":    true,
		".ps1xml":  true,
		".psc1":    true,
		".msh":     true,
		".msh1":    true,
		".msh2":    true,
		".mshxml":  true,
		".msh1xml": true,
		".msh2xml": true,
		".scf":     true,
		".lnk":     true,
		".inf":     true,
		".reg":     true,
		".doc":     true,
		".docx":    true,
		".xls":     true,
		".xlsx":    true,
		".ppt":     true,
		".pptx":    true,
	}

	// Check for suspicious patterns in filename
	suspiciousPatterns := []string{
		"malware",
		"virus",
		"trojan",
		"backdoor",
		"keylogger",
		"rootkit",
		"spyware",
		"adware",
		"ransomware",
		"botnet",
		"exploit",
		"payload",
		"inject",
		"bypass",
		"crack",
		"keygen",
		"patch",
		"hack",
		"cracked",
		"pirated",
	}

	// Check extension
	if suspiciousExtensions[fileExt] {
		return true
	}

	// Check filename patterns
	fileNameLower := strings.ToLower(fileName)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(fileNameLower, pattern) {
			return true
		}
	}

	return false
}

// getExecutableSeverity returns severity level for executable files
func (c *BinaryFileChecker) getExecutableSeverity(fileExt string) string {
	criticalExtensions := map[string]bool{
		".exe": true,
		".dll": true,
		".bin": true,
		".jar": true,
		".war": true,
		".ear": true,
		".msi": true,
		".pkg": true,
		".dmg": true,
	}

	if criticalExtensions[fileExt] {
		return "critical"
	}
	return "high"
}

// getExecutableDescription returns description for executable files
func (c *BinaryFileChecker) getExecutableDescription(fileExt string) string {
	descriptions := map[string]string{
		".exe":   "Windows executable file",
		".dll":   "Dynamic link library",
		".so":    "Shared object library",
		".dylib": "macOS dynamic library",
		".bin":   "Binary executable file",
		".jar":   "Java archive",
		".war":   "Web application archive",
		".ear":   "Enterprise application archive",
		".app":   "macOS application bundle",
		".deb":   "Debian package",
		".rpm":   "RPM package",
		".msi":   "Windows installer",
		".pkg":   "macOS package",
		".dmg":   "macOS disk image",
		".iso":   "ISO disk image",
		".img":   "Disk image",
		".raw":   "Raw disk image",
		".vmdk":  "VMware disk image",
		".vdi":   "VirtualBox disk image",
		".qcow2": "QEMU disk image",
		".ova":   "Open Virtual Appliance",
		".ovf":   "Open Virtualization Format",
	}

	if desc, exists := descriptions[fileExt]; exists {
		return desc
	}
	return fmt.Sprintf("Executable file: %s", fileExt)
}

// getLargeFileSeverity returns severity level for large files
func (c *BinaryFileChecker) getLargeFileSeverity(fileSize int64) string {
	sizeMB := float64(fileSize) / (1024 * 1024)

	if sizeMB > 100 {
		return "critical"
	} else if sizeMB > 50 {
		return "high"
	} else if sizeMB > 20 {
		return "medium"
	}
	return "low"
}

// isInGitignore checks if a file is ignored by .gitignore
func (c *BinaryFileChecker) isInGitignore(repoPath, filePath string) bool {
	cmd := exec.Command("git", "check-ignore", filePath)
	cmd.Dir = repoPath
	err := cmd.Run()
	return err == nil
}

// calculateScore calculates security score based on findings
func (c *BinaryFileChecker) calculateScore(report *BinaryAuditReport) int {
	score := 100

	// Deduct points for executable files
	for _, file := range report.ExecutableFiles {
		if !file.InGitignore {
			switch file.Severity {
			case "critical":
				score -= 25
			case "high":
				score -= 20
			case "medium":
				score -= 15
			case "low":
				score -= 10
			}
		}
	}

	// Deduct points for large files
	for _, file := range report.LargeFiles {
		if !file.InGitignore {
			switch file.Severity {
			case "critical":
				score -= 20
			case "high":
				score -= 15
			case "medium":
				score -= 10
			case "low":
				score -= 5
			}
		}
	}

	// Deduct points for suspicious files
	for _, file := range report.SuspiciousFiles {
		if !file.InGitignore {
			score -= 30
		}
	}

	// Deduct points for files in history
	for _, file := range report.ExecutableFiles {
		if file.InHistory {
			score -= 15
		}
	}

	for _, file := range report.SuspiciousFiles {
		if file.InHistory {
			score -= 20
		}
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

// GetCategory returns the category of this checker
func (c *BinaryFileChecker) GetCategory() string {
	return "Security"
}
