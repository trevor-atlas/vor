package tests

import (
	"github.com/trevor-atlas/vor/git"
	"testing"
)

var sadEnv = &sadMockGetter{}
var happyEnv = &happyMockGetter{}
var sadOSHandler = &sadMockOSHandler{}
var happyOSHandler = &happyMockOSHandler{}

func TestNewGitHappyPath(t *testing.T) {
	g := git.New(happyOSHandler, happyEnv)
	expected := happyOSHandler

	if expected != g.Sys {
		t.Errorf("client.New was incorrect:\ngot: %s\nwant: %s", g.Sys, expected)
	}
}

func TestNewGitHappyPath2(t *testing.T) {
	g := git.New(happyOSHandler, happyEnv)
	expected := "/usr/local/bin/git"
	if expected != g.Path {
		t.Errorf("client.New was incorrect:\ngot: %s\nwant: %s", g.Path, expected)
	}
}

func TestNewGitSadPath(t *testing.T) {
	g := git.New(sadOSHandler, sadEnv)
	expected := "/usr/local/bin/git"
	if expected == g.Path {
		t.Errorf("client.New should fail! got: %s\nwant: %s", g.Path, expected)
	}
}

func TestNewGitSadPath2(t *testing.T) {
	g := git.New(sadOSHandler, sadEnv)
	expected := "/usr/local/bin/git"
	if expected == g.Path {
		t.Errorf("client.New should fail! got: %s\nwant: %s", g.Path, expected)
	}
}
