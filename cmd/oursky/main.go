package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/vibegear/oursky/pkg/agent"
	"github.com/vibegear/oursky/pkg/ctrl"
	"github.com/vibegear/oursky/pkg/provider"
	"github.com/vibegear/oursky/pkg/provider/docker"
	"github.com/vibegear/oursky/pkg/templates"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func updateOursky() error {
	const repo = "oursky/vendatta"

	// Detect platform
	osName := runtime.GOOS
	arch := runtime.GOARCH

	var binaryName string
	switch osName {
	case "linux", "darwin":
		binaryName = fmt.Sprintf("oursky-%s-%s", osName, arch)
	case "windows":
		binaryName = fmt.Sprintf("oursky-%s-%s.exe", osName, arch)
	default:
		return fmt.Errorf("unsupported OS: %s", osName)
	}

	// Get latest release
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo))
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return fmt.Errorf("failed to parse release: %w", err)
	}

	fmt.Printf("Latest version: %s\n", release.TagName)

	// Download binary
	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, release.TagName, binaryName)
	fmt.Printf("Downloading from %s\n", downloadURL)

	resp, err = http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Get current binary path
	currentPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Create temp file in /tmp
	tempPath := fmt.Sprintf("/tmp/oursky-update-%d", os.Getpid())
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempPath)

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	tempFile.Close()

	// Make executable on Unix
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tempPath, 0755); err != nil {
			return fmt.Errorf("failed to make executable: %w", err)
		}
	}

	// Check if we can write to the directory
	dir := filepath.Dir(currentPath)
	testFile := filepath.Join(dir, ".oursky_write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("cannot write to %s. Please run with sudo or reinstall to a user-writable directory like ~/.local/bin", dir)
	}
	f.Close()
	os.Remove(testFile)

	// Backup current binary
	backupPath := currentPath + ".backup"
	if err := os.Rename(currentPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Replace with new binary
	if err := os.Rename(tempPath, currentPath); err != nil {
		// Try to restore backup
		os.Rename(backupPath, currentPath)
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	fmt.Printf("Successfully updated to %s\n", release.TagName)
	fmt.Printf("Backup saved at %s\n", backupPath)

	return nil
}

func main() {
	var providers []provider.Provider
	dProvider, err := docker.NewDockerProvider()
	if err == nil {
		providers = append(providers, dProvider)
	}

	controller := ctrl.NewBaseController(providers)

	rootCmd := &cobra.Command{
		Use:   "oursky",
		Short: "Oursky Dev Environment Manager",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize .oursky in the current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return controller.Init(context.Background())
		},
	}

	devCmd := &cobra.Command{
		Use:   "dev [branch]",
		Short: "Start a development session for a branch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return controller.Dev(context.Background(), args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List active sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			sessions, err := controller.List(context.Background())
			if err != nil {
				return err
			}
			for _, s := range sessions {
				fmt.Printf("%s\t%s\t%s\n", s.Labels["oursky.session.id"], s.Provider, s.Status)
			}
			return nil
		},
	}

	killCmd := &cobra.Command{
		Use:   "kill [session-id]",
		Short: "Stop and destroy a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return controller.Kill(context.Background(), args[0])
		},
	}

	agentCmd := &cobra.Command{
		Use:   "agent [session-id]",
		Short: "Start MCP agent gateway for a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sessions, _ := controller.List(context.Background())
			var targetSession *provider.Session
			for _, s := range sessions {
				if s.ID == args[0] || s.Labels["oursky.session.id"] == args[0] {
					targetSession = &s
					break
				}
			}
			if targetSession == nil {
				return fmt.Errorf("session %s not found", args[0])
			}

			p, ok := controller.Providers[targetSession.Provider]
			if !ok {
				return fmt.Errorf("provider %s not found", targetSession.Provider)
			}

			s := agent.NewAgentServer(targetSession.ID, p)
			return s.Serve()
		},
	}

	templatesCmd := &cobra.Command{
		Use:   "templates",
		Short: "Manage AI agent templates",
	}

	templatesPullCmd := &cobra.Command{
		Use:   "pull [url]",
		Short: "Pull templates from a git repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := templates.NewManager(".oursky")
			repo := templates.TemplateRepo{
				URL: args[0],
			}

			// Check if branch flag is provided
			if branch, _ := cmd.Flags().GetString("branch"); branch != "" {
				repo.Branch = branch
			}

			return manager.PullRepo(repo)
		},
	}
	templatesPullCmd.Flags().String("branch", "", "Branch to pull from")

	templatesListCmd := &cobra.Command{
		Use:   "list",
		Short: "List pulled template repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := templates.NewManager(".oursky")
			repos, err := manager.ListRepos()
			if err != nil {
				return err
			}

			if len(repos) == 0 {
				fmt.Println("No template repositories pulled")
				return nil
			}

			fmt.Println("Pulled template repositories:")
			for _, repo := range repos {
				fmt.Printf("  - %s\n", repo)
			}
			return nil
		},
	}

	templatesMergeCmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge templates from all sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := templates.NewManager(".oursky")
			data, err := manager.Merge(".oursky")
			if err != nil {
				return err
			}

			fmt.Printf("Merged %d skills, %d rules, %d commands\n",
				len(data.Skills), len(data.Rules), len(data.Commands))
			return nil
		},
	}

	templatesCmd.AddCommand(templatesPullCmd, templatesListCmd, templatesMergeCmd)

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update oursky to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateOursky()
		},
	}

	rootCmd.AddCommand(initCmd, devCmd, listCmd, killCmd, agentCmd, templatesCmd, updateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
