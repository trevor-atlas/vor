package commands

import (
	"github.com/spf13/cobra"
	"github.com/trevor-atlas/vor/jira"
	"github.com/trevor-atlas/vor/rest"
	"net/http"
	"time"
)

var issueJSON bool
var issue = &cobra.Command{
	Use:   "issue",
	Short: "get metadata for a specific issue",
	Long: `
	prints out an issue and its comments
	`,
	Run: func(cmd *cobra.Command, args []string) {
		get := jira.InstantiateHttpMethods(rest.NewHTTPClient(
			&http.Client{
				Transport:     nil,
				CheckRedirect: jira.RedirectHandler,
				Jar:           nil,
				Timeout:       time.Second * 10,
			}))
		issue := jira.GetIssue(args[0], get)
		if issueJSON {
			jira.PrintIssueJson(issue)
		} else {
			jira.PrintIssue(issue)
		}
	},
}

func init() {
	issue.Flags().BoolVarP(&issueJSON, "json", "j", false, "output issue as json")
	rootCmd.AddCommand(issue)
}
