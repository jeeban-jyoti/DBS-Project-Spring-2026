package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string

		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Table:", tableName)
	}

	fmt.Println("Connected to DB!")
}
