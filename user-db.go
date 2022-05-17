package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type User struct {
	id       int
	ip       string
	password string
}

func createUserTable() {
	db := GetConnection()

	createUserTableSQL := `CREATE TABLE user (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"ip" TEXT,
		"password" TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create user table...")
	statement, err := db.Prepare(createUserTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}

	defer statement.Close()

	statement.Exec() // Execute SQL Statements
	log.Println("user table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertUser(ip string, password string) {
	db := GetConnection()

	log.Println("Inserting user record ...")
	insertUserSQL := `INSERT INTO user(ip, password) VALUES (?, ?)`
	statement, err := db.Prepare(insertUserSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln("prepare: " + err.Error())
	}

	defer statement.Close()
	_, err = statement.Exec(ip, password)
	if err != nil {
		log.Fatalln("exec: " + err.Error())
	}
}

func getUser(ip string) User {
	db := GetConnection()
	row, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	var user User
	for row.Next() { // Iterate and fetch the records from result cursor
		if err := row.Scan(&user.id, &user.ip, &user.password); err != nil {
			log.Fatal(err)
		}
	}

	return user
}

func userChangePassword(ip, password string) {
	db := GetConnection()

	log.Println("Updating user record ...")
	sqlUpdateUser := "UPDATE user SET password = ? WHERE ip = ?"
	statement, err := db.Prepare(sqlUpdateUser)
	if err != nil {
		log.Fatalln("prepare: " + err.Error())
	}

	defer statement.Close()

	_, err = statement.Exec(password, ip)
	if err != nil {
		log.Fatalln("exec: " + err.Error())
	}
}
