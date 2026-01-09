package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Manager struct {
	RepoPath string
	BaseDir  string
}

func NewManager(repoPath string, baseDir string) *Manager {
	return &Manager{
		RepoPath: repoPath,
		BaseDir:  baseDir,
	}
}

func (m *Manager) Add(branch string) (string, error) {
	wtPath := filepath.Join(m.BaseDir, branch)
	if _, err := os.Stat(wtPath); err == nil {
		return wtPath, nil
	}

	currentBranch, err := m.getCurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	if currentBranch == branch {
		safeBranches := []string{"main", "master"}
		switched := false
		for _, safeBranch := range safeBranches {
			if m.branchExists(safeBranch) {
				if err := m.switchToBranch(safeBranch); err == nil {
					switched = true
					break
				}
			}
		}
		if !switched {
			return "", fmt.Errorf("cannot create worktree: currently on branch %s and no safe branch (main/master) found", branch)
		}
	}

	if !m.branchExists(branch) {
		cmd := exec.Command("git", "checkout", "-b", branch)
		cmd.Dir = m.RepoPath
		if output, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("failed to create branch %s: %w, output: %s", branch, err, string(output))
		}
		if currentBranch != "" && currentBranch != branch {
			m.switchToBranch(currentBranch)
		}
	}

	cmd := exec.Command("git", "worktree", "add", wtPath, branch)
	cmd.Dir = m.RepoPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to create worktree: %w, output: %s", err, string(output))
	}

	return wtPath, nil
}

func (m *Manager) getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = m.RepoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output[:len(output)-1]), nil
}

func (m *Manager) branchExists(branch string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	cmd.Dir = m.RepoPath
	return cmd.Run() == nil
}

func (m *Manager) switchToBranch(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = m.RepoPath
	return cmd.Run()
}

func (m *Manager) Remove(branch string) error {
	wtPath := filepath.Join(m.BaseDir, branch)

	cmd := exec.Command("git", "worktree", "remove", branch)
	cmd.Dir = m.RepoPath
	cmd.Run()

	return os.RemoveAll(wtPath)
}
