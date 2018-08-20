package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/spf13/viper"

	"encoding/base64"

	"github.com/trevor-atlas/vor/rest"
	"github.com/trevor-atlas/vor/system"
	"github.com/trevor-atlas/vor/utils"
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

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectHandler(req *http.Request, via []*http.Request) error {
	jiraUsername := viper.GetString("jira.username")
	jiraKey := viper.GetString("jira.apikey")
	req.Header.Add("Authorization", "Basic "+basicAuth(jiraUsername, jiraKey))
	return nil
}

func formatMultiline(message string, formatter func(string) string) string {
	var b strings.Builder
	// write to the string builder
	w := func(str string) {
		b.WriteString(formatter(str))
	}
	// write to the string builder with a new line
	wnl := func(str string) {
		b.WriteString(formatter(str) + "\n")
	}

	for _, line := range strings.Split(message, "\n") {
		if utf8.RuneCountInString(line) > max_len {
			if strings.Contains(line, "{code}") {
				w("```")
				for _, str := range strings.Split(line, "{code}") {
					str_len := utf8.RuneCountInString(str)
					if str_len > max_len {
						if str_len/2 > max_len {
							wnl(str[0 : str_len/3])
							wnl(" " + str[str_len/3:(str_len/3)*2])
							w(" " + str[(str_len/3)*2:])
						} else {
							wnl(str[0 : str_len/2])
							wnl(" " + str[str_len/2:])
						}
					} else {
						wnl(str)
					}
				}
				w("```")
			} else if strings.Contains(line, ". ") {
				for _, str := range strings.SplitAfter(line, ". ") {
					wnl(str)
				}
			} else {
				wnl(line[0 : max_len-4])
				wnl(" " + line[max_len-4:])
			}
		} else {
			wnl(line)
		}
	}
	return strings.Trim(b.String(), "\n")
}

func PrintIssues(issues JiraIssues) {
	orgName := viper.GetString("jira.orgname")
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
		if issue.Fields.IssueType.Name == "Sub-task" {
			continue
		}
		switch issue.Fields.Status.Name {
		case "To Do":
			toDo = append(toDo, issue)
		case "In Progress":
			inProg = append(inProg, issue)
		case "Review":
			review = append(review, issue)
		case "Verification":
			verify = append(verify, issue)
		default:
			continue
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
		fmt.Fprintln(w, cyan("Issue No.")+"\t "+cyan("Issue Type")+"\t "+cyan("URL"))
		fmt.Fprintln(w)
		for _, issue := range column {
			fmt.Fprintln(w,
				issue.Key+"\t "+
					issue.Fields.IssueType.Name+"\t "+
					issueURL+issue.Key)
			fmt.Fprintln(w)
		}
		w.Flush()
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

func PrintIssue(issue JiraIssue) string {
	orgName := viper.GetString("jira.orgname")
	var b strings.Builder
	w := b.WriteString
	pad := utils.PadOutput(2)
	wpnl := func(str string) {
		w(pad(str) + "\n")
	}
	issueURL := "" + orgName + ".atlassian.net/browse/" + issue.Key
	cyan := color.New(color.FgHiCyan).SprintFunc()
	blue := color.New(color.FgHiBlue).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()
	magenta := color.New(color.FgHiMagenta).SprintFunc()
	title, _ := BuildTitle(left_quote+issue.Fields.Summary+right_quote, 10)
	assignee := issue.Fields.Assignee.Name
	if assignee == "" {
		assignee = "<unassigned>"
	}

	w(title)
	wpnl(cyan("issue: ") + issue.Key)
	wpnl(cyan("type: ") + issue.Fields.IssueType.Name)
	wpnl(cyan("status: ") + issue.Fields.Status.Name)
	wpnl(cyan("reporter: ") + magenta(issue.Fields.Reporter.Name))
	wpnl(cyan("assignee: ") + magenta(assignee))
	wpnl(cyan("created: ") + yellow(time.Time(*issue.Fields.Created).Format("2006-01-02 15:04")+" ("+humanize.Time(time.Time(*issue.Fields.Created))) + ")")
	wpnl(cyan("updated: ") + yellow(humanize.Time(time.Time(*issue.Fields.Updated))))
	wpnl(cyan("url: ") + blue(issueURL))
	wpnl(cyan("description:"))
	w(formatMultiline(issue.Fields.Description, utils.PadOutput(4)) + "\n\n")

	if len(issue.Fields.Comment.Comments) > 0 {
		nestedPad := utils.PadOutput(4)
		wpnl(cyan("comments:"))
		for _, comment := range issue.Fields.Comment.Comments {
			w(nestedPad(cyan("author: ") + comment.Author.Name + "\n"))
			w(nestedPad(cyan("created: ") + yellow(time.Time(*comment.Created).Format("2006-01-02 15:04")+" ("+humanize.Time(time.Time(*comment.Created))+")\n")))
			w(nestedPad(cyan("updated: ")+yellow(humanize.Time(time.Time(*comment.Updated)))) + "\n")
			w(nestedPad(cyan("body:\n")))
			w(formatMultiline(comment.Body, utils.PadOutput(6)) + "\n\n")
		}
	}

	result := b.String()
	fmt.Println(result)
	return result
}

func get(url string) ([]byte, error) {
	username := viper.GetString("jira.username")
	apikey := viper.GetString("jira.apikey")

	client := rest.NewHTTPClient(
		&http.Client{
			Transport:     nil,
			CheckRedirect: redirectHandler,
			Jar:           nil,
			Timeout:       time.Second * 10,
		}).
		WithHeader("Accept", "application/json").
		WithHeader("Authorization", "Basic "+basicAuth(username, apikey)).
		URL(url)

	return client.GET()
}

func GetIssues() JiraIssues {
	defer system.ExecutionTimer(time.Now(), "GetIssues")
	orgname := viper.GetString("jira.orgname")
	if orgname == "" {
		system.Exit("jira.orgname config not found.")
	}
	username := viper.GetString("jira.username")
	if username == "" {
		system.Exit("jira.username config not found.")
	}
	apikey := viper.GetString("jira.apikey")
	if apikey == "" {
		system.Exit("jira.apikey config not found.")
	}
	url := "https://" + orgname + ".atlassian.net/rest/api/2/search?jql=assignee=currentuser()+order+by+status+asc&expand=fields"

	body, err := get(url)
	if err != nil {
		system.Exit("There was a problem making the request to the jira API in `GetIssues`")
	}

	parsed := JiraIssues{}

	parseError := json.Unmarshal(body, &parsed)
	if parseError != nil {
		fmt.Printf("There was a problem parsing the jira API response:\n%s\n", parseError)
		system.Exit("")
	}
	return parsed
}

func GetIssue(issueNumber string) JiraIssue {
	defer system.ExecutionTimer(time.Now(), "GetIssue")
	orgName := viper.GetString("jira.orgname")
	url := "https://" + orgName + ".atlassian.net/rest/api/2/issue/" + issueNumber + "?expand=fields"

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
