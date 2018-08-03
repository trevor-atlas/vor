package tests

import (
	"github.com/spf13/viper"
	"github.com/trevor-atlas/vor/system"
	"strings"
	"testing"
)

const TEST_DATA_PATH = "../test-data/"

func TestGetStringEnv(t *testing.T) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	env := system.NewENVGetter()
	result := env.String("PATH")
	expected := 10
	if expected > len(result) {
		t.Errorf("GetStringEnv was incorrect:\ngot: %d\nwant: %d", len(result), expected)
	}
}

func TestExists(t *testing.T) {
	sys := system.NewOSHandler()
	{
		result, _ := sys.Exists(TEST_DATA_PATH + "testdata.keep")
		expected := true
		if result != expected {
			t.Errorf("Exists was incorrect:\ngot: %t\nwant: %t", result, expected)
		}
	}
	{
		result, _ := sys.Exists(TEST_DATA_PATH + "notarealfile.json")
		expected := false
		if result != expected {
			t.Errorf("Exists was incorrect:\ngot: %t\nwant: %t", result, expected)
		}
	}
}
