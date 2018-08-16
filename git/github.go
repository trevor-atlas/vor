package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/system"
	"github.com/trevor-atlas/vor/utils"
	"io/ioutil"
	"net/http"
)

func GeneratePRName(branchName string) string {
	return utils.TitleCase(branchName)
}

func Post(url string, requestBody []byte) PullRequestResponse {
	githubAPIKey := system.GetString("github.apikey")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", "token "+githubAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		system.Exit("error parsing github response! %s", err)
	}
	parsed := PullRequestResponse{}
	if utils.Contains(string(body), "No commits between") {
		system.Exit("Your branch is not changed from the base branch!")
	}
	log := logger.New()
	log.Debug("response Status: %s", resp.Status)
	log.Debug("response Headers: %s", resp.Header)
	log.Debug("response Body: %s", string(body))

	parseError := json.Unmarshal(body, &parsed)
	if parseError != nil {
		fmt.Printf("error parsing json\n %s", parseError)
		panic(parseError)
	}
	fmt.Println(parsed)
	return parsed
}
