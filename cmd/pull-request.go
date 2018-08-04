package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/trevor-atlas/vor/git"
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/system"
	"github.com/trevor-atlas/vor/utils"
	"regexp"
)

var pullRequest = &cobra.Command{
	Use:   "pr",
	Short: "create a pull request with your current branch",
	Long:  `create a pull request with your current branch against the default origin or a configured one`,
	Run: func(cmd *cobra.Command, args []string) {
		sys := system.NewOSHandler()
		git := git.New(sys, system.Env)
		base := system.Env.GetString("git.pull-request-base")
		githubAPIKey := system.Env.GetString("github.apikey")
		owner := system.Env.GetString("github.owner")
		if githubAPIKey == "" {
			system.NewOSHandler().Exit("No github API key found in vor config (github.apikey)")
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
		logger.Debug("found matches %s", res)
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
		fmt.Println(gpOutput, err)
		// If the upstread is not set, do it for us! otherwise panic cause this is weird
		if err != nil {
			git.Call("push --set-upstream origin " + branch)
			// utils.ExitWithMessage("Something went wrong pushing to github:\n" + gpOutput)
		}

		b, err := json.Marshal(git.PullRequestBody{
			Title: branch,
			Body:  "Created automagically by Vor",
			Head:  owner + ":" + branch,
			Base:  base,
		})
		if err != nil {
			logger.Debug("there was a problem unmarshalling the git response %s", err)
			utils.ExitWithMessage("There was a problem marshaling the JSON to create the pull request!")
		}
		git.Post("https://api.github.com/repos/"+owner+"/"+repo+"/pulls", b)
	},
}

func init() {
	rootCmd.AddCommand(pullRequest)
}
