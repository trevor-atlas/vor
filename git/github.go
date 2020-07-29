package git

import (
	"fmt"
	"strings"


	"trevoratlas.com/vor/utils"
)

func GeneratePRName(branchName string) string {
	r := strings.NewReplacer(
		"-", " ",
		"/", " ",
		"\n", "",
		"\t", "")
	return utils.StringPipe(
		r.Replace,
		utils.Capitalize,
	)(branchName)
}

func Post(post func() ([]byte, error)) (PullRequestResponse, error) {
	response, requestErr := post()
	if requestErr != nil {
		utils.Error("Something went wrong talking to github: %s\n", requestErr)
		utils.Exit("")
	}
	parsed := PullRequestResponse{}
	parseError := parsed.Unmarshal(response)
	if parseError != nil {
		fmt.Printf("error parsing json\n %s", parseError)
		return parsed, parseError
	}

	return parsed, nil
}
