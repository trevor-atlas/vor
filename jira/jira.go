package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"

	"github.com/trevor-atlas/vor/system"
	"github.com/trevor-atlas/vor/utils"
)

const (
	top_left           string = "\u256D"
	top_right          string = "\u256E"
	bottom_left        string = "\u2570"
	bottom_right       string = "\u256F"
	x_line             string = "\u2500"
	y_line             string = "\u2502"
	bottom_left_sharp  string = "\u2514"
	bottom_right_sharp string = "\u2518"
	shade              string = "\u2591"
	left_quote         string = "\u201C"
	right_quote        string = "\u201D"
	thick_underscore   string = "\u2581"
	max_len            int    = 70
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectHandler(req *http.Request, via []*http.Request) error {
	jiraUsername := system.GetString("jira.username")
	jiraKey := system.GetString("jira.apikey")
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
	orgName := system.GetString("jira.orgname")
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
	orgName := system.GetString("jira.orgname")
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

	w(title)
	wpnl(cyan("issue: ") + issue.Key)
	wpnl(cyan("type: ") + issue.Fields.IssueType.Name)
	wpnl(cyan("status: ") + issue.Fields.Status.Name)
	wpnl(cyan("reporter: ") + magenta(issue.Fields.Reporter.Name))
	wpnl(cyan("assignee: ") + magenta(issue.Fields.Assignee.Name))
	wpnl(cyan("created: ") + yellow(time.Time(*issue.Fields.Created).Format("2006-01-02 15:04")))
	wpnl(cyan("updated: ") + yellow(humanize.Time(time.Time(*issue.Fields.Updated))))
	wpnl(cyan("url: ") + blue(issueURL))
	wpnl(cyan("description:"))
	w(formatMultiline(issue.Fields.Description, utils.PadOutput(4)) + "\n\n")

	if len(issue.Fields.Comment.Comments) > 0 {
		nestedPad := utils.PadOutput(4)
		wpnl(cyan("comments:"))
		for _, comment := range issue.Fields.Comment.Comments {
			w(nestedPad(cyan("author: ") + comment.Author.Name + "\n"))
			w(nestedPad(cyan("created: ")+yellow(time.Time(*comment.Created).Format("2006-01-02 15:04"))) + "\n")
			w(nestedPad(cyan("updated: ")+yellow(humanize.Time(time.Time(*comment.Updated)))) + "\n")
			w(nestedPad(cyan("body:\n")))
			w(formatMultiline(comment.Body, utils.PadOutput(6)) + "\n\n")
		}
	}

	result := b.String()
	fmt.Println(result)
	return result
}

type HTTP struct{}

func Get(url string) (*http.Response, error) {
	jiraUsername := system.GetString("jira.username")
	jiraKey := system.GetString("jira.apikey")

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
	defer system.ExecutionTimer(time.Now(), "GetIssues")
	orgname := system.GetString("jira.orgname")
	if orgname == "" {
		system.Exit("jira.orgname config not found.")
	}
	username := system.GetString("jira.username")
	if username == "" {
		system.Exit("jira.username config not found.")
	}
	apikey := system.GetString("jira.apikey")
	if apikey == "" {
		system.Exit("jira.apikey config not found.")
	}
	url := "https://" + orgname + ".atlassian.net/rest/api/2/search?jql=assignee=currentuser()+order+by+status+asc&expand=fields"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(username, apikey))
	body, err := HTTP{}.FetchWithHeaders(url, *req, redirectHandler)
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
	orgName := system.GetString("jira.orgname")
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

// UnmarshalJSON will transform the JIRA time into a time.Time
// during the transformation of the JIRA JSON response
func (t *Time) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil
	}
	ti, err := time.Parse("\"2006-01-02T15:04:05.999-0700\"", string(b))
	if err != nil {
		return err
	}
	*t = Time(ti)
	return nil
}

// UnmarshalJSON will transform the JIRA date into a time.Time
// during the transformation of the JIRA JSON response
func (t *Date) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil
	}
	ti, err := time.Parse("\"2006-01-02\"", string(b))
	if err != nil {
		return err
	}
	*t = Date(ti)
	return nil
}
