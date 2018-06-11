package util

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func trace() {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	shouldStop := false
	fmt.Println("\nstack trace:")
	for shouldStop != true {
		frame, more := frames.Next()
		fmt.Printf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
			shouldStop = !more
	}
}

func loggerError(message string, rest ...interface{}) {
	fmt.Printf(message, rest...)
}

func shellExec(cmd string, wg *sync.WaitGroup) (string, error) {
	defer wg.Done() // Need to signal to waitgroup that this goroutine is done

	parts := strings.Fields(cmd)
	head := parts[0]
	args := parts[1:]

	out, err := exec.Command(head, args...).Output()
	if err != nil {
		loggerError("ERROR: calling exec with '%s'\ncommand failed with %s\n", cmd, err)
		trace()
	}
	result := string(out)
	return result, err
}

func gitStatus(wg *sync.WaitGroup) {
	res, err := shellExec("/usr/local/bin/hub log", wg)
	if err != nil {
		loggerError("Git Status failed, are you in a valid git repository?")
	}
	fmt.Println(res)
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	gitStatus(wg)
}
