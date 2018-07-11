package jira

import (
	"github.com/fatih/color"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/trevor-atlas/vor/utils"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectHandler(req *http.Request, via []*http.Request) error {
	jiraUsername := utils.GetStringEnv("jira.username")
	jiraKey := utils.GetStringEnv("jira.apikey")
	req.Header.Add("Authorization", "Basic "+basicAuth(jiraUsername, jiraKey))
	return nil
}

func formatMultiline(message string, formatter func(string) string) string {
	maxLen := 120
	var b strings.Builder
	for _, line := range strings.Split(message, "\n") {
		if len(line) > maxLen && strings.Contains(line, ". ") {
			for _, str := range strings.SplitAfter(line, ". ") {
				b.WriteString(formatter(str) + "\n")
			}
		} else if len(line) > maxLen {
			b.WriteString(formatter(line[0:maxLen]) + "\n")
			b.WriteString(formatter(line[maxLen:]) + "\n")
		} else {
			b.WriteString(formatter(line) + "\n")
		}
	}
	return strings.Trim(b.String(), "\n")
}

func PrintIssues(issues JiraIssues) {
	orgName := utils.GetStringEnv("jira.orgname")
	var b strings.Builder
	w := b.WriteString
	divider := "\n--------------------------------\n"
	pad := utils.PadOutput(2)
	cyan := color.New(color.FgHiCyan).SprintFunc()
	// red := color.New(color.FgHiRed).SprintFunc()

	issueURL := "" + orgName + ".atlassian.net/browse/"

	var todo []JiraIssue
	var inProgress []JiraIssue
	var done []JiraIssue

	for _, issue := range issues.Issues {
		if issue.Fields.Status.StatusCategory.Name == "To Do" {
			todo = append(todo, issue)
		}
		if issue.Fields.Status.StatusCategory.Name == "In Progress" {
			inProgress = append(inProgress, issue)
		}
		if issue.Fields.Status.StatusCategory.Name == "Done" {
			done = append(done, issue)
		}
	}

	sortedIssues := append(todo, inProgress...)
	sortedIssues = append(sortedIssues, done...)

	var issueColumn string
	var alreadyPrintedColumn bool

	for _, issue := range sortedIssues{
		// issueLabel := issue.Fields.IssueType.Name
		if issueColumn == "" || issueColumn != issue.Fields.Status.StatusCategory.Name {
			issueColumn = issue.Fields.Status.StatusCategory.Name
			alreadyPrintedColumn = false
		}

		if issueColumn == issue.Fields.Status.StatusCategory.Name && !alreadyPrintedColumn {
			w("\n")
			w(cyan(issueColumn))
			w(divider)
			alreadyPrintedColumn = true
		}
		w(cyan(issue.Key) + " " + issueURL + issue.Key + "\n")
		w(pad(issue.Fields.Summary))
		w("\n\n")
	}

	fmt.Println(b.String())
}

func PrintIssue(issue JiraIssue) {
	var b strings.Builder
	w := b.WriteString
	divider := "\n--------------------------------\n"
	pad := utils.PadOutput(4)

	w(issue.Fields.IssueType.Name + "\n")
	w(divider + "Title: " + issue.Fields.Summary + "\n")
	w("Description: " + formatMultiline(issue.Fields.Description, pad) + "\n")
	w("created on: " + issue.Fields.Created + "\n")

	if len(issue.Fields.Comment.Comments) > 0 {
		w("Comments:" + divider)
		for _, comment := range issue.Fields.Comment.Comments {
			w(pad("Author: " + comment.Author.Name + "\n"))
			w(formatMultiline("\""+comment.Body+"\"", pad))
			w(divider)
		}
	}
	fmt.Println(b.String())
}

func Get(url string) (*http.Response, error) {
	jiraUsername := utils.GetStringEnv("jira.username")
	jiraKey := utils.GetStringEnv("jira.apikey")

	client := &http.Client{
		CheckRedirect: redirectHandler,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(jiraUsername, jiraKey))

	resp, err := client.Do(req)
	return resp, err
}

func GetIssues() JiraIssues {
	orgName := utils.GetStringEnv("jira.orgname")
	url := "https://" + orgName + ".atlassian.net/rest/api/2/search?jql=assignee=currentuser()&expand=fields"

	resp, err := Get(url)
	if err != nil {
		fmt.Printf("error making request")
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	parsed := JiraIssues{}

	parseError := json.Unmarshal(body, &parsed)
	if parseError != nil {
		fmt.Printf("error parsing json\n %s", parseError)
		panic(parseError)
	}
	return parsed
}

func GetIssue(issueNumber string) JiraIssue {
	orgName := utils.GetStringEnv("jira.orgname")
	url := "https://" + orgName + ".atlassian.net/rest/api/2/issue/" + issueNumber + "?expand=fields"

	resp, err := Get(url)
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
