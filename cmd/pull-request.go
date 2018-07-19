package commands

import (
	"encoding/json"
	"fmt"
	"github.com/trevor-atlas/vor/utils"
	"github.com/trevor-atlas/vor/logger"
	"regexp"
	"github.com/trevor-atlas/vor/git"
	"github.com/spf13/cobra"
)

var pullRequest = &cobra.Command{
	Use:   "pr",
	Short: "create a pull request with your current branch",
	Long: `create a pull request with your current branch against the default origin or a configured one`,
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureAvailability()
		base := utils.GetStringEnv("git.pull-request-base")
		githubAPIKey := utils.GetStringEnv("github.apikey")
		owner := utils.GetStringEnv("github.owner")
		if githubAPIKey == "" {
			utils.ExitWithMessage("No github API key found in vor config (github.apikey)")
		}
		if base == "" {
			fmt.Println("Repository base not found in config, falling back to origin/master...")
			base = "master"
		}
		remotes, _ := git.Call("remote -v")

		/* remotes will (likely) look like:
		origin	https://github.com/owner/project.git (fetch)
		origin	https://github.com/owner/project.git (push)
		*/
		matcher, _ := regexp.Compile(`github.com\/(.*)\/(.*)\.git`)
		res := matcher.FindAllStringSubmatch(remotes, 1)[0]
		logger.Debug("found matches", res)
		if owner == "" {
			owner = res[1]
		}
		repo := res[2]


		logger.Debug("git remote owner: " + owner + ", base: " + base)
		// git.StashExistingChanges()

	/*
	* master
	platform/AQ-4329/story/notifications-page
	*/
		branches, _ := git.Call("branch")
		branchMatcher, _ := regexp.Compile("\\* (.*)")
		branch := branchMatcher.FindAllStringSubmatch(branches, -1)[0][1]

		fmt.Println("s", branch)

		gpOutput, err := git.Call("push -u")
		if err != nil {
			utils.ExitWithMessage("Something went wrong pushing to github:\n" + gpOutput)
		}

		b, err := json.Marshal(git.PullRequestBody{
				Title: "currentBranchName",
				Body: "Created automagically by Vor",
				Head: owner+":"+branch,
				Base: base,
		})
		if err != nil {
			logger.Debug()
			utils.ExitWithMessage("There was a problem marshaling the JSON to create the pull request!")
		}
		  git.Post("https://api.github.com/repos/"+owner+"/"+repo+"/pulls", b)
	},
}

func init() {
	rootCmd.AddCommand(pullRequest)
}

