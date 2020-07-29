package git

import (
	"sync"

	"trevoratlas.com/vor/utils"

	"github.com/spf13/viper"
)

type NativeGit interface {
	Call(string) (string, error)
	Stash() bool
	UnStash(string)
}

type GitClient struct {
	Path string
}

var once sync.Once
var client GitClient

func New() GitClient {
	once.Do(func() {
		localGit := viper.GetString("git.path")
		exists, fsErr := utils.Exists(localGit)
		if fsErr != nil || !exists {
			utils.Exit("Could not find local git client at " + "\"" + localGit + "\"")
		}
		client.Path = localGit

		_, gitErr := utils.Exec(client.Path + " status")
		if gitErr != nil {
			utils.Exit("Invalid git repository")
		}
	})
	return client
}

// Call – call a Client command by name
// you can pass arguments as well E.G:
// Client.Call("checkout -b my-branch-name")
// returns the text output of the command and a standard error (if any)
func (git GitClient) Call(command string) (string, error) {
	utils.Debug("calling 'git " + command + "'")
	return utils.Exec(git.Path + " " + command)
}

// Stash – stash changes if the working directory is unclean
func (git GitClient) Stash() {
	_, err := git.Call("stash")
	if err != nil {
		utils.Exit("error stashing changes")
	}
}

// UnStash – unstash the top most stash (called after a Stash())
func (git GitClient) UnStash() {
	_, err := git.Call("stash apply")
	if err != nil {
		utils.Exit("error unstashing changes")
	}
}
