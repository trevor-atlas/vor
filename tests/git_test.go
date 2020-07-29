package tests

import (
	"strings"
	"testing"
	"trevoratlas.com/vor/git"

	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("git.path", "/usr/local/bin/git")
}

func TestNewGit(t *testing.T) {
	gc := git.New()
	expected := viper.GetString("git.path")

	if expected != gc.Path {
		t.Errorf("client.New was incorrect:\ngot: %s\nwant: %s", gc.Path, expected)
	}
}
