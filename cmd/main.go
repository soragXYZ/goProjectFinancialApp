package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"financialApp/api/router"
	"financialApp/config"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Check env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var s = []string{"CLIENT_ID", "CLIENT_SECRET", "DBNAME", "DBUSER", "DBPASS"}
	for _, env := range s {
		_, envExists := os.LookupEnv(env)
		if !envExists {
			log.Fatal(env, " does not exist")
		}
	}

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: os.Getenv("DBNAME"),
	}

	// Get a database handle.
	config.DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := config.DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to DB...")

	router.New()
}
