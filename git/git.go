package git

import (
	"sync"

	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/utils"
)

func IsGitAvailable() (bool, error) {
	localGitPath := utils.GetStringEnv("git.path")
	return utils.Exists(localGitPath)
}

func IsGitRepo() bool {
	_, err := Call("status")
	return err == nil
}

// EnsureAvailability exits the program if:
// A. it cannot find git in the local filesystem, or
// B. you are not in a git repository
// otherwise it does nothing (noop)
func EnsureAvailability() {
	localGit, _ := IsGitAvailable()
	inRepo := IsGitRepo()

	if !localGit {
		utils.ExitWithMessage("could not find local git at \"/usr/local/bin/git\"")
	}

	if !inRepo {
		utils.ExitWithMessage("Git Status failed, are you in a valid git repository?")
	}
}

// Call – call a git command by name
// you can pass arguments as well E.G:
// git.Call("checkout -b my-branch-name")
// returns the text output of the command and a standard error (if any)
func Call(command string) (string, error) {
	localGitPath := utils.GetStringEnv("git.path")
	logger.Debug("calling git command " + command + " git path is: \"" + localGitPath + "\"")
	wg := new(sync.WaitGroup)
	wg.Add(2)
	return utils.ShellExec(localGitPath+" "+command, wg)
}
