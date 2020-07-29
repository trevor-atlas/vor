package jira

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dustin/go-humanize"

	"trevoratlas.com/vor/formatters"

	"unicode/utf8"

	"trevoratlas.com/vor/utils"

	"log"
)

const (
	top_left     string = "\u256D"
	top_right    string = "\u256E"
	bottom_left  string = "\u2570"
	bottom_right string = "\u256F"
	x_line       string = "\u2500"
	y_line       string = "\u2502"
	shade        string = "\u2591"
	left_quote   string = "\u201C"
	right_quote  string = "\u201D"
	max_len      int    = 70
)

func formatMultiline(message string) string {
	r := strings.NewReplacer(
		"\\{", "{",
		"\\}", "}",
		"\\-", "-",
		"\\!", "!")
	message = r.Replace(message)

	inlineCode := regexp.MustCompile(`{{([a-zA-Z.\-_$@#%^&*()=+\s~/\\]*)}}$`)

	split := strings.SplitAfter(message, "}}")

	for i, str := range split {
		if utils.Contains(str, "}}") || utils.Contains(str, "{{") {
			s := inlineCode.ReplaceAllString(str, formatCode(`$1`))
			split[i] = s
		}
	}

	message = strings.Join(split, "")

	split = strings.SplitAfter(message, "{code}")

	for i, str := range split {
		if utils.Contains(str, "{code}") {
			split[i] = formatCode(str)
		}
	}

	message = strings.Join(split, "")

	return message
}

