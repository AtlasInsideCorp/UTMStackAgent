package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func createDB() {
	//os.Remove("UTMStack.db")
	_, err := os.Stat("UTMStack.db")
	if os.IsNotExist(err) {
		log.Println("Creating UTMStack.db...")
		file, err := os.Create("UTMStack.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("UTMStack.db created")

		createUserTable()
		createLogTable()
		createSettingsTable()
		createModulesTable()

		insertModule(Module{Id: 0, Name: "test 1", Image: "assets/media/logos/UTMStack.svg", Enable: 1, Date: time.Now()})
		insertModule(Module{Id: 0, Name: "test 2", Image: "assets/media/logos/UTMStack.svg", Enable: 1, Date: time.Now()})
		insertModule(Module{Id: 0, Name: "test 3", Image: "assets/media/logos/UTMStack.svg", Enable: 1, Date: time.Now()})
		insertModule(Module{Id: 0, Name: "test 4", Image: "assets/media/logos/UTMStack.svg", Enable: 1, Date: time.Now()})
		insertModule(Module{Id: 0, Name: "test 5", Image: "assets/media/logos/UTMStack.svg", Enable: 1, Date: time.Now()})
		insertModule(Module{Id: 0, Name: "test 6", Image: "assets/media/logos/UTMStack.svg", Enable: 1, Date: time.Now()})
	}
	return
}

func GetConnection() *sql.DB {
	if db != nil {
		return db
	}

	var err error
	// Conexión a la base de datos
	db, err = sql.Open("sqlite3", "UTMStack.db") // Open the created SQLite File
	checkErr(err)

	return db
}

func checkErr(err error, args ...string) {
	if err != nil {
		log.Println("Error")
		log.Println("%q: %s", err, args)
	}
}
