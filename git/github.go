package git

import (
	"fmt"
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/system"
	"github.com/trevor-atlas/vor/utils"
)

func GeneratePRName(branchName string) string {
	return utils.Capitalize(branchName)
}

func Post(post func() ([]byte, error)) (PullRequestResponse, error) {
	response, requestErr := post()
	if requestErr != nil {
		logger.Error("Something went wrong talking to github: %s\n", requestErr)
		system.Exit("")
	}
	parsed := PullRequestResponse{}
	parseError := parsed.Unmarshal(response)
	if parseError != nil {
		fmt.Printf("error parsing json\n %s", parseError)
		return parsed, parseError
	}

	return parsed, nil
}
