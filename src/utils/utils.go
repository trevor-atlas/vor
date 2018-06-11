package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// GetEnv does stuff
func GetEnv(name string) string {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == name {
			return pair[1]
		}
	}
	return ""
}

// Trace does stuff
func Trace() {
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

// LoggerError does stuff
func LoggerError(message string, rest ...interface{}) {
	fmt.Printf(message, rest...)
}

// ShellExec does stuff
func ShellExec(cmd string, wg *sync.WaitGroup) (string, error) {
	defer wg.Done() // Need to signal to waitgroup that this goroutine is done

	parts := strings.Fields(cmd)
	head := parts[0]
	args := parts[1:]

	out, err := exec.Command(head, args...).Output()
	if err != nil {
		LoggerError("ERROR: calling exec with '%s'\ncommand failed with %s\n", cmd, err)
		Trace()
	}
	result := string(out)
	return result, err
}

// LeftPad does stuff
func LeftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}
