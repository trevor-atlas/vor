package commands

import (
	"github.com/spf13/cobra"

	"trevoratlas.com/vor/jira"
)

var issuesJSON bool
var issues = &cobra.Command{
	Use:   "issues",
	Short: "list your jira issues",
	Long: `
	List each of your assigned issues in jira
	`,
	Run: func(cmd *cobra.Command, args []string) {
		service := jira.Service{}
		issues := service.GetIssues()

		if issuesJSON {
			jira.PrintIssuesJson(issues)
		} else {
			jira.PrintIssues(issues)
		}
	},
}

func init() {
	issues.Flags().BoolVarP(&issuesJSON, "json", "j", false, "output issues as json")
	rootCmd.AddCommand(issues)
}
