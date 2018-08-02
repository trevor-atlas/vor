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
		issue := jira.GetIssue(args[0])
		jira.PrintIssue(issue)
	},
}

func init() {
	rootCmd.AddCommand(issue)
}
