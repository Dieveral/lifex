package commands

import (
	"fmt"
	"testing"
)

func testCommandParse(text string, expected *Command, t *testing.T) {

	command := Parse(text)

	if command.Name != expected.Name {
		t.Fatalf("Command Name is not correct.\nExpected: '%s'\nGot: '%s'", expected.Name, command.Name)
	} else if command.Target != expected.Target {
		t.Fatalf("Command Target is not correct.\nExpected: '%s'\nGot: '%s'", expected.Target, command.Target)
	} else if len(command.Args) != len(expected.Args) {
		t.Fatalf("Command Args count is not correct.\nExpected: '%d'\nGot: '%d'", len(expected.Args), len(command.Args))
	}
}

// Test parsing command with name only
func TestCommandParseNameOnly(t *testing.T) {
	name := "help"
	text := fmt.Sprintf("%s", name)

	expected := NewCommand()
	expected.Name = name

	testCommandParse(text, expected, t)
}

// Test parsing command with name only with carret move
func TestCommandParseNameOnlyWithCarretMove(t *testing.T) {
	name := "help"
	text := fmt.Sprintf("%s\r\n", name)

	expected := NewCommand()
	expected.Name = name

	testCommandParse(text, expected, t)
}

// Test parsing command with name only
func TestCommandParseNameAndTargetOnly(t *testing.T) {
	name := "show"
	target := "company"
	text := fmt.Sprintf("%s %s", name, target)

	expected := NewCommand()
	expected.Name = "show"
	expected.Target = "company"

	testCommandParse(text, expected, t)
}
