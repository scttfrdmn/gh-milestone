package cmd

import (
	"fmt"
	"strconv"

	"github.com/scttfrdmn/gh-milestone/pkg/api"
	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:   "close <milestone-number|title>",
	Short: "Close a milestone",
	Long: `Close a milestone (sets state to closed).

Examples:
  gh milestone close 1
  gh milestone close "v1.0.0"`,
	Args: cobra.ExactArgs(1),
	RunE: runClose,
}

func runClose(cmd *cobra.Command, args []string) error {
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
		State: "closed",
	}

	milestone, err := client.UpdateMilestone(milestoneNumber, input)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Closed milestone %q #%d\n", milestone.Title, milestone.Number)

	return nil
}
