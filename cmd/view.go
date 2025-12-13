package cmd

import (
	"fmt"
	"strconv"

	"github.com/scttfrdmn/gh-milestone-manager/pkg/api"
	"github.com/scttfrdmn/gh-milestone-manager/pkg/format"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view <milestone-number|title>",
	Short: "View milestone details",
	Long: `View detailed information about a milestone.

Examples:
  gh milestone view 1
  gh milestone view "v1.0.0"`,
	Args: cobra.ExactArgs(1),
	RunE: runView,
}

func runView(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(repoFlag)
	if err != nil {
		return err
	}

	// Try to parse as number first
	var milestone *api.Milestone
	if number, err := strconv.Atoi(args[0]); err == nil {
		milestone, err = client.GetMilestone(number)
		if err != nil {
			return err
		}
	} else {
		// Try to find by title
		milestone, err = client.FindMilestoneByTitle(args[0])
		if err != nil {
			return err
		}
	}

	output := format.FormatMilestoneDetail(milestone)
	fmt.Println(output)

	return nil
}
