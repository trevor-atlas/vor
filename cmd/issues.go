package commands

import (
	"github.com/trevor-atlas/vor/jira"
	"github.com/spf13/cobra"
)

var issues = &cobra.Command{
	Use:   "issues",
	Short: "list your jira issues",
	Long: `
	List each of your assigned issues in jira
	`,
	Run: func(cmd *cobra.Command, args []string) {
		issues := jira.GetIssues()
		jira.PrintIssues(issues)
	},
}

func init() {
	rootCmd.AddCommand(issues)
}