package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	// Change to ur own db name, user and password
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "choy050823"
	dbname   = "travel_forum"
)

var DB *sql.DB

func InitDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")

	// Execute schema.sql on first run
	// executeSchema()
}

func executeSchema() {
	// Read the schema.sql file
	schema, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		log.Fatal("Failed to read schema.sql:", err)
	}

	// Execute the schema.sql script
	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal("Failed to execute schema.sql:", err)
	}

	fmt.Println("Schema executed successfully!")
}