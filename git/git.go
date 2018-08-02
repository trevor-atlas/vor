package git

import (
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/utils"
)

type NativeGit interface {
	Call(string) (string, error)
	Stash() bool
	UnStash(string)
}

type Git struct{
	path string
}

func NewGit() *Git {
	g := new(Git)
	osutil := utils.OS{}
	localGit := utils.ENV{}.String("git.path")
	exists, fsErr := osutil.Exists(localGit)
	if fsErr != nil || !exists {
		osutil.Exit("Could not find local git at " + localGit)
	}
	g.path = localGit

	_, gitErr := osutil.Exec(g.path + "status")
	if gitErr != nil {
		osutil.Exit("Invalid git repository")
	}
	return g
}

// Call – call a git command by name
// you can pass arguments as well E.G:
// git.Call("checkout -b my-branch-name")
// returns the text output of the command and a standard error (if any)
func (git *Git) Call(command string) (string, error) {
	logger.Debug("calling 'git " + command + "'")
	return utils.OS{}.Exec(git.path+" "+command)
}

// Stash – stash changes if the working directory is unclean
func (git *Git) Stash() (didStash bool) {
	cmdOutput, _ := git.Call("status")
	contains := func(substr string) bool { return utils.Contains(cmdOutput, substr) }

	if contains("deleted") || contains("modified") || contains("untracked") {
		affirmed := utils.OS{}.Confirm("Working directory is not clean. Stash changes?")
		if !affirmed {
			return false
		}
		git.Call("stash")
		return true
	}
	return false
}

// UnStash – unstash the top most stash (called after a Stash())
func (git *Git) UnStash(message string) {
	affirm := utils.OS{}.Confirm(message)
	if affirm {
		git.Call("stash apply")
	}
}
