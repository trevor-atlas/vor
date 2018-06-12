package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trevor-atlas/vor/src/git"
	"github.com/trevor-atlas/vor/src/utils"
)

func stashExistingChanges() {
	cmdOutput, _ := git.Call("status")
	if utils.CaseInsensitiveContains(cmdOutput, "deleted") || utils.CaseInsensitiveContains(cmdOutput, "modified") {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("working directory is not clean. Stash changes? (Y/N)")
		text, _ := reader.ReadString('\n')
		if utils.CaseInsensitiveContains(text, "N") {
			utils.ExitWithMessage("")
		}
		git.Call("stash")
	}
}

func createBranch(args []string) {
	// branchTemplate := "{jira-issue-number}/{jira-issue-type}/{jira-issue-title}"
	jiraKey := viper.Get("VOR_JIRA_API_KEY")
	fmt.Println(jiraKey)

}

// steps for branch:
// 1. check for existing local git
// 2. stash changes if they exist
// 3. get JIRA info for ticket, create branch
var branch = &cobra.Command{
	Use:   "branch",
	Short: "create a branch for a jira issue",
	Long:  `creates a git branch for a given jira issue with the default template of {jira-issue-number}/{jira-issue-type}/{jira-issue-title}`,
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureAvailability()
		stashExistingChanges()
		createBranch(args)
	},
}

func init() {
	rootCmd.AddCommand(branch)
}
