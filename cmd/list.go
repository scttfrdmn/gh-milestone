package cmd

import (
	"fmt"

	"github.com/scttfrdmn/gh-milestone-manager/pkg/api"
	"github.com/scttfrdmn/gh-milestone-manager/pkg/format"
	"github.com/spf13/cobra"
)

var (
	listState string
	listSort  string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List milestones",
	Long: `List milestones in a repository.

Examples:
  gh milestone list
  gh milestone list --state closed
  gh milestone list --state all --sort due-on`,
	RunE: runList,
}

func init() {
	listCmd.Flags().StringVarP(&listState, "state", "s", "open", "Filter by state (open, closed, all)")
	listCmd.Flags().String("sort", "due_on", "Sort by (due_on, completeness)")
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(repoFlag)
	if err != nil {
		return err
	}

	milestones, err := client.ListMilestones(listState, listSort)
	if err != nil {
		return err
	}

	output := format.FormatMilestoneTable(milestones)
	fmt.Println(output)

	return nil
}
