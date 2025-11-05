package cmd

import (
	"fmt"
	"strconv"

	"github.com/scttfrdmn/gh-milestone/pkg/api"
	"github.com/spf13/cobra"
)

var reopenCmd = &cobra.Command{
	Use:   "reopen <milestone-number|title>",
	Short: "Reopen a milestone",
	Long: `Reopen a milestone (sets state to open).

Examples:
  gh milestone reopen 1
  gh milestone reopen "v1.0.0"`,
	Args: cobra.ExactArgs(1),
	RunE: runReopen,
}

func runReopen(cmd *cobra.Command, args []string) error {
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

	input := api.MilestoneInput{
		State: "open",
	}

	milestone, err := client.UpdateMilestone(milestoneNumber, input)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Reopened milestone %q #%d\n", milestone.Title, milestone.Number)

	return nil
}
