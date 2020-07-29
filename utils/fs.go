package utils

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Exists returns whether the given file or directory exists or not
func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

func Exec(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	head := parts[0]
	args := parts[1:]
	wg := new(sync.WaitGroup)
	wg.Add(2)
	defer wg.Done()

	out, err := exec.Command(head, args...).Output()
	if err != nil {
		Error("ERROR: calling exec with '%s':\ncommand failed with %s\n", cmd, err)
	}

	return string(out), err
}

func WalkUpFS(cutoff string) (match string, err error) {
	// the directory we started in
	startDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// split it into directories
	path := strings.Split(startDir, "/")
	cutoff = filepath.Clean(cutoff)

	for i := len(path) - 1; i > 0; i-- {
		if cutoff == path[i] || i == 0 { // if we are at the cutoff directory, stop
			return "", nil
		} else { // otherwise, walk the path from i
			// it's worth noting that slice indices with ':' are length based,
			// not index based (normally i + 1 would be out of range on the first iteration)
			here := strings.Join(path[:i + 1], "/")
			Debug("Walking " + here)
			matches, err := filepath.Glob(here + "/vor.*")

			if err != nil {
				Debug("Walk error [%v]\n", err)
				return "", err
			}

			if len(matches) > 0 {
				Debug("Walker done, found match at '%s'", matches[0])
				return matches[0], nil
			}
		}
	}
	Debug("Walker done, no matches found")
	return "", nil
}

func CreateXDGConfig(config []byte) string {
	home, homeErr := homedir.Dir()
	if homeErr != nil {
		Exit("Vor encountered an error attempting to read from the filesystem")
	}

	if _, err := os.Stat(home + "/.config"); os.IsNotExist(err) {
		os.Mkdir(home + "/.config", 0755)
	}

	if _, err := os.Stat(home + "/.config/vor"); os.IsNotExist(err) {
		os.Mkdir(home + "/.config/vor", 0755)
	}
	if _, err := os.Stat(home + "/.config/vor/vor.yaml"); os.IsNotExist(err) {
		f, createErr := os.Create(home + "/.config/vor/vor.yaml")
		if createErr != nil {
			Exit("vor encountered an error attempting to write a default config to " + home + "/.config/vor/vor.yaml")
		}
		defer f.Close()
		_, writeErr := f.Write(config)
		if writeErr != nil {
			Exit("vor encountered an error attempting to write a default config to " + home + "/.config/vor/vor.yaml")
		}
		f.Sync()
		fmt.Println("Created a new default configuration at ~/.config/vor/vor.yaml")

	}
	return home + "/.config/vor"
}
