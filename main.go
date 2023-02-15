package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"lifex/commands"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const dbDriver string = "mysql"
const dbName string = "lifex"
const user string = "Diever"
const password string = "D7BKJUcTsTqNA!A323#FofbK"

func main() {

	fmt.Println("Lifex - is an application for showing your's everyday expenses statistic.")

	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@/%s", user, password, dbName))
	if err != nil {
		fmt.Println("ERROR: Database connaction failed.")
		fmt.Println(err)
		return
	} else {
		fmt.Println("Database connection established.")
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	for command := readCommand(reader); strings.ToLower(command.Name) != "exit"; command = readCommand(reader) {
		err := executeCommand(command, db)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}

func executeCommand(command *commands.Command, db *sql.DB) error {

	switch strings.ToLower(command.Name) {
	case "check":
		fmt.Println("...checking database connection")
		if err := db.Ping(); err != nil {
			fmt.Println("Database connection failed.")
			return err
		} else {
			fmt.Println("SUCCESS: Database is connected.")
		}
	case "help":
		fmt.Println("...list of available commands:")
		fmt.Println("check - checks database connection")
		fmt.Println("help - shows help")
		fmt.Println("exit - closes application")
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}

	return nil
}

func readCommand(reader *bufio.Reader) *commands.Command {

	fmt.Println()
	fmt.Print("> ")

	line, _ := reader.ReadString('\n')
	return parseCommand(line)
}

// Correct command has template: Name Target Argument1 = Value 1 Argument2 = Value 2
// Command can have no target and/or no aarguments
func parseCommand(text string) *commands.Command {

	var command = commands.NewCommand()

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
