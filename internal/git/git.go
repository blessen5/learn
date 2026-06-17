package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsRepo checks if the given path is inside a git repository.
func IsRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

// Init runs git init in the given path.
func Init(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Add stages a file.
func Add(repoRoot, filepath string) error {
	cmd := exec.Command("git", "add", filepath)
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// AddFiles stages specific files.
func AddFiles(repoRoot string, files []string) error {
	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// AddAll stages all changes.
func AddAll(repoRoot string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Commit creates a commit with the given message.
func Commit(repoRoot, message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Status returns modified/untracked files.
func Status(repoRoot string) ([]string, error) {
	if !IsRepo(repoRoot) {
		return nil, fmt.Errorf("not a git repository: %s", repoRoot)
	}

	// Use git diff and ls-files for reliable parsing
	var files []string

	// Modified/staged files
	cmd := exec.Command("git", "diff", "--name-only", "--cached")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if line != "" {
				files = append(files, line)
			}
		}
	}

	// Unstaged modifications
	cmd = exec.Command("git", "diff", "--name-only")
	cmd.Dir = repoRoot
	out, err = cmd.Output()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if line != "" {
				files = append(files, line)
			}
		}
	}

	// Untracked files
	cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard")
	cmd.Dir = repoRoot
	out, err = cmd.Output()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if line != "" {
				files = append(files, line)
			}
		}
	}

	return files, nil
}

// IsAvailable checks if git is installed.
func IsAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
