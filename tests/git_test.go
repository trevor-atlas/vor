package tests

import (
	"github.com/trevor-atlas/vor/git"
	"testing"
	"github.com/spf13/viper"
	"strings"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("global.git.path", "/usr/local/bin/git")
}

func TestNewGit(t *testing.T) {
	gc := git.New()
	expected := "/usr/local/bin/git"

	if expected != gc.Path {
		t.Errorf("client.New was incorrect:\ngot: %s\nwant: %s", gc.Path, expected)
	}
}
