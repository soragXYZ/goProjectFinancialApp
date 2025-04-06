package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"financialApp/config"
)

func CreatePermanentUserToken(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/permanentUserToken/" {
		http.NotFound(w, r)
		return
	}

	// check if one permanent user token already exists in DB or not
	var entries uint
	row := config.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM authToken)")
	if err := row.Scan(&entries); err != nil {
		log.Fatal("SELECT EXISTS (SELECT 1 FROM authToken): %v", err)
		http.Error(w, "Error in SELECT EXISTS (SELECT 1 FROM authToken)", http.StatusInternalServerError)
		return
	}

	if entries != 0 {
		http.Error(w, "Permanent user token already exists", http.StatusConflict)
		return
	}

	initToken := AuthTokenInitRequest{
		Client_id:     os.Getenv("CLIENT_ID"),
		Client_secret: os.Getenv("CLIENT_SECRET"),
	}

	const url string = "https://testfinary-sandbox.biapi.pro/2.0/auth/init"
	jsonBody, err := json.Marshal(initToken)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
	}

	var authToken AuthToken

	err = json.NewDecoder(resp.Body).Decode(&authToken)
	if err != nil {
		log.Fatal(err)
	}

	_, err = config.DB.Exec(
		"INSERT INTO authToken (auth_token, id_user) VALUES (?, ?)",
		authToken.Auth_token, authToken.Id_user,
	)
	if err != nil {
		log.Fatal("INSERT INTO authToken: %v", err)
	}

	jsonBody, err = json.Marshal(authToken)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBody)

}

func CreateTemporaryUserToken(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/temporaryUserToken/" {
		http.NotFound(w, r)
		return
	}

	// check if one permanent user token already exists in DB or not
	var permanentUserToken string
	row := config.DB.QueryRow("SELECT auth_token FROM authToken LIMIT 1")
	if err := row.Scan(&permanentUserToken); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Permanent user token does not exists", http.StatusNotFound)
			return
		}

		log.Fatal("SELECT auth_token FROM authToken LIMIT 1: %v", err)
		http.Error(w, "Error in SELECT auth_token FROM authToken LIMIT 1", http.StatusInternalServerError)
		return
	}

	const url string = "https://testfinary-sandbox.biapi.pro/2.0/auth/token/code"
	var bearer string = "Bearer " + permanentUserToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", bearer)
	fmt.Println(bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
	}

	var authCode AuthCode

	err = json.NewDecoder(resp.Body).Decode(&authCode)
	if err != nil {
		log.Fatal(err)
	}

	// Clean old temporary token from DB and insert the new one
	_, err = config.DB.Exec("DELETE from authCode")
	if err != nil {
		log.Fatal("DELETE authCode: %v", err)
	}

	_, err = config.DB.Exec(
		"INSERT INTO authCode (auth_code, expires_in) VALUES (?, ?)",
		authCode.Auth_Code, authCode.Expires_in,
	)
	if err != nil {
		log.Fatal("INSERT INTO authCode: %v", err)
	}

	jsonBody, err := json.Marshal(authCode)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBody)
}

func GetTemporaryUserToken(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/temporaryUserToken/" {
		http.NotFound(w, r)
		return
	}

	var authCode AuthCode

	row := config.DB.QueryRow("SELECT * FROM authCode LIMIT 1")
	if err := row.Scan(&authCode.Auth_Code, &authCode.Expires_in); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Token does not exist", http.StatusNotFound)
			return
		}

		log.Fatal("SELECT * FROM authCode LIMIT 1: %v", err)
		http.Error(w, "Error in SELECT * FROM authCode LIMIT 1", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(authCode)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonBody)
}

func GetPermanentUserToken(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/permanentUserToken/" {
		http.NotFound(w, r)
		return
	}

	var authToken AuthToken

	row := config.DB.QueryRow("SELECT * FROM authToken LIMIT 1")
	if err := row.Scan(&authToken.Auth_token, &authToken.Id_user); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Token does not exist", http.StatusNotFound)
			return
		}

		log.Fatal("SELECT * FROM authToken LIMIT 1: %v", err)
		http.Error(w, "Error in SELECT * FROM authToken LIMIT 1", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(authToken)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonBody)
}

func DeletePermanentUserToken(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/permanentUserToken/" {
		http.NotFound(w, r)
		return
	}

	_, err := config.DB.Exec("DELETE from authToken")
	if err != nil {
		log.Fatal("DELETE authToken: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
