package cmd

import (
	"fmt"
	"strconv"

	"github.com/scttfrdmn/gh-milestone-manager/pkg/api"
	"github.com/spf13/cobra"
)

var (
	editTitle       string
	editDescription string
	editDueDate     string
	editState       string
)

var editCmd = &cobra.Command{
	Use:   "edit <milestone-number|title>",
	Short: "Edit a milestone",
	Long: `Edit an existing milestone.

Examples:
  gh milestone edit 1 --title "v1.0.0 - Updated"
  gh milestone edit 1 --due-date "2026-01-15"
  gh milestone edit "v1.0.0" --state closed`,
	Args: cobra.ExactArgs(1),
	RunE: runEdit,
}

func init() {
	editCmd.Flags().StringVarP(&editTitle, "title", "t", "", "New milestone title")
	editCmd.Flags().StringVarP(&editDescription, "description", "d", "", "New milestone description")
	editCmd.Flags().String("due-date", "", "New due date (YYYY-MM-DD)")
	editCmd.Flags().String("state", "", "New state (open or closed)")
}

func runEdit(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(repoFlag)
	if err != nil {
		return err
	}

	// Find the milestone
	var milestoneNumber int
	if number, err := strconv.Atoi(args[0]); err == nil {
		milestoneNumber = number
	} else {
		milestone, err := client.FindMilestoneByTitle(args[0])
		if err != nil {
			return err
		}
		milestoneNumber = milestone.Number
	}

	// Build input from flags that were set
	input := api.MilestoneInput{}

	if cmd.Flags().Changed("title") {
		input.Title = editTitle
	}

	if cmd.Flags().Changed("description") {
		input.Description = editDescription
	}

	if cmd.Flags().Changed("due-date") {
		dueDate, _ := cmd.Flags().GetString("due-date")
		parsedDueDate, err := api.ParseDueDate(dueDate)
		if err != nil {
			return err
		}
		input.DueOn = parsedDueDate
	}

	if cmd.Flags().Changed("state") {
		state, _ := cmd.Flags().GetString("state")
		input.State = state
	}

	milestone, err := client.UpdateMilestone(milestoneNumber, input)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Updated milestone %q #%d\n", milestone.Title, milestone.Number)
	fmt.Printf("  %s\n", milestone.URL)

	return nil
}
