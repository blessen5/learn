package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"learn/internal/config"
	"learn/internal/git"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push notes to remote",
	Long:  "Push committed notes to the configured git remote.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if !git.IsRepo(cfg.Repo.Root) {
			return fmt.Errorf("not a git repository: %s", cfg.Repo.Root)
		}

		gitCmd := exec.Command("git", "push")
		gitCmd.Dir = cfg.Repo.Root
		gitCmd.Stdin = os.Stdin
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		return gitCmd.Run()
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull notes from remote",
	Long:  "Pull latest notes from the configured git remote.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if !git.IsRepo(cfg.Repo.Root) {
			return fmt.Errorf("not a git repository: %s", cfg.Repo.Root)
		}

		gitCmd := exec.Command("git", "pull")
		gitCmd.Dir = cfg.Repo.Root
		gitCmd.Stdin = os.Stdin
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		return gitCmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
}
