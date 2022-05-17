package main

import (
	"log"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type LogObj struct {
	Id         int
	LogType    string
	Desc       string
	Date       time.Time
	DateFormat string
}

func createLogTable() {
	db := GetConnection()

	createLogsTableSQL := `CREATE TABLE logs (
		"Id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"LogType" TEXT NOT NULL,
		"Desc" TEXT NOT NULL,
		"Date" DATETIME NOT NULL
	  );` // SQL Statement for Create Table

	log.Println("Create logs table...")
	statement, err := db.Prepare(createLogsTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}

	defer statement.Close()

	statement.Exec() // Execute SQL Statements
	log.Println("logs table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertLog(logObj LogObj) {
	db := GetConnection()

	log.Println("Inserting log record ...")
	insertLogSQL := `INSERT INTO logs(logType, desc, date) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertLogSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln("prepare: " + err.Error())
	}

	defer statement.Close()
	_, err = statement.Exec(logObj.LogType, logObj.Desc, logObj.Date)
	if err != nil {
		log.Fatalln("exec: " + err.Error())
	}
}

func getLogs() []LogObj {
	db := GetConnection()
	row, err := db.Query("SELECT * FROM logs")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	var logs []LogObj
	for row.Next() { // Iterate and fetch the records from result cursor
		l := LogObj{}
		if err := row.Scan(&l.Id, &l.LogType, &l.Desc, &l.Date); err != nil {
			log.Fatal(err)
		}

		l.DateFormat = l.Date.Format("Mon Jan 2 15:04:05")
		logs = append(logs, l)
	}
	sort.Slice(logs[:], func(i, j int) bool {
		return logs[i].Date.After(logs[j].Date)
	})
	return logs
}
