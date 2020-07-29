package tests

import (
	"strings"
	"testing"


	"github.com/spf13/viper"
)

const TEST_DATA_PATH = "../test-data/"

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func TestGetStringEnvExists(t *testing.T) {
	result := viper.GetString("PATH")
	expected := 10
	if expected > len(result) {
		t.Errorf("GetStringEnv was incorrect:\ngot: '%d'\nwant: '%d'", len(result), expected)
	}
}

func TestGetStringFallsBackToGlobal(t *testing.T) {
	viper.SetDefault("testthing", "this should match")
	result := viper.GetString("testthing")
	expected := "this should match"
	if expected != result {
		t.Errorf("GetStringEnv was incorrect:\ngot: '%s'\nwant: '%s'", result, expected)
	}
}

func TestGetStringEnvDoesNotExist(t *testing.T) {
	result := viper.GetString("SOMETHING_THAT_DOES_NOT_EXIST")
	expected := ""
	if expected != result {
		t.Errorf("GetStringEnv was incorrect:\ngot: '%s'\nwant: '%s'", result, expected)
	}
}

func TestExistsTruthy(t *testing.T) {
	result, _ := system.Exists(TEST_DATA_PATH + "testdata.keep")
	expected := true
	if result != expected {
		t.Errorf("Exists was incorrect:\ngot: %t\nwant: %t", result, expected)
	}
}

func TestExistsFalsy(t *testing.T) {
	result, _ := system.Exists(TEST_DATA_PATH + "notarealfile.json")
	expected := false
	if result != expected {
		t.Errorf("Exists was incorrect:\ngot: %t\nwant: %t", result, expected)
	}
}
