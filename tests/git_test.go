package tests

import (
	"github.com/trevor-atlas/vor/git"
	"testing"
	"github.com/spf13/viper"
	"strings"
	"github.com/trevor-atlas/vor/system"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("global.git.path", "/usr/local/bin/git")
}

func TestNewGit(t *testing.T) {
	gc := git.New()
	expected := system.GetString("git.path")

	if expected != gc.Path {
		t.Errorf("client.New was incorrect:\ngot: %s\nwant: %s", gc.Path, expected)
	}
}
