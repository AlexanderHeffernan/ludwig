package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// CreateWorktree creates a new worktree for the branch
func CreateWorktree(branchName string) error {
	worktreePath := "../worktrees/" + branchName
	cmd := exec.Command("git", "worktree", "add", worktreePath, branchName)
	cmd.Dir = getRepoRoot()
	return cmd.Run()
}

// WorktreeExists checks if a worktree already exists for the branch
func WorktreeExists(branchName string) (bool, error) {
	cmd := exec.Command("git", "worktree", "list")
	cmd.Dir = getRepoRoot()
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	// Check if branchName appears in the output
	return strings.Contains(string(output), branchName), nil
}

// GenerateBranchName creates a unique branch name from a task name
// Extracts first 2-3 words, converts to kebab-case, ensures uniqueness
func GenerateBranchName(taskName string) (string, error) {
	// Extract first 2-3 meaningful words
	words := extractWords(taskName)
	if len(words) == 0 {
		return "", fmt.Errorf("unable to generate branch name from task: %s", taskName)
	}

	// Take first 2-3 words
	maxWords := 3
	if len(words) < maxWords {
		maxWords = len(words)
	}
	baseName := strings.Join(words[:maxWords], "-")

	// Convert to kebab-case and limit length
	baseName = toKebabCase(baseName)
	if len(baseName) > 40 {
		baseName = baseName[:40]
	}

	// Add ludwig/ prefix
	baseName = "ludwig/" + baseName

	// Check for duplicates and append counter if needed
	branchName := baseName
	counter := 1
	for {
		exists, err := WorktreeExists(branchName)
		if err != nil {
			return "", err
		}
		if !exists {
			return branchName, nil
		}
		branchName = fmt.Sprintf("%s-%d", baseName, counter)
		counter++
	}
}

// extractWords extracts meaningful words from a task name
func extractWords(taskName string) []string {
	// Split by spaces and special chars
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	parts := re.Split(taskName, -1)

	var words []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && len(part) > 2 { // Skip single letters and numbers only
			words = append(words, strings.ToLower(part))
		}
	}
	return words
}

// toKebabCase converts a string to kebab-case
func toKebabCase(s string) string {
	// Replace underscores and spaces with hyphens
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")

	// Remove non-alphanumeric except hyphens
	re := regexp.MustCompile("[^a-z0-9-]")
	s = re.ReplaceAllString(s, "")

	// Remove multiple consecutive hyphens
	re = regexp.MustCompile("-+")
	s = re.ReplaceAllString(s, "-")

	// Trim hyphens from start/end
	s = strings.Trim(s, "-")

	return s
}

// getRepoRoot returns the root directory of the repository
func getRepoRoot() string {
	// For now, assume current working directory is within the repo
	// In production, you might want to use git to find the root
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}
