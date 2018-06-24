package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/utils"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("username1", "password123"))
	return nil
}

func padOutput(padding int) func(string) string {
	return func(str string) string {
		return utils.LeftPad(str, " ", padding)
	}
}

func PrintIssue(issue JiraIssue) {
	pad := padOutput(4)
	var formatted string
	for _, line := range strings.Split(issue.Fields.Description, "\n") {
		formatted += pad(line) + "\n"
	}
	fmt.Println(
		"//////////////////////////////////////////////////////////////////\n" +
			pad(issue.Fields.IssueType.Name+" – created on: "+issue.Fields.Created+"\n") +
			"\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\" + "\n" +
			pad(issue.Fields.Summary+"\n") +
			formatted)

	if len(issue.Fields.Comment.Comments) > 0 {
		fmt.Println(
			"//////////////////////////////////////////////////////////////////\n" +
				pad("COMMENTS:\n") +
				"\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\" + "\n")
		for _, comment := range issue.Fields.Comment.Comments {
			fmt.Println(pad(comment.Author.Name))
			fmt.Println(pad(comment.Body))
			fmt.Println()
		}
	}
}

func GetJiraIssue(issueNumber string) JiraIssue {
	jiraUsername := utils.GetStringEnv("jira.username")
	jiraKey := utils.GetStringEnv("jira.apikey")
	orgName := utils.GetStringEnv("jira.orgname")
	logger.Debug("jira username: " + jiraUsername)
	logger.Debug("jirakey: " + jiraKey)
	logger.Debug("orgname: " + orgName)
	url := "https://" + orgName + ".atlassian.net/rest/api/2/issue/" + issueNumber + "?expand=fields"
	logger.Debug("built jira url: " + url)

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(jiraUsername, jiraKey))

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
	return parsed
}
