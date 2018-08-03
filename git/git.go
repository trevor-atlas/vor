package git

import (
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/system"
	"github.com/trevor-atlas/vor/utils"
)

type NativeGit interface {
	Call(string) (string, error)
	Stash() bool
	UnStash(string)
}

type Git struct {
	Path string
	Sys  system.OSHandler
}

func New(sys system.OSHandler, env system.Envloader) *Git {
	g := new(Git)
	localGit := env.String("git.path")
	exists, fsErr := sys.Exists(localGit)
	if fsErr != nil || !exists {
		sys.Exit("Could not find local git at " + "\"" + localGit + "\"")
	}
	g.Path = localGit
	g.Sys = sys

	_, gitErr := sys.Exec(g.Path + " status")
	if gitErr != nil {
		sys.Exit("Invalid git repository")
	}
	return g
}

// Call – call a git command by name
// you can pass arguments as well E.G:
// git.Call("checkout -b my-branch-name")
// returns the text output of the command and a standard error (if any)
func (git *Git) Call(command string) (string, error) {
	logger.Debug("calling 'git " + command + "'")
	return git.Sys.Exec(git.Path + " " + command)
}

// Stash – stash changes if the working directory is unclean
func (git *Git) Stash() (didStash bool) {
	cmdOutput, _ := git.Call("status")
	contains := func(substr string) bool { return utils.Contains(cmdOutput, substr) }

	if contains("deleted") || contains("modified") || contains("untracked") {
		affirmed := git.Sys.Confirm("Working directory is not clean. Stash changes?")
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
	affirm := git.Sys.Confirm(message)
	if affirm {
		git.Call("stash apply")
	}
}
