package commands

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
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

type jiraUser struct {
	Active       bool   `json:"active"`
	TimeZone     string `json:"timeZone"`
	DisplayName  string `json:"displayName"`
	Name         string
	EmailAddress string
	AvatarUrls   map[string]interface{} `json:"-"`
	AccountId    string
	Key          string
	Self         string
}

type jiraComment struct {
	ID           string
	Self         string
	Author       jiraUser
	Body         string
	UpdateAuthor jiraUser
	Created      string
	Updated      string
	Total        int
}

type JiraIssue struct {
	ID     string `json:"id"`
	Self   string `json:"self"` // url to request this issue
	Key    string `json:"key"`  // AQ-XXXX
	Fields struct {
		Summary           string // title of jira issue
		Created           string `json:"created"` // 2018-05-25T04:18:06.836-0500
		Updated           string `json:"updated"` // 2018-06-11T22:23:03.606-0500
		Description       string // description of Jira issue
		Customfield_12022 struct {
			Value string // team name
		}
		Reporter jiraUser
		Assignee jiraUser
		Comment  struct {
			Comments []jiraComment
		}
		Priority struct {
			Name string `json:"priority"` // Medium
		}
		IssueType struct {
			Name    string `json:"name"` // Bug, Task, Story
			Subtask bool   `json:"subtask"`
			IconURL string `json:"iconUrl"`
		} `json:"issuetype"`
	} `json:"fields"`
}

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
		fmt.Printf("error making request")
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	parsed := JiraIssue{}

	parseError := json.Unmarshal(body, &parsed)
	if parseError != nil {
		fmt.Printf("error parsing json\n %s", parseError)
		panic(parseError)
	}
	fmt.Println(string(body))
	fmt.Println(
		"\n" +
			parsed.Fields.IssueType.Name + "\n" +
			parsed.Fields.Summary + "\n" +
			parsed.Fields.Description)
	fmt.Println()
	fmt.Println("comments:\n========================================================")
	for i := range parsed.Fields.Comment.Comments {
		comment := parsed.Fields.Comment.Comments[i]
		fmt.Println(comment.Author.Name)
		fmt.Println(comment.Body)
		fmt.Println(`--------------------------------------------------`)
		fmt.Println()
	}
}

func buildJiraURL(args []string) string {
	// projectName := viper.GetString("VOR_PROJECT_NAME")
	urlParts := []string{"https://", "aquicore", ".atlassian.net/rest/api/2/issue/", args[0], "?expand=fields"}
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
		// stashExistingChanges()
		createBranch(args)
	},
}

func init() {
	rootCmd.AddCommand(branch)
}
