package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"financialApp/config"
)

func CreatePermanentUserToken(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/createPermanentUserToken/" {
		http.NotFound(w, r)
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
		"INSERT INTO authToken (auth_token, token_type, id_user, expires_in) VALUES (?, ?, ?, ?)",
		authToken.Auth_token, authToken.Token_type, authToken.Id_user, authToken.Expires_in,
	)
	if err != nil {
		log.Fatal("INSERT INTO authToken: %v", err)
	}

	return

}

func GetPermanentUserToken(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/getPermanentUserToken/" {
		http.NotFound(w, r)
		return
	}

	var authToken AuthToken

	row := config.DB.QueryRow("SELECT * FROM authToken LIMIT 1")
	// check syntax why err at the end with ;
	if err := row.Scan(&authToken.Auth_token, &authToken.Token_type, &authToken.Id_user, &authToken.Expires_in); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Token does not exist", http.StatusNotFound)
			return
		}

		log.Fatal("SELECT TOP 1 * FROM authToken: %v", err)
		http.Error(w, "Error in select top 1", http.StatusInternalServerError)
	}

	jsonBody, err := json.Marshal(authToken)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonBody)
}
