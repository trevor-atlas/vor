package commands

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"trevoratlas.com/vor/git"
	"trevoratlas.com/vor/jira"

	"trevoratlas.com/vor/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

// TODO: this should be generic per issue provider
func generateIssueTag(issue jira.JiraIssue) string {
	issueType := utils.LowerKebabCase(issue.Fields.IssueType.Name)
	priority := utils.LowerKebabCase(issue.Fields.Priority.Name)

	switch issueType {
	case "bug":
		if priority == "blocker" {
			return "blocker"
		}
		return "bug"
	case "story", "task":
		return "feature"
	default:
		return issueType
	}
}

func generateBranchName(issue jira.JiraIssue) string {
	branchTemplate := viper.GetString("git.branchtemplate")
	projectName := viper.GetString("projectname")
	author := viper.GetString("author")
	template := strings.Split(branchTemplate, "/")

	// TODO add output to logger if an entry is empty
	// this is a user error but still confusing
	for i := range template {
		switch template[i] {
		case "{projectname}":
			template[i] = utils.LowerKebabCase(projectName)
			break
		case "{jira-issue-number}":
			template[i] = issue.Key
			break
		case "{jira-issue-type}":
			template[i] = generateIssueTag(issue)
			break
		case "{jira-issue-title}":
			template[i] = utils.LowerKebabCase(issue.Fields.Summary)
			break
		case "{date}":
			template[i] = time.Now().Format("06-01-02")
			break
		case "{author}":
			template[i] = utils.LowerKebabCase(author)
			break
		}
	}

	// remove empty entries
	var cleanTemplate []string
	for _, item := range template {
		if item != "" {
			cleanTemplate = append(cleanTemplate, item)
		}
	}

	branchName := strings.Join(cleanTemplate, "/")
	utils.Debug("build branch name: " + branchName)
	return branchName
}

func createBranch(args []string) (branchName string) {
	gc := git.New()
	utils.Debug("cli args: ", args)
	jiraService := jira.Service{}
	issue := jiraService.GetIssue(args[0])
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
				utils.Exit("error calling local git")
			}
			fmt.Println("checked out existing local branch: '" + cyan(branch) + "'")
			return
		}
	}
	_, err := gc.Call("checkout -b " + newBranchName)
	if err != nil {
		utils.Exit("error calling local git")
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
		if len(args) == 0 {
			utils.Exit("Must provide an issue key `vor branch XX-1234`")
		}
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
			shouldStash := utils.Confirm("Working directory is not clean. Stash changes?")
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
		affirm := utils.Confirm(branch + " created.\nwould you like to re-apply your stashed changes?")
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
