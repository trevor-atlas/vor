package commands

import (
	"github.com/spf13/cobra"
	"trevoratlas.com/vor/jira"
	"trevoratlas.com/vor/utils"
)

var issueJSON bool
var issue = &cobra.Command{
	Use:   "issue",
	Short: "get metadata for a specific issue",
	Long: `
	prints out an issue and its comments
	`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckArgs(args, "an issue number XX-1234")
		service := jira.Service{}
		issue := service.GetIssue(args[0])

		// if issueJSON {
		// 	jira.PrintIssueJson(issue)
		// } else {
			jira.PrintIssue(issue)

		// }
	},
}

func init() {
	issue.Flags().BoolVarP(&issueJSON, "json", "j", false, "output issue as json")
	rootCmd.AddCommand(issue)
}
