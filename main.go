package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"financialApp/api/resource/auth"
	"fmt"
	"log"
	"net/http"
	"os"

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
	// var err error
	var db *sql.DB
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	initToken := auth.AuthTokenInitRequest{
		Client_id:     os.Getenv("CLIENT_ID"),
		Client_secret: os.Getenv("CLIENT_SECRET"),
	}

	const url string = "https://testfinary-sandbox.biapi.pro/2.0/auth/init"
	jsonBody, err := json.Marshal(initToken)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	var authToken auth.AuthToken
	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(&authToken)
		if err != nil {
			// return HTTP 400 bad request
		}
		fmt.Printf("%s\n", authToken.Auth_token)
		fmt.Printf("%s\n", authToken.Token_type)
		fmt.Printf("%d\n", authToken.Id_user)
		fmt.Printf("%d\n", authToken.Expires_in)
	}

	result, err := db.Exec("INSERT INTO authToken (auth_token, token_type, id_user, expires_in) VALUES (?, ?, ?, ?)", authToken.Auth_token, authToken.Token_type, authToken.Id_user, authToken.Expires_in)
	if err != nil {
		log.Fatal("INSERT INTO authToken: %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		log.Fatal("INSERT INTO authToken: %v", err)
	}
	fmt.Printf("Inserted ID: %d\n", id)
}
