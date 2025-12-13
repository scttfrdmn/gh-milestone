package cmd

import (
	"fmt"

	"github.com/scttfrdmn/gh-milestone-manager/pkg/api"
	"github.com/spf13/cobra"
)

var (
	createTitle       string
	createDescription string
	createDueDate     string
	createState       string
)

var createCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new milestone",
	Long: `Create a new milestone in a repository.

Examples:
  gh milestone create "v1.0.0" --title "First Release" --due-date "2025-12-31"
  gh milestone create "Sprint 5" --description "Q1 Sprint" --due-date "2025-03-01"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCreate,
}

func init() {
	createCmd.Flags().StringVarP(&createTitle, "title", "t", "", "Milestone title (required)")
	createCmd.Flags().StringVarP(&createDescription, "description", "d", "", "Milestone description")
	createCmd.Flags().String("due-date", "", "Due date (YYYY-MM-DD)")
	createCmd.Flags().String("state", "open", "State (open or closed)")
	createCmd.MarkFlagRequired("title")
}

func runCreate(cmd *cobra.Command, args []string) error {
	// If positional arg provided, use it as title
	if len(args) > 0 {
		createTitle = args[0]
	}

	if createTitle == "" {
		return fmt.Errorf("title is required")
	}

	client, err := api.NewClient(repoFlag)
	if err != nil {
		return err
	}

	// Parse due date if provided
	dueDate, _ := cmd.Flags().GetString("due-date")
	parsedDueDate := ""
	if dueDate != "" {
		parsedDueDate, err = api.ParseDueDate(dueDate)
		if err != nil {
			return err
		}
	}

	state, _ := cmd.Flags().GetString("state")

	input := api.MilestoneInput{
		Title:       createTitle,
		Description: createDescription,
		DueOn:       parsedDueDate,
		State:       state,
	}

	milestone, err := client.CreateMilestone(input)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Created milestone %q #%d\n", milestone.Title, milestone.Number)
	fmt.Printf("  %s\n", milestone.URL)

	return nil
}
