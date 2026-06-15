package cmd

import (
	"fmt"
	"path/filepath"
	"sort"

	"learn/internal/config"
	"learn/internal/file"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes by category",
	Long:  "Display a tree view of all notes organized by category.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		categories := file.ListCategories(cfg.Repo.Root)
		if len(categories) == 0 {
			fmt.Println("No categories found.")
			return nil
		}

		sort.Strings(categories)
		totalNotes := 0

		for _, cat := range categories {
			catDir := filepath.Join(cfg.Repo.Root, cat)
			files, err := file.ListMarkdownFiles(catDir)
			if err != nil || len(files) == 0 {
				continue
			}

			fmt.Printf("\n%s/ (%d)\n", cat, len(files))
			for _, f := range files {
				name := filepath.Base(f)
				fmt.Printf("  %s\n", name)
				totalNotes++
			}
		}

		fmt.Printf("\nTotal: %d notes across %d categories\n", totalNotes, len(categories))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
