package jira

import (
			"time"
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
		Created           *Time  `json:"created"` // 2018-05-25T04:18:06.836-0500
		Updated           *Time  `json:"updated"` // 2018-06-11T22:23:03.606-0500
		Description       string // description of Jira issue
		Reporter          jiraUser
		Assignee          jiraUser
		Customfield_12022 struct {
			Value string // team name
		}
		Comment struct {
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
			Description    string
			Name           string
			StatusCategory struct {
				Key  string
				Name string
				ID   int
			}
		}
		Project struct {
			Key  string
			Name string
		}
	} `json:"fields"`
}

type JiraIssues struct {
	Issues []JiraIssue
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
