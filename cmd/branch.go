package commands

import (
	"github.com/fatih/color"
	"fmt"
	"strings"

	"github.com/trevor-atlas/vor/jira"
	"github.com/trevor-atlas/vor/logger"

	"github.com/spf13/cobra"
	"github.com/trevor-atlas/vor/git"
	"github.com/trevor-atlas/vor/utils"
)

func generateIssueTag(issue jira.JiraIssue) string {
	switch issue.Fields.IssueType.Name {
	case "bug":
		if issue.Fields.Priority.Name == "blocker" {
			return "break"
		}
	return "bug"
	case "story", "task": return "feature"
	default: return "feature"
	}
}

func generateBranchName(issue jira.JiraIssue) string {
	branchTemplate := utils.GetStringEnv("branchtemplate")
	projectName := utils.GetStringEnv("projectname")
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

func createBranch(args []string) {
	logger.Debug("cli args: ", args)
	issue := jira.GetIssue(args[0])
	newBranchName := generateBranchName(issue)
	fmt.Println(newBranchName)
	localBranches, _ := git.Call("branch")
	cyan := color.New(color.FgHiCyan).SprintFunc()
	replacer := strings.NewReplacer(
		" ", "",
		"\r", "",
		"\n", "")
	for _, branch := range strings.Split(localBranches, "\n") {
		if strings.ToLower(replacer.Replace(branch)) == strings.ToLower(newBranchName) {
			git.Call("checkout " + branch)
			fmt.Println("checked out existing local branch: '" + cyan(branch) + "'")
			return
		}
	}
	git.Call("checkout -b " + newBranchName)
	fmt.Println("checked out new local branch: '" + cyan(newBranchName) + "'")
}

var branch = &cobra.Command{
	Use:   "branch",
	Short: "create a branch for a jira issue",
	Long: `
	creates a git branch for a given jira issue
	NOTE: (you must include the project key E.G. XX-1234)
	vor branch XX-4321
	`,
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureAvailability()
		git.StashExistingChanges()
		createBranch(args)
	},
}

func init() {
	rootCmd.AddCommand(branch)
}
