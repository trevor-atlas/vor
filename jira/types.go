package jira

import (
	"time"
	"net/http"
	"io/ioutil"
)

// Time represents the Time definition of JIRA as a time.Time of go
type Time time.Time

// Date represents the Date definition of JIRA as a time.Time of go
type Date time.Time

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
	Created      *Time
	Updated      *Time
	Total        int
}

// JiraIssue describes the response for a single jira issue
type JiraIssue struct {
	ID     string `json:"id"`
	Self   string `json:"self"` // url to request this issue
	Key    string `json:"key"`  // AQ-XXXX
	Fields struct {
		Summary           string // title of jira issue
		Created           *Time `json:"created"` // 2018-05-25T04:18:06.836-0500
		Updated           *Time `json:"updated"` // 2018-06-11T22:23:03.606-0500
		Description       string // description of Jira issue
		Reporter jiraUser
		Assignee jiraUser
		Customfield_12022 struct {
			Value string // team name
		}
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
		}
		Status struct {
			Description string
			Name string
			StatusCategory struct {
				Key string
				Name string
				ID int
			}
		}
		Project struct {
			Key string
			Name string
		}
	} `json:"fields"`
}

type JiraIssues struct {
	Issues []JiraIssue
}

type HttpResponseFetcher interface {
	Fetch(url string) ([]byte, error)
}

func (h *HTTP) Fetch(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (h *HTTP) FetchWithHeaders(url string, req http.Request, redirectHandler func(req *http.Request, via []*http.Request) error) ([]byte, error) {
	client := &http.Client{
		CheckRedirect: redirectHandler,
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(&req)
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
