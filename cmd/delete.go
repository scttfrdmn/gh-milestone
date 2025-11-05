package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scttfrdmn/gh-milestone/pkg/api"
	"github.com/spf13/cobra"
)

var deleteYes bool

var deleteCmd = &cobra.Command{
	Use:   "delete <milestone-number|title>",
	Short: "Delete a milestone",
	Long: `Delete a milestone from a repository.

WARNING: This action cannot be undone.

Examples:
  gh milestone delete 1
  gh milestone delete "v1.0.0"
  gh milestone delete 1 --yes  # Skip confirmation`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation prompt")
}

func runDelete(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(repoFlag)
	if err != nil {
		return err
	}

	// Find the milestone
	var milestone *api.Milestone
	if number, err := strconv.Atoi(args[0]); err == nil {
		milestone, err = client.GetMilestone(number)
		if err != nil {
			return err
		}
	} else {
		milestone, err = client.FindMilestoneByTitle(args[0])
		if err != nil {
			return err
		}
	}

	// Confirm deletion unless --yes flag is set
	if !deleteYes {
		fmt.Printf("? Are you sure you want to delete milestone %q #%d? (y/N) ", milestone.Title, milestone.Number)
		total := milestone.OpenIssues + milestone.ClosedIssues
		if total > 0 {
			fmt.Printf("\n  This will unassign %d issue(s) from the milestone.\n", total)
		}

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	err = client.DeleteMilestone(milestone.Number)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Deleted milestone %q #%d\n", milestone.Title, milestone.Number)

	return nil
}
