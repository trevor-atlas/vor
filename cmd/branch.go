package commands

import (
	"fmt"
	"github.com/fatih/color"
	"strings"

	"github.com/trevor-atlas/vor/jira"
	"github.com/trevor-atlas/vor/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trevor-atlas/vor/git"
	"github.com/trevor-atlas/vor/rest"
	"github.com/trevor-atlas/vor/system"
	"github.com/trevor-atlas/vor/utils"
	"net/http"
	"time"
)

func generateIssueTag(issue jira.JiraIssue) string {
	switch issue.Fields.IssueType.Name {
	case "bug":
		if issue.Fields.Priority.Name == "blocker" {
			return "blocker"
		}
		return "bug"
	case "story", "task":
		return "feature"
	default:
		return "feature"
	}
}

func generateBranchName(issue jira.JiraIssue) string {
	branchTemplate := viper.GetString("branchtemplate")
	projectName := viper.GetString("projectname")
	templateParts := strings.Split(branchTemplate, "/")

	for i := range templateParts {
		switch templateParts[i] {
		case "{projectname}":
			templateParts[i] = projectName
			break
		case "{jira-issue-number}":
			templateParts[i] = issue.Key
			break
		case "{jira-issue-type}":
			templateParts[i] = generateIssueTag(issue)
			break
		case "{jira-issue-title}":
			templateParts[i] = utils.LowerKebabCase(issue.Fields.Summary)
			break
		}
	}

	branchName := strings.Join(templateParts, "/")
	logger.Debug("build branch name: " + branchName)
	return branchName
}

func createBranch(args []string) (branchName string) {
	gc := git.New()
	logger.Debug("cli args: ", args)
	get := jira.InstantiateHttpMethods(rest.NewHTTPClient(
		&http.Client{
			Transport:     nil,
			CheckRedirect: jira.RedirectHandler,
			Jar:           nil,
			Timeout:       time.Second * 10,
		}))
	issue := jira.GetIssue(args[0], get)
	newBranchName := generateBranchName(issue)
	fmt.Println(newBranchName)
	localBranches, _ := gc.Call("branch")
	cyan := color.New(color.FgHiCyan).SprintFunc()
	replacer := strings.NewReplacer(
		"", "",
		"'", "",
		"\"", "",
		",", "",
		" ", "",
		"\r", "",
		"\n", "")
	for _, branch := range strings.Split(localBranches, "\n") {
		if strings.ToLower(replacer.Replace(branch)) == strings.ToLower(newBranchName) {
			_, err := gc.Call("checkout " + branch)
			if err != nil {
				system.Exit("error calling local git")
			}
			fmt.Println("checked out existing local branch: '" + cyan(branch) + "'")
			return
		}
	}
	_, err := gc.Call("checkout -b " + newBranchName)
	if err != nil {
		system.Exit("error calling local git")
	}
	fmt.Println("checked out new local branch: '" + cyan(newBranchName) + "'")
	return newBranchName
}

var affirmAll bool
var declineAll bool
var branch = &cobra.Command{
	Use:   "branch",
	Short: "create a branch for a jira issue",
	Long: `
	creates a git branch for a given jira issue
	NOTE: (you must include the project key E.G. XX-1234)
	vor branch XX-4321
	`,
	Run: func(cmd *cobra.Command, args []string) {
		makeBranch(args)
	},
}

func makeBranch(args []string) {
	var didStash bool
	gc := git.New()
	cmdOutput, _ := gc.Call("status")
	c := func(substr string) bool { return utils.Contains(cmdOutput, substr) }

	if c("deleted") || c("modified") || c("untracked") {
		if !affirmAll || !declineAll {
			shouldStash := system.Confirm("Working directory is not clean. Stash changes?")
			if shouldStash {
				gc.Stash()
				didStash = true
			}
		} else if affirmAll {
			gc.Stash()
			didStash = true
		} else if declineAll {
			didStash = false
		}
	}
	branch := createBranch(args)
	if didStash && !affirmAll && !declineAll {
		affirm := system.Confirm(branch + " created.\nwould you like to re-apply your stashed changes?")
		if affirm {
			gc.UnStash()
		}
	} else if affirmAll {
		gc.UnStash()
	}
}

func init() {
	branch.Flags().BoolVarP(&affirmAll, "yes", "y", false, "skip prompts and affirm")
	branch.Flags().BoolVarP(&declineAll, "no", "n", false, "skip prompts and decline")
	rootCmd.AddCommand(branch)
}
