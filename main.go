package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/router"
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

	database.InitDB(cfg)

	//dummy check for using the variable
	if err := database.DB.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to DB!")

	http.HandleFunc("/api/v1/", router.Route)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
