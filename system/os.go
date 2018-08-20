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

func ExecutionTimer(start time.Time, name string) {
	log := logger.New()
	elapsed := time.Since(start)
	log.Debug("%s took %s", name, elapsed)
}

// Confirm prompt the user with {message} (Y/N) and return true for Y, false for N (case insensitive)
func Confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message + " [Y/N]")
	text, _ := reader.ReadString('\n')
	return utils.Contains(text, "Y") || utils.Contains(text, "yes")
}

func Exit(message string) {
	color.Red(message + "\ncanceling operation...")
	os.Exit(1)
}

// Exists returns whether the given file or directory exists or not
func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

func Exec(cmd string) (string, error) {
	log := logger.New()
	parts := strings.Fields(cmd)
	head := parts[0]
	args := parts[1:]
	wg := new(sync.WaitGroup)
	wg.Add(2)
	defer wg.Done()

	out, err := exec.Command(head, args...).Output()
	if err != nil {
		log.Error("ERROR: calling exec with '%s':\ncommand failed with %s\n", cmd, err)
	}

	return string(out), err
}
