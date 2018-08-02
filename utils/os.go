package utils

import (
	"time"
	"github.com/trevor-atlas/vor/logger"
	"bufio"
	"os"
	"fmt"
	"github.com/fatih/color"
	"strings"
	"sync"
	"os/exec"
)

type SideEffects interface {
	Confirm(string) bool
	Exit(string)
	Exists(string) (bool, error)
	Exec(string)(string, error)
	ExecutionTimer(time.Time, string)
}

type OS struct{}

func (o *OS) ExecutionTimer(start time.Time, name string) {
	elapsed := time.Since(start)
	logger.Debug("%s took %s", name, elapsed)
}

// Confirm prompt the user with {message} (Y/N) and return true for Y, false for N (case insensitive)
func (o *OS) Confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message + " [Y/N]")
	text, _ := reader.ReadString('\n')
	return Contains(text, "Y") || Contains(text, "yes")
}

func (o *OS) Exit(message string) {
	color.Red(message + "\ncanceling operation...")
	os.Exit(1)
}

// Exists returns whether the given file or directory exists or not
func (o *OS) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (o *OS) Exec(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	head := parts[0]
	args := parts[1:]
	wg := new(sync.WaitGroup)
	wg.Add(2)
	defer wg.Done()

	out, err := exec.Command(head, args...).Output()
	if err != nil {
		logger.Error("ERROR: calling exec with '%s':\ncommand failed with %s\n", cmd, err)
	}

	return string(out), err
}

