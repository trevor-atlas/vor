package git

import (
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/system"
	"github.com/trevor-atlas/vor/utils"
	"sync"
)

type NativeGit interface {
	Call(string) (string, error)
	Stash() bool
	UnStash(string)
}

type gitClient struct {
	Path string
}

var once sync.Once
var client gitClient

func New() gitClient {
	once.Do(func() {
		localGit := system.GetString("git.path")
		exists, fsErr := system.Exists(localGit)
		if fsErr != nil || !exists {
			system.Exit("Could not find local git client at " + "\"" + localGit + "\"")
		}
		client.Path = localGit

		_, gitErr := system.Exec(client.Path + " status")
		if gitErr != nil {
			system.Exit("Invalid git repository")
		}
	})
	return client
}

// Call – call a Client command by name
// you can pass arguments as well E.G:
// Client.Call("checkout -b my-branch-name")
// returns the text output of the command and a standard error (if any)
func (git gitClient) Call(command string) (string, error) {
	log := logger.New()
	log.Debug("calling 'Client " + command + "'")
	return system.Exec(git.Path + " " + command)
}

// Stash – stash changes if the working directory is unclean
func (git gitClient) Stash() (didStash bool) {
	cmdOutput, _ := git.Call("status")
	c := func(substr string) bool { return utils.Contains(cmdOutput, substr) }

	if c("deleted") || c("modified") || c("untracked") {
		affirmed := system.Confirm("Working directory is not clean. Stash changes?")
		if !affirmed {
			return false
		}
		git.Call("stash")
		return true
	}
	return false
}

// UnStash – unstash the top most stash (called after a Stash())
func (git gitClient) UnStash(message string) {
	affirm := system.Confirm(message)
	if affirm {
		git.Call("stash apply")
	}
}
