package commands

import (
	"github.com/trevor-atlas/vor/jira"
	"github.com/spf13/cobra"
)

var issue = &cobra.Command{
	Use:   "issue",
	Short: "get metadata for a specific issue",
	Long: `
	prints out an issue and its comments
	`,
	Run: func(cmd *cobra.Command, args []string) {
		issues := jira.GetIssues()
		for _, issue := range issues.Issues {
			if issue.Key == args[0] {
				jira.PrintIssue(issue)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(issue)
}