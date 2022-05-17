package main

import (
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type Module struct {
	Id         int
	Name       string
	Image      string
	Enable     int
	Date       time.Time
	DateFormat string
}

func createModulesTable() {
	db := GetConnection()

	createModulesTableSQL := `CREATE TABLE modules (
		"Id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Name" TEXT NOT NULL,
		"Image" TEXT NOT NULL,
		"Enable" INTEGER,
		"Date" DATETIME NOT NULL
	  );` // SQL Statement for Create Table

	log.Println("Create modules table...")
	statement, err := db.Prepare(createModulesTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}

	defer statement.Close()

	statement.Exec() // Execute SQL Statements
	log.Println("modules table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertModule(module Module) {
	db := GetConnection()

	log.Println("Inserting module record ...")
	insertModuleSQL := `INSERT INTO modules(name, image, enable, date) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertModuleSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln("prepare: " + err.Error())
	}

	defer statement.Close()
	_, err = statement.Exec(module.Name, module.Image, module.Enable, module.Date)
	if err != nil {
		log.Fatalln("exec: " + err.Error())
	}
}

func getModules() []Module {
	db := GetConnection()
	row, err := db.Query("SELECT * FROM modules")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	var modules []Module
	for row.Next() { // Iterate and fetch the records from result cursor
		m := Module{}
		if err := row.Scan(&m.Id, &m.Name, &m.Image, &m.Enable, &m.Date); err != nil {
			log.Fatal(err)
		}

		m.DateFormat = m.Date.Format("Mon Jan 2 15:04:05")
		modules = append(modules, m)
	}

	return modules
}

func updateModule(module Module) {
	db := GetConnection()

	log.Println("Updating modules record ...")
	sqlUpdateModule := "UPDATE modules SET name = ?, image = ?, enable = ?, date = ? WHERE id = ?"
	statement, err := db.Prepare(sqlUpdateModule)
	if err != nil {
		log.Fatalln("prepare: " + err.Error())
	}

	defer statement.Close()

	_, err = statement.Exec(module.Name, module.Image, module.Enable, module.Date, module.Id)
	if err != nil {
		log.Fatalln("exec: " + err.Error())
	}
}
