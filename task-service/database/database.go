package database

import (
	"database/sql"
	"log"
    "os"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
    connStr := os.Getenv("DATABASE_URL")
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Unable to reach the database: %v", err)
    }

    log.Println("Connected to the PostgreSQL database successfully from task service!" + connStr)
}

func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database connection is not initialized. Call InitDB() first.")
	}
	return db
}

func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully.")
		}
	}
}