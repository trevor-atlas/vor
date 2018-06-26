package utils

import (
	"strings"
	"github.com/spf13/viper"
	"testing"
)

const TEST_DATA_PATH = "../test-data/"

func TestKebab(t *testing.T) {
	result := KebabCase("this is a string")
	expected := "this-is-a-string"
	if result != expected {
		t.Errorf("KebabCase was incorrect, got: %s, want: %s", result, expected)
	}
}

func TestLowerKebab(t *testing.T) {
	result := LowerKebabCase("THIS IS A STRING")
	expected := "this-is-a-string"
	if result != expected {
		t.Errorf("LowerKebabCase was incorrect:\ngot: %s\nwant: %s", result, expected)
	}
}

func TestCaseInsensitiveContains(t *testing.T) {
	{
		result := CaseInsensitiveContains("anything", "A")
		expected := true
		if result != expected {
			t.Errorf("CaseInsensitiveContains was incorrect:\ngot: %t\nwant: %t", result, expected)
		}
	}
	{
		result := CaseInsensitiveContains("anything", "X")
		expected := false
		if result != expected {
			t.Errorf("CaseInsensitiveContains was incorrect:\ngot: %t\nwant: %t", result, expected)
		}
	}
}

func TestExists(t *testing.T) {
	{
		result, _ := Exists(TEST_DATA_PATH + "testdata.keep")
		expected := true
		if result != expected {
			t.Errorf("Exists was incorrect:\ngot: %t\nwant: %t", result, expected)
		}
	}
	{
		result, _ := Exists(TEST_DATA_PATH + "notarealfile.json")
		expected := false
		if result != expected {
			t.Errorf("Exists was incorrect:\ngot: %t\nwant: %t", result, expected)
		}
	}
}

func TestLeftPad(t *testing.T) {
	result := LeftPad("GONNA PAD YA", "B ", 2)
	expected := "B B GONNA PAD YA"
	if result != expected {
		t.Errorf("LeftPad was incorrect:\ngot: %s\nwant: %s", result, expected)
	}
}

func TestRightPad(t *testing.T) {
	result := RightPad("GONNA PAD YA", " B", 2)
	expected := "GONNA PAD YA B B"
	if result != expected {
		t.Errorf("RightPad was incorrect:\ngot: %s\nwant: %s", result, expected)
	}
}

func TestGetStringEnv(t *testing.T) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	result := GetStringEnv("PATH")
	expected := 10
	if expected > len(result) {
		t.Errorf("GetStringEnv was incorrect:\ngot: %d\nwant: %d", len(result), expected)
	}
}