package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/trevor-atlas/vor/git"
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/system"
	"regexp"
)

var pullRequest = &cobra.Command{
	Use:   "pr",
	Short: "create a pull request with your current branch",
	Long:  `create a pull request with your current branch against the default origin or a configured one`,
	Run: func(cmd *cobra.Command, args []string) {
		gc := git.New()
		log := logger.New()
		base := system.GetString("git.pull-request-base")
		githubAPIKey := system.GetString("github.apikey")
		owner := system.GetString("github.owner")
		if githubAPIKey == "" {
			system.Exit("No github API key found in vor config (github.apikey)")
		}
		if base == "" {
			fmt.Println("Repository base not found in config, falling back to origin/master...")
			base = "master"
		}
		remotes, _ := gc.Call("remote -v")

		/* remotes will (likely) look like:
		origin	https://github.com/owner/project.git (fetch)
		origin	https://github.com/owner/project.git (push)
		*/
		matcher, _ := regexp.Compile(`github.com\/(.*)\/(.*)\.git`)
		res := matcher.FindAllStringSubmatch(remotes, 1)[0]
		log.Debug("found matches %s", res)
		if owner == "" {
			owner = res[1]
		}
		repo := res[2]

		log.Debug("git remote owner: " + owner + ", base: " + base)
		// git.StashExistingChanges()

		/*
			* master
			platform/AQ-4329/story/notifications-page
		*/
		branches, _ := gc.Call("branch")
		branchMatcher, _ := regexp.Compile("\\* (.*)")
		branch := branchMatcher.FindAllStringSubmatch(branches, -1)[0][1]

		fmt.Println("s", branch)

		gpOutput, err := gc.Call("push -u")
		fmt.Println(gpOutput, err)
		// If the upstread is not set, do it for us! otherwise panic cause this is weird
		if err != nil {
			gc.Call("push --set-upstream origin " + branch)
			// utils.ExitWithMessage("Something went wrong pushing to github:\n" + gpOutput)
		}

		b, err := json.Marshal(git.PullRequestBody{
			Title: branch,
			Body:  "Created automagically by Vor",
			Head:  owner + ":" + branch,
			Base:  base,
		})
		if err != nil {
			log.Debug("there was a problem unmarshalling the git response %s", err)
			system.Exit("There was a problem marshaling the JSON to create the pull request!")
		}
		git.Post("https://api.github.com/repos/"+owner+"/"+repo+"/pulls", b)
	},
}

func init() {
	rootCmd.AddCommand(pullRequest)
}
