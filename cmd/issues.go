package commands

import (
	"github.com/spf13/cobra"
	"net/http"
	"time"
	"trevoratlas.com/vor/jira"
	"trevoratlas.com/vor/rest"
)

var issuesJSON bool
var issues = &cobra.Command{
	Use:   "issues",
	Short: "list your jira issues",
	Long: `
	List each of your assigned issues in jira
	`,
	Run: func(cmd *cobra.Command, args []string) {
		issues := jira.Service.GetIssues()

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
