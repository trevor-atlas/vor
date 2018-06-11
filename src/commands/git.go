package commands

import (
	"fmt"
	"github.com/trevor-atlas/vor/src/utils"
)

func GitStatus(wg *sync.WaitGroup) {
	res, err := ShellExec("/usr/local/bin/hub log", wg)
	if err != nil {
		LoggerError("Git Status failed, are you in a valid git repository?")
	}
	fmt.Println(res)
}