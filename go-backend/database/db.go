// package database

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"

// 	_ "github.com/lib/pq"
// )

// const (
// 	// Change to ur own db name, user and password
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "choy050823"
// 	dbname   = "travel_forum"
// )

// var DB *sql.DB

// func InitDB() {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)

// 	var err error
// 	DB, err = sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = DB.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Successfully connected to the database!")

// 	// Execute schema.sql on first run
// 	// executeSchema()
// }

// func executeSchema() {
// 	// Read the schema.sql file
// 	schema, err := os.ReadFile("./database/schema.sql")
// 	if err != nil {
// 		log.Fatal("Failed to read schema.sql:", err)
// 	}

// 	// Execute the schema.sql script
// 	_, err = DB.Exec(string(schema))
// 	if err != nil {
// 		log.Fatal("Failed to execute schema.sql:", err)
// 	}

// 	fmt.Println("Schema executed successfully!")
// }

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dbUser, dbPassword, dbName, dbHost, dbPort string) {
// func InitDB() {
	// local test
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory:", err)
	}

	envPath := filepath.Join(rootDir, "..", ".env") // Go up one level for .env
	err = godotenv.Load(envPath)
	if err != nil {
			log.Println("No .env file found, using environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL") // Get from environment variable
    if dbURL == "" {
        log.Fatal("DATABASE_URL environment variable not set")
    }

	// var err error
	DB, err = sql.Open("postgres", psqlInfo)
	// DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")
}
