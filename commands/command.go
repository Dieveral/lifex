package commands

import (
	"fmt"
	"strings"
)

type Command struct {
	Name   string
	Target string
	Args   map[string]string
}

func NewCommand() *Command {
	command := Command{Name: "", Args: make(map[string]string)}
	return &command
}

// Correct command has template: Name Target Argument1 = Value 1 Argument2 = Value 2
// Command can have no target and/or no aarguments
func Parse(text string) *Command {

	var command = NewCommand()

	text = strings.TrimSpace(text)

	var index int
	var value string
	if index = strings.Index(text, "="); index >= 0 {
		value = strings.TrimSpace(text[:index])
		text = strings.TrimSpace(text[index+1:])
	} else {
		value = text
		text = ""
	}

	vals := strings.Split(value, " ")

	var argName string
	switch len(vals) {
	case 0:
		return command
	case 1:
		if index < 0 {
			command.Name = vals[0]
		} // if it is not command name but argument name, return empty command
		return command
	case 2:
		command.Name = vals[0]
		if index < 0 {
			command.Target = vals[1]
			return command
		}
		argName = vals[1]
	default:
		command.Name = vals[0]
		command.Target = vals[1]
		// skip not argument name values
		if index < 0 {
			return command
		}
		argName = vals[len(vals)-1]
	}

	for index = strings.Index(text, "="); index >= 0; index = strings.Index(text, "=") {

		value = strings.TrimSpace(text[:index])

		spaceIndex := strings.LastIndex(value, " ")
		if spaceIndex < 0 {
			command.AddArgument(argName, "")
			argName = value
		} else {
			command.AddArgument(argName, value[:spaceIndex])
			argName = value[spaceIndex+1:]
		}

		text = strings.TrimSpace(text[index+1:])
	}

	command.AddArgument(argName, text) // Add last argument

	return command
}

func (command Command) String() string {

	if command.Name == "" {
		return ""
	} else {

		result := command.Name

		if command.Target != "" {
			result = fmt.Sprintf("%s %s", result, command.Target)
		}

		for key, value := range command.Args {
			result = fmt.Sprintf("%s %s=\"%s\"", result, key, value)
		}

		return result
	}
}

func (command *Command) AddArgument(name string, value string) {
	command.Args[name] = value
}
