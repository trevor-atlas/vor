package commands

import (
	"github.com/spf13/cobra"
	"github.com/trevor-atlas/vor/jira"
	"github.com/trevor-atlas/vor/rest"
	"net/http"
	"time"
)

var issuesJSON bool
var issues = &cobra.Command{
	Use:   "issues",
	Short: "list your jira issues",
	Long: `
	List each of your assigned issues in jira
	`,
	Run: func(cmd *cobra.Command, args []string) {
		get := jira.InstantiateHttpMethods(rest.NewHTTPClient(
			&http.Client{
				Transport:     nil,
				CheckRedirect: jira.RedirectHandler,
				Jar:           nil,
				Timeout:       time.Second * 10,
			}))
		issues := jira.GetIssues(get)
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
