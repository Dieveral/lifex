package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Lifex - is an application for showing yours everyday expenses statistic.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for command := readCommand(reader); command != "exit"; command = readCommand(reader) {
		err := executeCommand(command)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}

func executeCommand(command string) error {
	switch command {
	case "help":
		fmt.Println("List of available commands:")
		fmt.Println("help - shows help")
		fmt.Println("exit - closes application")
		fmt.Println()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println()
	}
	return nil
}

func readCommand(reader *bufio.Reader) string {
	fmt.Print("> ")
	line, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(line))
}
