package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/trevor-atlas/vor/env"
	"github.com/trevor-atlas/vor/rest"
	"github.com/trevor-atlas/vor/utils"
	"net/http"
	"regexp"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trevor-atlas/vor/git"
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/system"
)

var gc git.GitClient
var prMessage string
var prTitle string
var prJson bool

func getRemotesMeta() (owner, repo string) {
	/* remotes will (likely) look like:
	origin	https://github.com/owner/project.git (fetch)
	origin	https://github.com/owner/project.git (push)
	*/
	remotes, _ := gc.Call("remote -v")
	matcher, _ := regexp.Compile(`github.com\/(.*)\/(.*)\.git`)
	res := matcher.FindAllStringSubmatch(remotes, 1)[0]
	logger.Debug("found matches %s", res)
	_owner := viper.GetString("github.owner")
	if owner == "" {
		owner = res[1]
	}
	_repo := res[2]

	return _owner, _repo
}

func pr(args []string) {
	gc = git.New()
	owner, repo := getRemotesMeta()
	branch := getLocalBranchName()
	if branch == "" {
		system.Exit("something went wrong getting the local branch name!")
	}
	if prTitle == "" {
		prTitle = git.GeneratePRName(branch)
	}

	pushResult, pushErr := gc.Call("push")
	// If the upstread is not set, do it for us! otherwise panic cause this is weird
	if pushErr != nil {
		_, err := gc.Call("push --set-upstream origin " + branch)
		if err != nil {
			system.Exit("error calling local git")
		}
		system.Exit("Something went wrong pushing to github:\n" + pushResult)
	}

	body, marshalErr := json.Marshal(git.PullRequestBody{
		Title: prTitle,
		Body:  prMessage,
		Head:  owner + ":" + branch,
		Base:  env.PULL_REQUEST_BASE,
	})
	if marshalErr != nil {
		system.Exit("There was a problem marshaling the JSON to create the pull request!")
	}
	res, resErr := git.Post(
		rest.NewHTTPClient(&http.Client{Timeout: time.Second * 10}).
			URL("https://api.github.com/repos/"+owner+"/"+repo+"/pulls").
			WithHeader("Authorization", "token "+env.GITHUB_APIKEY).
			WithHeader("Content-Type", "application/json").
			BODY(bytes.NewBuffer(body)).
			POST,
		)

	if resErr != nil {
		system.Exit("")
	}

	if res.Errors != nil {
		logger.Error("Something went wrong creating your pull request")
		if utils.Contains(res.Message, "No commits between") {
			system.Exit("Your branch is not changed from the base branch!")
		}
		if res.Message == "Validation Failed" {
			logger.Info("A pull request already exists for this branch")
			fmt.Println("https://github.com/"+owner+"/"+repo+"/pull/"+branch)
		} else {
			res.PrintJSON()
		}
		system.Exit("")
	}

	if prJson {
		res.PrintJSON()
	} else {
		fmt.Printf("Pull request created at: %s\n", res.HTMLURL)
	}
}

func getLocalBranchName() string {
	/*
		* master
		XX-1234/story/notifications
	*/
	branches, _ := gc.Call("branch")
	branchMatcher, _ := regexp.Compile("\\* (.*)")
	branch := branchMatcher.FindAllStringSubmatch(branches, -1)[0][1]
	return branch
}

var pullRequest = &cobra.Command{
	Use:     "pull-request",
	Aliases: []string{"pr", "pull"},
	Example: "vor pr",
	Short:   "create a pull request with your current branch",
	Long:    `create a pull request with your current branch against the default origin or a configured one`,
	Run: func(cmd *cobra.Command, args []string) {
		pr(args)
	},
}

func init() {
	pullRequest.Flags().BoolVarP(&prJson, "json", "j", false, "output json response")
	pullRequest.Flags().StringVarP(&prTitle, "title", "t", "", "optional title for pull request (defaults to branch name)")
	pullRequest.Flags().StringVarP(&prMessage, "message", "m", "Created automagically by Vor", "optional message for pull request description")
	rootCmd.AddCommand(pullRequest)
}
