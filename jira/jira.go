package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

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
	// var builder strings.Builder
	// b := builder.WriteString
	divider := "\n--------------------------------\n"
	// pad := utils.PadOutput(2)
	cyan := color.New(color.FgHiCyan).SprintFunc()
	// red := color.New(color.FgHiRed).SprintFunc()
	issueURL := "" + orgName + ".atlassian.net/browse/"
	toDo := []JiraIssue{}
	inProg := []JiraIssue{}
	review := []JiraIssue{}
	verify := []JiraIssue{}

	for _, issue := range issues.Issues {
		switch issue.Fields.Status.Name {
		case "To Do":
			toDo = append(toDo, issue)
		case "In Progress":
			inProg = append(inProg, issue)
		case "Review":
			review = append(review, issue)
		case "Verification":
			verify = append(verify, issue)
		default: continue
		}
	}

	columns := [][]JiraIssue{toDo, inProg, review, verify}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 12, 8, 2, '\t', ' ')

	for _, column := range columns {
		if column == nil || len(column) < 1 {
			continue
		}
		fmt.Println()
		fmt.Print(cyan(column[0].Fields.Status.StatusCategory.Name))
		fmt.Println(divider)
		fmt.Fprintln(w, cyan("Issue No.") + "\t " + cyan("Issue Type") + "\t " + cyan("URL"))
		fmt.Fprintln(w)
		for _, issue := range column {
			fmt.Fprintln(w,
				issue.Key + "\t "+
				issue.Fields.IssueType.Name + "\t "+
				issueURL+issue.Key)
			fmt.Fprintln(w)

			// fmt.Println(issue.Fields.Summary)
		}
		w.Flush()
	}
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
	url := "https://" + orgName + ".atlassian.net/rest/api/2/search?jql=assignee=currentuser()+order+by+status+asc&expand=fields"

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

func Filter(vs []JiraIssue, f func(JiraIssue) bool) []JiraIssue {
	vsf := make([]JiraIssue, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
