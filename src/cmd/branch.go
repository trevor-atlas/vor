package commands

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("username1", "password123"))
	return nil
}

func createBranch(args []string) {
	jiraKey := viper.GetString("VOR_JIRA_API_KEY")
	branchTemplate := viper.GetString("VOR_BRANCH_TEMPLATE")
	projectName := viper.GetString("VOR_PROJECT_NAME")
	templateParts := strings.Split(branchTemplate, "/")
	for i := range templateParts {
		switch templateParts[i] {
		case "{project-name}":
			templateParts[i] = projectName
			break
		case "{jira-issue-number}":
			templateParts[i] = "AQ-XXXX"
			break
		case "{jira-issue-type}":
			templateParts[i] = "bug"
			break
		case "{jira-issue-title}":
			templateParts[i] = "some-jira-issue-title"
		}
	}
	fmt.Println("build issue template: " + strings.Join(templateParts, "/"))

	url := buildJiraURL(args)
	fmt.Println("built jira url: " + url)

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth("trevor.allen@aquicore.com", jiraKey))

	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
}

func buildJiraURL(args []string) string {
	projectName := viper.GetString("VOR_PROJECT_NAME")
	urlParts := []string{"https://", projectName, ".atlassian.net/rest/api/2/issue/", args[0], "?expand=fields"}
	return strings.Join(urlParts, "")
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
