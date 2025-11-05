package cmd

import (
	"github.com/spf13/cobra"
)

var (
	repoFlag string
)

var rootCmd = &cobra.Command{
	Use:   "milestone",
	Short: "Manage GitHub milestones",
	Long: `A comprehensive milestone management tool for GitHub CLI.

Examples:
  gh milestone list
  gh milestone create "v1.0.0" --title "First Release" --due-date "2025-12-31"
  gh milestone view 1
  gh milestone edit 1 --due-date "2026-01-15"
  gh milestone close 1`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&repoFlag, "repo", "R", "", "Repository (owner/repo)")

	// Add subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(closeCmd)
	rootCmd.AddCommand(reopenCmd)
}
