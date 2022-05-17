package main

import (
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type Settings struct {
	Id                  int
	Server              string
	Key                 string
	ValidateCertificate int
	Date                time.Time
	DateFormat          string
}

func createSettingsTable() {
	db := GetConnection()

	createSettingsTableSQL := `CREATE TABLE settings (
		"Id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Server" TEXT NOT NULL,
		"Key" TEXT NOT NULL,
		"Date" DATETIME NOT NULL,
		"ValidateCertificate" INTEGER
	  );` // SQL Statement for Create Table

	log.Println("Create settings table...")
	statement, err := db.Prepare(createSettingsTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}

	defer statement.Close()

	statement.Exec() // Execute SQL Statements
	log.Println("settings table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertSettings(settings Settings) {
	db := GetConnection()

	log.Println("Inserting settings record ...")
	insertSettingsSQL := `INSERT INTO settings(server, key, date, validatecertificate) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertSettingsSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln("prepare: " + err.Error())
	}

	defer statement.Close()
	_, err = statement.Exec(settings.Server, settings.Key, settings.Date, settings.ValidateCertificate)
	if err != nil {
		log.Fatalln("exec: " + err.Error())
	}
}

func getSettings() Settings {
	db := GetConnection()
	row, err := db.Query("SELECT * FROM settings")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	var settings Settings
	for row.Next() { // Iterate and fetch the records from result cursor
		s := Settings{}
		if err := row.Scan(&s.Id, &s.Server, &s.Key, &s.Date, &s.ValidateCertificate); err != nil {
			log.Fatal(err)
		}

		s.DateFormat = s.Date.Format("Mon Jan 2 15:04:05")
		settings = s
		break
	}

	return settings
}

func updateSettings(settings Settings) {
	db := GetConnection()

	log.Println("Updating settings record ...")
	sqlUpdateSettings := "UPDATE settings SET server = ?, key = ?, date = ?, validatecertificate = ? WHERE id = ?"
	statement, err := db.Prepare(sqlUpdateSettings)
	if err != nil {
		log.Fatalln("prepare: " + err.Error())
	}

	defer statement.Close()

	_, err = statement.Exec(settings.Server, settings.Key, settings.Date, settings.ValidateCertificate, settings.Id)
	if err != nil {
		log.Fatalln("exec: " + err.Error())
	}
}
