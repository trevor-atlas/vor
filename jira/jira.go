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

func PrintIssue(issue JiraIssue) {
	var b strings.Builder
	w := b.WriteString
	divider := "\n--------------------------------\n"
	pad := utils.PadOutput(4)

	w(issue.Fields.IssueType.Name + "\n")
	w(divider + "Title: " +issue.Fields.Summary+"\n")
	w("Description: "+formatMultiline(issue.Fields.Description, pad) + "\n")
	w("created on: "+ issue.Fields.Created + "\n")

	if len(issue.Fields.Comment.Comments) > 0 {
		w("Comments:" + divider)
		for _, comment := range issue.Fields.Comment.Comments {
			w(pad("Author: " + comment.Author.Name + "\n"))
			w(formatMultiline("\"" +comment.Body+ "\"", pad))
			w(divider)
		}
	}
	fmt.Println(b.String())
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
		CheckRedirect: redirectHandler,
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
