package commands

import (
	"github.com/spf13/cobra"
	"net/http"
	"time"
	"trevoratlas.com/vor/jira"
	"trevoratlas.com/vor/rest"
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