func PrintIssueJson(issue JiraIssue) {
	data, err := json.MarshalIndent(issue, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}

func PrintIssuesJson(issues JiraIssues) {
	data, err := json.MarshalIndent(issues, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}

type GroupedIssues map[string][]JiraIssue

type FieldType int

const (
	DateField FieldType = iota
	PersonField
	URLField
	GenericField
)

func PrintKeyValue(key string, value string, t FieldType) {
	key = formatters.CYAN(strings.ToLower(key)) + ": "
	switch t {
	case DateField:
		fmt.Print(key + formatters.YELLOW(value) + "\n")
	case PersonField:
		fmt.Print(key + formatters.MAGENTA(value) + "\n")
	case URLField:
		fmt.Print(key + formatters.BLUE(value) + "\n")
	case GenericField:
		fmt.Print(key + value + "\n")
	}
}

func formatCode(s string) string {
	return "\033[30;5;47m" + s + "\033[0m"
}

func PrintIssues(issues JiraIssues) {
	divider := "\n--------------------------------\n"
	issueURL := "" + utils.JIRA_ORGNAME + ".atlassian.net/browse/"
	grouped := make(GroupedIssues)

	for _, issue := range issues.Issues {
		if issue.Fields.IssueType.Name == "Sub-task" {
			continue
		}
		if issue.Fields.IssueType.Name == "Epic" {
			continue
		}
		if issue.Fields.Status.StatusCategory.Name == "Done" {
			continue
		}
		grouped[issue.Fields.Status.Name] = append(grouped[issue.Fields.Status.Name], issue)
	}

	for _, value := range grouped {
		if value == nil || len(value) < 1 {
			continue
		}
		for _, issue := range value {
			created := time.Time(*issue.Fields.Created).Format("2006-01-02 15:04") + " (" + humanize.Time(time.Time(*issue.Fields.Created)) + ")"
			updated := humanize.Time(time.Time(*issue.Fields.Updated))

			fmt.Println()
			fmt.Println(issue.Key + "/" + issue.Fields.IssueType.Name + "/" + issue.Fields.Summary)
			PrintKeyValue("assignee", issue.Fields.Assignee.DisplayName, PersonField)
			PrintKeyValue("reporter", issue.Fields.Reporter.DisplayName, PersonField)
			PrintKeyValue("created", created, DateField)
			PrintKeyValue("updated", updated, DateField)
			PrintKeyValue("status", issue.Fields.Status.StatusCategory.Name, GenericField)
			PrintKeyValue("URL", issueURL+issue.Key, URLField)
			fmt.Println(divider)
		}
	}
}

// BuildTitle creates a nice top-level outlined title
// FIXME: don't pass colorized output to this method yet, it messes up the math and I'm not sure of a good way to account for it yet
func BuildTitle(title string, maxPadding int) (formattedTitle string, length int) {
	r := strings.Repeat
	if maxPadding > max_len {
		maxPadding = max_len
	}
	if utf8.RuneCountInString(title) > max_len {
		title = title[0:max_len-3] + "..."
	}
	titleLen := utf8.RuneCountInString(title)

	padAmount := 0
	if titleLen < max_len {
		padAmount = max_len - titleLen - 2
		if padAmount > maxPadding {
			padAmount = maxPadding
		}
	}
	padding := " " + r(shade, padAmount) + " "
	if utf8.RuneCountInString(padding) < 3 {
		padding = " "
	}
	paddingSize := utf8.RuneCountInString(padding)
	// the 2 here is accounting for the added spaces
	hBorder := r(x_line, titleLen+(paddingSize*2))
	result := top_left + hBorder + top_right + "\n" +
		y_line + padding + title + padding + y_line + "\n" +
		bottom_left + hBorder + bottom_right + "\n"
	return result, utf8.RuneCountInString(padding + title + padding)
}

func getIssueURL(issueID string) string {
	return "https://" + utils.JIRA_ORGNAME + ".atlassian.net/browse/" + issueID
}

func PrintIssue(issue JiraIssue) {
	fmt.Println("made it here")
	title, _ := BuildTitle(left_quote+issue.Fields.Summary+right_quote, 10)
	assignee := issue.Fields.Assignee.DisplayName
	if assignee == "" {
		assignee = "<unassigned>"
	}

	created := time.Time(*issue.Fields.Created).Format("2006-01-02 15:04") + " (" + humanize.Time(time.Time(*issue.Fields.Created)) + ")"
	updated := humanize.Time(time.Time(*issue.Fields.Updated))

	fmt.Println(title)
	PrintKeyValue("issue", issue.Key, GenericField)
	PrintKeyValue("type", issue.Fields.IssueType.Name, GenericField)
	PrintKeyValue("status", issue.Fields.Status.Name, GenericField)
	PrintKeyValue("reporter", issue.Fields.Reporter.DisplayName, PersonField)
	PrintKeyValue("assignee", assignee, PersonField)
	PrintKeyValue("created", created, DateField)
	PrintKeyValue("updated", updated, DateField)
	PrintKeyValue("url", getIssueURL(issue.Key), URLField)

	if issue.Fields.Description != "" {
		PrintKeyValue("description", formatMultiline(issue.Fields.Description), GenericField)
	} else {
		PrintKeyValue("description", "none", GenericField)
	}

	if len(issue.Fields.Comment.Comments) > 0 {
		nestedPad := utils.PadOutput(4)
		fmt.Println(formatters.CYAN("comments:"))
		for _, comment := range issue.Fields.Comment.Comments {
			commentCreated := time.Time(*comment.Created).Format("2006-01-02 15:04") + " (" + humanize.Time(time.Time(*comment.Created)) + ")\n"
			commentUpdated := humanize.Time(time.Time(*comment.Updated))
			PrintKeyValue(nestedPad("author"), comment.Author.Name, PersonField)
			PrintKeyValue(nestedPad("created"), commentCreated, DateField)
			PrintKeyValue(nestedPad("updated"), commentUpdated, DateField)
			PrintKeyValue(nestedPad("body"), comment.Body, GenericField)
		}
	}
}

func GetIssues(get func(url string) ([]byte, error)) JiraIssues {
	url := "https://" + utils.JIRA_ORGNAME + ".atlassian." +
		"net/rest/api/2/search?jql=assignee=currentuser()+order+by+status+asc&expand=fields"

	body, err := get(url)
	if err != nil {
		utils.Exit("There was a problem making the request to the jira API in `GetIssues`")
	}

	parsed := JiraIssues{}
	parseError := json.Unmarshal(body, &parsed)
	if parseError != nil {
		fmt.Printf("There was a problem parsing the jira API response:\n%s\n", parseError)
		utils.Exit("")
	}
	return parsed
}

func GetIssue(issueNumber string, get func(url string) ([]byte, error)) JiraIssue {
	url := "https://" + utils.JIRA_ORGNAME + ".atlassian." +
		"net/rest/api/2/issue/" + issueNumber + "?expand=fields"

	res, err := get(url)
	if err != nil {
		fmt.Printf("error making request")
		panic(err)
	}

	parsed := JiraIssue{}
	parseError := json.Unmarshal(res, &parsed)
	if parseError != nil {
		fmt.Printf("error parsing json\n %s", parseError)
		panic(parseError)
	}
	return parsed
}
