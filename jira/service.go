package jira

import (
	"encoding/json"
	"fmt"
	"trevoratlas.com/vor/rest"
	"trevoratlas.com/vor/utils"
)

type Service struct{}

func (s *Service) GetIssues() JiraIssues {
	url := "https://" + utils.JIRA_ORGNAME + ".atlassian." +
		"net/rest/api/2/search?jql=assignee=currentuser()+order+by+status+asc&expand=fields"
	client := rest.New().
		WithHeader("Accept", "application/json").
		WithBasicAuth(utils.JIRA_USERNAME, utils.JIRA_APIKEY).
		// WithHandler(func(req *http.Request, via []*http.Request) error {
		// req.Header.Add("Authorization", "Basic "+BasicAuth(utils.JIRA_USERNAME, utils.JIRA_APIKEY))
		// return nil
		//}).
		URL(url)

	body, err := client.GET()
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

func (s *Service) GetIssue(issueNumber string) JiraIssue {
	url := "https://" + utils.JIRA_ORGNAME + ".atlassian." +
		"net/rest/api/2/issue/" + issueNumber + "?expand=fields"

	client := rest.New().
		WithHeader("Accept", "application/json").
		WithBasicAuth(utils.JIRA_USERNAME, utils.JIRA_APIKEY).
		// WithHandler(func(req *http.Request, via []*http.Request) error {
		// req.Header.Add("Authorization", "Basic "+BasicAuth(utils.JIRA_USERNAME, utils.JIRA_APIKEY))
		// return nil
		// }).
		URL(url)

	res, err := client.GET()
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
