package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/fatih/color"
	"github.com/trevor-atlas/vor/logger"
)


func PadOutput(padding int) func(string) string {
	return func(str string) string {
		return LeftPad(str, " ", padding)
	}
}


func KebabCase(s string) string {
	r := strings.NewReplacer(
" ", "-",
		":", "",
		"/", "-",
		"\n", "",
		"	", "-")
	return r.Replace(s)
}

func LowerKebabCase(s string) string {
	return strings.ToLower(KebabCase(s))
}

// GetStringEnv get a specific env and fallback to the global env if it does not exist locally
func GetStringEnv(specificEnv string) string {
	env := viper.GetString(specificEnv)
	if env == "" {
		env = viper.GetString("global." + specificEnv)
	}
	return env
}

// PromptYesNo prompt the user with {message} (Y/N) and return true for Y, false for N (case insensitive)
func PromptYesNo(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message + " (Y/N)")
	text, _ := reader.ReadString('\n')
	return CaseInsensitiveContains(text, "N")
}

func ExitWithMessage(message string) {
	color.Red(message + "\ncanceling operation...")
	os.Exit(1)
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

// Exists returns whether the given file or directory exists or not
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Trace spits out the current stack trace when called
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

// ShellExec attempts to execute a given shell command
func ShellExec(cmd string, wg *sync.WaitGroup) (string, error) {
	defer wg.Done() // Need to signal to waitgroup that this goroutine is done

	parts := strings.Fields(cmd)
	head := parts[0]
	args := parts[1:]

	out, err := exec.Command(head, args...).Output()
	if err != nil {
		logger.Error("ERROR: calling exec with '%s'\ncommand failed with %s\n", cmd, err)
		Trace()
	}
	result := string(out)
	return result, err
}

// LeftPad pad s with pLen number of padStr
func LeftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}

// RightPad pad s with pLen number of padStr
func RightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, pLen)
}
