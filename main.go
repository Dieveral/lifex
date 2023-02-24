package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"lifex/commands"
	"lifex/entities"
	"os"
	"strconv"
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
			case "city":
				if name, ok := command.Args["name"]; ok {
					if id, err := addCity(name, db); err == nil {
						fmt.Printf("City '%s' was successfully added with Id=%d\n", name, id)
					} else {
						return err
					}
				} else {
					return fmt.Errorf("city name is not specified")
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

					id, err := strconv.ParseInt(val, 10, 0)
					if err != nil {
						return err
					}

					if err = showCompanyById(id, db); err != nil {
						return err
					}
				} else if name, ok := command.Args["name"]; ok {

					if err := showCompanyByName(name, db); err != nil {
						return err
					}
				} else {

					if err := showAllCompanies(db); err != nil {
						return err
					}
				}
			case "city":
				if val, ok := command.Args["id"]; ok {

					id, err := strconv.ParseInt(val, 10, 0)
					if err != nil {
						return err
					}

					if err = showCityById(id, db); err != nil {
						return err
					}
				} else if name, ok := command.Args["name"]; ok {

					if err := showCityByName(name, db); err != nil {
						return err
					}
				} else {

					if err := showAllCities(db); err != nil {
						return err
					}
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

func companyExists(name string, db *sql.DB) bool {

	row := db.QueryRow("select Count(1) from lifex.company where Name=?", name)

	var res int
	row.Scan(&res)

	if res == 1 {
		return true
	} else {
		return false
	}
}

// If operation is successful, returns id of the added company
func addCompany(name string, db *sql.DB) (int64, error) {

	if companyExists(name, db) {
		return 0, fmt.Errorf("company '%s' already exists", name)
	}

	result, err := db.Exec("insert into lifex.company (Name) values (?)", name)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func showCompanyById(id int64, db *sql.DB) error {

	row := db.QueryRow("select Id, Name from lifex.company where Id=?", id)
	company := entities.Company{}
	err := row.Scan(&company.Id, &company.Name)
	if err != nil {
		return err
	}

	idWidth := max(4, len(fmt.Sprint(company.Id)))
	nameWidth := max(6, strlen(company.Name))

	printCompanyHeader(idWidth, nameWidth)
	printCompanyInfo(company, idWidth, nameWidth)
	printCompanyFooter(idWidth, nameWidth)

	return nil
}

func showCompanyByName(name string, db *sql.DB) error {

	row := db.QueryRow("select Id, Name from lifex.company where Name=?", name)
	company := entities.Company{}
	err := row.Scan(&company.Id, &company.Name)
	if err != nil {
		return err
	}

	idWidth := max(4, len(fmt.Sprint(company.Id)))
	nameWidth := max(6, strlen(company.Name))

	printCompanyHeader(idWidth, nameWidth)
	printCompanyInfo(company, idWidth, nameWidth)
	printCompanyFooter(idWidth, nameWidth)

	return nil
}

func showAllCompanies(db *sql.DB) error {

	var companies []entities.Company

	rows, err := db.Query("select Id, Name from lifex.company")
	if err != nil {
		return err
	}
	defer rows.Close()

	idWidth := 4
	nameWidth := 6

	for rows.Next() {
		var company entities.Company
		rows.Scan(&company.Id, &company.Name)
		companies = append(companies, company)

		idWidth = max(idWidth, len(fmt.Sprint(company.Id)))
		nameWidth = max(nameWidth, strlen(company.Name))
	}

	printCompanyHeader(idWidth, nameWidth)
	for _, company := range companies {
		printCompanyInfo(company, idWidth, nameWidth)
	}
	printCompanyFooter(idWidth, nameWidth)

	return nil
}

func printCompanyInfo(company entities.Company, idWidth int, nameWidth int) {

	idw := len(fmt.Sprint(company.Id))
	nw := strlen(company.Name)

	fmt.Printf("|%s%d|%s%s|\n", strings.Repeat(" ", idWidth-idw), company.Id, company.Name, strings.Repeat(" ", nameWidth-nw))
}

func printCompanyHeader(idWidth int, nameWidth int) {

	idIndent := (idWidth - 2) / 2
	nameIndent := (nameWidth - 4) / 2

	fmt.Printf(" %s %s \n", strings.Repeat("_", idWidth), strings.Repeat("_", nameWidth))
	fmt.Printf("|%sID%s|%sName%s|\n", strings.Repeat(" ", idWidth-idIndent-2), strings.Repeat(" ", idIndent), strings.Repeat(" ", nameWidth-nameIndent-4), strings.Repeat(" ", nameIndent))
	fmt.Printf("|%s|%s|\n", strings.Repeat("_", idWidth), strings.Repeat("_", nameWidth))
}

func printCompanyFooter(idWidth int, nameWidth int) {
	fmt.Printf("|%s|%s|\n", strings.Repeat("_", idWidth), strings.Repeat("_", nameWidth))
}

func cityExists(name string, db *sql.DB) bool {

	row := db.QueryRow("select Count(1) from lifex.city where Name=?", name)

	var res int
	row.Scan(&res)

	if res == 1 {
		return true
	} else {
		return false
	}
}

func addCity(name string, db *sql.DB) (int64, error) {

	if cityExists(name, db) {
		return 0, fmt.Errorf("city '%s' already exists", name)
	}

	result, err := db.Exec("insert into lifex.city (Name) values (?)", name)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func showCityById(id int64, db *sql.DB) error {

	row := db.QueryRow("select Id, Name from lifex.city where Id=?", id)
	city := entities.City{}
	err := row.Scan(&city.Id, &city.Name)
	if err != nil {
		return err
	}

	idWidth := max(4, len(fmt.Sprint(city.Id)))
	nameWidth := max(6, strlen(city.Name))

	printCityHeader(idWidth, nameWidth)
	printCityInfo(city, idWidth, nameWidth)
	printCityFooter(idWidth, nameWidth)

	return nil
}

func showCityByName(name string, db *sql.DB) error {

	row := db.QueryRow("select Id, Name from lifex.city where Name=?", name)
	city := entities.City{}
	err := row.Scan(&city.Id, &city.Name)
	if err != nil {
		return err
	}

	idWidth := max(4, len(fmt.Sprint(city.Id)))
	nameWidth := max(6, strlen(city.Name))

	printCityHeader(idWidth, nameWidth)
	printCityInfo(city, idWidth, nameWidth)
	printCityFooter(idWidth, nameWidth)

	return nil
}

func showAllCities(db *sql.DB) error {

	var cities []entities.City

	rows, err := db.Query("select Id, Name from lifex.city")
	if err != nil {
		return err
	}
	defer rows.Close()

	idWidth := 4
	nameWidth := 6

	for rows.Next() {
		var city entities.City
		rows.Scan(&city.Id, &city.Name)
		cities = append(cities, city)

		idWidth = max(idWidth, len(fmt.Sprint(city.Id)))
		nameWidth = max(nameWidth, strlen(city.Name))
	}

	printCityHeader(idWidth, nameWidth)
	for _, city := range cities {
		printCityInfo(city, idWidth, nameWidth)
	}
	printCityFooter(idWidth, nameWidth)

	return nil
}

func printCityInfo(city entities.City, idWidth int, nameWidth int) {

	idw := len(fmt.Sprint(city.Id))
	nw := strlen(city.Name)

	fmt.Printf("|%s%d|%s%s|\n", strings.Repeat(" ", idWidth-idw), city.Id, city.Name, strings.Repeat(" ", nameWidth-nw))
}

func printCityHeader(idWidth int, nameWidth int) {

	idIndent := (idWidth - 2) / 2
	nameIndent := (nameWidth - 4) / 2

	fmt.Printf(" %s %s \n", strings.Repeat("_", idWidth), strings.Repeat("_", nameWidth))
	fmt.Printf("|%sID%s|%sName%s|\n", strings.Repeat(" ", idWidth-idIndent-2), strings.Repeat(" ", idIndent), strings.Repeat(" ", nameWidth-nameIndent-4), strings.Repeat(" ", nameIndent))
	fmt.Printf("|%s|%s|\n", strings.Repeat("_", idWidth), strings.Repeat("_", nameWidth))
}

func printCityFooter(idWidth int, nameWidth int) {
	fmt.Printf("|%s|%s|\n", strings.Repeat("_", idWidth), strings.Repeat("_", nameWidth))
}

func max(value int, newValue int) int {

	if newValue > value {
		return newValue
	} else {
		return value
	}
}

func strlen(text string) int {
	return len([]rune(text))
}
