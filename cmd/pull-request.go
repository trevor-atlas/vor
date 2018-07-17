package commands

import (
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
		// githubAPIKey := utils.GetStringEnv("github.apikey")
		// if githubAPIKey == "" {
			// utils.ExitWithMessage("No github API key found in vor config (github.apikey)")
		// }
		remotes, err := git.Call("remote -v")
		if err != nil {
			fmt.Println("err", err)
		}

		/* remotes will (likely) look like:
		origin	https://github.com/owner/project.git (fetch)
		origin	https://github.com/owner/project.git (push)
		*/
		re, _ := regexp.Compile(`github.com\/(.*)\/.*`)
		res := re.FindAllStringSubmatch(remotes, 1)
		fmt.Println(res)
		owner := res[0][1]

		if base == "" {
			fmt.Println("Repository base not found in config, falling back to origin/master...")
			base = "master"
		}
		logger.Debug("git remote owner: " + owner + ", base: " + base)
		// git.StashExistingChanges()


	},
}

func init() {
	rootCmd.AddCommand(pullRequest)
}

