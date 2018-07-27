package git

import (
	"encoding/json"
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/utils"
	"io/ioutil"
	"fmt"
	"bytes"
	"net/http"
)

func GeneratePRName(branchName string) string {
	return utils.TitleCase(branchName)
}

func Post (url string, requestBody []byte) PullRequestResponse {
	githubAPIKey := utils.GetStringEnv("github.apikey")
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
    req.Header.Set("Authorization", "token " + githubAPIKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	parsed := PullRequestResponse{}
	// if utils.CaseInsensitiveContains(string(resp.Body), "No commits between") {
		// utils.ExitWithMessage("Your branch is not changed from the base branch!")
	// }
	logger.Debug("response Status:", resp.Status)
	logger.Debug("response Headers:", resp.Header)
	logger.Debug("response Body:", string(body))

	parseError := json.Unmarshal(body, &parsed)
	if parseError != nil {
		fmt.Printf("error parsing json\n %s", parseError)
		panic(parseError)
	}
	fmt.Println(parsed)
	return parsed

}