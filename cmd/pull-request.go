package commands

import (
			"github.com/spf13/cobra"
	"github.com/trevor-atlas/vor/git"
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/system"
	"fmt"
	"regexp"
	"encoding/json"
)

var gc git.GitClient
var log logger.Logger

func getRemotesMeta() (owner, repo string) {
	/* remotes will (likely) look like:
	origin	https://github.com/owner/project.git (fetch)
	origin	https://github.com/owner/project.git (push)
	*/
	remotes, _ := gc.Call("remote -v")
	matcher, _ := regexp.Compile(`github.com\/(.*)\/(.*)\.git`)
	res := matcher.FindAllStringSubmatch(remotes, 1)[0]
	log.Debug("found matches %s", res)
	_owner := system.GetString("github.owner")
	if owner == "" {
		owner = res[1]
	}
	_repo := res[2]

	return _owner, _repo
}

func pr(args []string) {
	var prMessage string
	var prTitle string

	rootCmd.Flags().StringVarP(&prMessage, "message", "m", "Created automagically by Vor", "optional message for pull request description")

	githubAPIKey := system.GetString("github.apikey")
	if githubAPIKey == "" {
		system.Exit("No github API key found in vor config (github.apikey)")
	}
	base := system.GetString("git.pull-request-base")
	if base == "" {
		fmt.Println("Repository base not found in config, falling back to origin/master...")
		base = "master"
	}
	owner, repo := getRemotesMeta()
	
	branch := getLocalBranchName()
	if branch != "" {
		rootCmd.Flags().StringVarP(&prTitle, "title", "t", branch, "optional title for pull request (defaults to branch name)")
	} else {
		system.Exit("something went wrong getting the local branch name!")
	}

	gpOutput, err := gc.Call("push -u")
	fmt.Println(gpOutput, err)
	// If the upstread is not set, do it for us! otherwise panic cause this is weird
	if err != nil {
		_, err := gc.Call("push --set-upstream origin " + branch)
		if err != nil {
			system.Exit("error calling local git")
		}
		// utils.ExitWithMessage("Something went wrong pushing to github:\n" + gpOutput)
	}

	b, err := json.Marshal(git.PullRequestBody{
		Title: prTitle,
		Body:  prMessage,
		Head:  owner + ":" + branch,
		Base:  base,
	})
	if err != nil {
		log.Debug("there was a problem unmarshalling the git response %s", err)
		system.Exit("There was a problem marshaling the JSON to create the pull request!")
	}
	git.Post("https://api.github.com/repos/"+owner+"/"+repo+"/pulls", b)
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
	Use:   "pull-request",
	Aliases: []string{"pr", "pull"},
	Example: "vor pr",
	Short: "create a pull request with your current branch",
	Long:  `create a pull request with your current branch against the default origin or a configured one`,
	Run: func(cmd *cobra.Command, args []string) {

		pr(args)
	},
}

func init() {
	gc = git.New()
	log = *logger.New()
	rootCmd.AddCommand(pullRequest)
}
