package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"trevoratlas.com/vor/formatters"
)

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing.
func CheckArgs(got []string, wanted ...string) {
	if len(got) < len(wanted) {
		Exit("Not enough arguments, wanted: " + strings.Join(wanted[:],", "))
	}
}

// ExecutionTimer a useful perf tool,
// call this at the start of a suspect method to trace how long it takes to execute
func ExecutionTimer(start time.Time, name string) {
	elapsed := time.Since(start)
	Debug("%s took %s", name, elapsed)
}

// Confirm prompt the user with {message} (Y/N) and return true for Y, false for N (case insensitive)
func Confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message + " [Y/N]")
	text, _ := reader.ReadString('\n')
	return Contains(text, "Y") || Contains(text, "yes")
}

func Exit(format string, rest ...interface{}) {
	fmt.Println(formatters.RED(format, rest...))
	os.Exit(1)
}

