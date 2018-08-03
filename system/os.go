package system

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/utils"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type OSHandler interface {
	Confirm(string) bool
	Exit(string)
	Exists(string) (bool, error)
	Exec(string) (string, error)
	ExecutionTimer(time.Time, string)
}

// ensure Sys implements the OSHandler interface
var _ OSHandler = (*Sys)(nil)

type Sys struct{}

func NewOSHandler() *Sys {
	return new(Sys)
}

func (s *Sys) ExecutionTimer(start time.Time, name string) {
	elapsed := time.Since(start)
	logger.Debug("%s took %s", name, elapsed)
}

// Confirm prompt the user with {message} (Y/N) and return true for Y, false for N (case insensitive)
func (s *Sys) Confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message + " [Y/N]")
	text, _ := reader.ReadString('\n')
	return utils.Contains(text, "Y") || utils.Contains(text, "yes")
}

func (s *Sys) Exit(message string) {
	color.Red(message + "\ncanceling operation...")
	os.Exit(1)
}

// Exists returns whether the given file or directory exists or not
func (s *Sys) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (s *Sys) Exec(cmd string) (string, error) {
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
