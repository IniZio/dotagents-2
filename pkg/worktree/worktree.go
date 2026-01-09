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

	cmd := exec.Command("git", "worktree", "add", wtPath, branch)
	cmd.Dir = m.RepoPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to create worktree: %w, output: %s", err, string(output))
	}

	return wtPath, nil
}

func (m *Manager) Remove(branch string) error {
	wtPath := filepath.Join(m.BaseDir, branch)
	
	cmd := exec.Command("git", "worktree", "remove", branch)
	cmd.Dir = m.RepoPath
	cmd.Run()

	return os.RemoveAll(wtPath)
}
