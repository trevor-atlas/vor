package commands

import (
	"github.com/spf13/viper"
	"fmt"
	"os"
	"strings"
	"bufio"
	"sync"
	"github.com/spf13/cobra"
	"github.com/trevor-atlas/vor/src/utils"
)

func exitWithMessage(message string) {
	fmt.Println(message + "\ncanceling operation...")
	os.Exit(1)
}

func CaseInsensitiveContains(s, substr string) bool {
    s, substr = strings.ToUpper(s), strings.ToUpper(substr)
    return strings.Contains(s, substr)
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

// returns wether of not git is available
func ensureGitAvailability() {
	localGit, _ := exists("/usr/local/bin/git");
	if localGit {
		_, err := callGit("status")
		if err != nil {
			exitWithMessage("Git Status failed, are you in a valid git repository?")
		}
		return
	}
	exitWithMessage("could not find local git at \"/usr/local/bin/git\"")
}

func callGit(command string) (string, error) {
	fmt.Println("debug: calling git command " + command)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	return utils.ShellExec("/usr/local/bin/git " + command, wg)
}

func stashExistingChanges() {
	cmdOutput, _ := callGit("status")
	if CaseInsensitiveContains(cmdOutput, "deleted") || CaseInsensitiveContains(cmdOutput, "modified") {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("working directory is not clean. Stash changes? (Y/N)")
		text , _ := reader.ReadString('\n')
		if CaseInsensitiveContains(text, "N") {
			exitWithMessage("")
		}
		callGit("stash")
	}
}

func createBranch(args []string) {
	// branchTemplate := "{jira-issue-number}/{jira-issue-type}/{jira-issue-title}"
	jira := viper.Get("jira-api-key")
	fmt.Println(jira)

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
		ensureGitAvailability()
		stashExistingChanges()
		createBranch(args)
	},
}

func init() {
	rootCmd.AddCommand(branch)
}