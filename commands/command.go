package commands

import "fmt"

type Command struct {
	Name   string
	Target string
	Args   map[string]string
}

func NewCommand() *Command {
	command := Command{Name: "", Args: make(map[string]string)}
	return &command
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
