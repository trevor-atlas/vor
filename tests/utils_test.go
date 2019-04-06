package tests

import (
	"testing"
	"trevoratlas.com/vor/utils"
)

func TestKebab(t *testing.T) {
	result := utils.KebabCase("this is a string")
	expected := "this-is-a-string"
	if result != expected {
		t.Errorf("KebabCase was incorrect, got: %s, want: %s", result, expected)
	}
}

func TestLowerKebab(t *testing.T) {
	result := utils.LowerKebabCase("THIS IS A STRING")
	expected := "this-is-a-string"
	if result != expected {
		t.Errorf("LowerKebabCase was incorrect:\ngot: %s\nwant: %s", result, expected)
	}
}

func TestCaseInsensitiveContains(t *testing.T) {
	{
		result := utils.Contains("anything", "A")
		expected := true
		if result != expected {
			t.Errorf("contains was incorrect:\ngot: %t\nwant: %t", result, expected)
		}
	}
	{
		result := utils.Contains("anything", "X")
		expected := false
		if result != expected {
			t.Errorf("contains was incorrect:\ngot: %t\nwant: %t", result, expected)
		}
	}
}

func TestLeftPad(t *testing.T) {
	result := utils.LeftPad("GONNA PAD YA", "B ", 2)
	expected := "B B GONNA PAD YA"
	if result != expected {
		t.Errorf("LeftPad was incorrect:\ngot: %s\nwant: %s", result, expected)
	}
}
