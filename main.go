package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"lifex/commands"
	"lifex/entities"
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
	case "add":
		if command.Target == "" {
			return fmt.Errorf("command target is not specified")
		} else {
			switch strings.ToLower(command.Target) {
			case "company":
				if name, ok := command.Args["name"]; ok {
					if id, err := addCompany(name, db); err == nil {
						fmt.Printf("Company '%s' was successfully added with Id=%d\n", name, id)
					} else {
						return err
					}
				} else {
					return fmt.Errorf("company name is not specified")
				}
			}
		}
	case "check":
		fmt.Println("...checking database connection")
		if err := db.Ping(); err != nil {
			fmt.Println("Database connection failed.")
			return err
		} else {
			fmt.Println("SUCCESS: Database is connected.")
		}
	case "getid":
		if command.Target == "" {
			return fmt.Errorf("command target is not specified")
		} else {
			switch strings.ToLower(command.Target) {
			case "company":
				if name, ok := command.Args["name"]; ok {
					if id, err := getCompanyId(name, db); err == nil {
						fmt.Printf("Company '%s' has Id=%d\n", name, id)
					} else {
						return err
					}
				} else {
					return fmt.Errorf("company name is not specified")
				}
			}
		}
	case "help":
		switch strings.ToLower(command.Target) {
		case "":
			fmt.Println("...list of available commands:")
			fmt.Println("add [table] <column1>=<value1> <column2>=<value2> ... - adds record to [table] with <column>=<value> pairs")
			fmt.Println("check - checks database connection")
			fmt.Println("exit - closes application")
			fmt.Println("help <command> - shows help for <command>. If command is not specified, total help is shown")
		case "add":

		}
	case "show":
		if command.Target == "" {
			return fmt.Errorf("command target is not specified")
		} else {
			switch strings.ToLower(command.Target) {
			case "company":
				if val, ok := command.Args["id"]; ok {
					atoi
					if id, err := getCompanyId(name, db); err == nil {
						fmt.Printf("Company '%s' has Id=%d\n", name, id)
					} else {
						return err
					}
				} else {
					return fmt.Errorf("company name is not specified")
				}
			}
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}

	return nil
}

func readCommand(reader *bufio.Reader) *commands.Command {

	fmt.Println()
	fmt.Print("> ")

	line, _ := reader.ReadString('\n')
	return commands.Parse(line)
}

func getCompanyId(name string, db *sql.DB) (int64, error) {

	row := db.QueryRow("select Id from lifex.company where Name=?", name)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

// If operation is successful, returns id of the added company
func addCompany(name string, db *sql.DB) (int64, error) {

	result, err := db.Exec("insert into lifex.company (Name) values (?)", name)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func showCompanyById(id int64, db *sql.DB) {

	fmt.Println(" ________ ________ ")
	fmt.Println("|   ID   |  Name  |")
	fmt.Println("|________|________|")

	row := db.QueryRow("select Id, Name from lifex.company where Id=?", id)
	company := entities.Company{}
	row.Scan(&company.Id, &company.Name)

	fmt.Printf("|%8d|%8s|\n", company.Id, company.Name)
	fmt.Println("|________|________|")
}
