package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://wizardcore:wizardcore@localhost:5432/wizardcore?sslmode=disable"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Read seed SQL file
	sqlBytes, err := os.ReadFile("internal/database/seeds/seed.sql")
	if err != nil {
		log.Fatal("Failed to read seed.sql:", err)
	}
	sql := string(sqlBytes)

	// Execute the SQL
	_, err = db.Exec(sql)
	if err != nil {
		log.Fatal("Failed to execute seed SQL:", err)
	}

	fmt.Println("Seed data inserted successfully.")
}