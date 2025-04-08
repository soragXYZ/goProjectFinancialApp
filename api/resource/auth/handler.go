package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"financialApp/config"
)

func CreatePermanentUserToken(w http.ResponseWriter, r *http.Request) {

	// Check if one permanent user token already exists in DB or not
	var entries uint
	var query string = "SELECT EXISTS (SELECT 1 FROM authToken)"
	row := config.DB.QueryRow(query)
	if err := row.Scan(&entries); err != nil {
		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, "", http.StatusInternalServerError)
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

	// Get a permanent user token from Powens API and store it in DB
	const url string = "https://testfinary-sandbox.biapi.pro/2.0/auth/init"
	jsonBody, err := json.Marshal(initToken)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot marshal initToken")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		config.Logger.Error().Err(err).Msg("Error in post request")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		config.Logger.Error().Msg(resp.Status)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var authToken AuthToken

	err = json.NewDecoder(resp.Body).Decode(&authToken)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot decode resp.Body")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	query = "INSERT INTO authToken (auth_token, id_user) VALUES (?, ?)"
	_, err = config.DB.Exec(query, authToken.Auth_token, authToken.Id_user)
	if err != nil {
		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	jsonBody, err = json.Marshal(authToken)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot marshal authToken")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBody)
}

func CreateTemporaryUserToken(w http.ResponseWriter, r *http.Request) {

	// Check if one permanent user token already exists in DB or not
	var permanentUserToken string
	var query string = "SELECT auth_token FROM authToken LIMIT 1"
	row := config.DB.QueryRow(query)
	if err := row.Scan(&permanentUserToken); err != nil {

		if err == sql.ErrNoRows {
			http.Error(w, "Permanent user token does not exists", http.StatusNotFound)
			return
		}

		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Get a temporary user token from Powens API and store it in DB
	const url string = "https://testfinary-sandbox.biapi.pro/2.0/auth/token/code"
	var bearer string = "Bearer " + permanentUserToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot create new get request")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot execute request")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		config.Logger.Error().Msgf("Powens did not answer correctly: %s", resp.Status)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var authCode AuthCode

	err = json.NewDecoder(resp.Body).Decode(&authCode)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot decode resp.Body")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Delete old temporary token from DB and insert the new one
	query = "DELETE from authCode"
	_, err = config.DB.Exec(query)
	if err != nil {
		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	query = "INSERT INTO authCode (auth_code, expires_in) VALUES (?, ?)"
	_, err = config.DB.Exec(query, authCode.Auth_Code, authCode.Expires_in)
	if err != nil {
		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(authCode)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot marshal authCode")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBody)
}

func GetTemporaryUserToken(w http.ResponseWriter, r *http.Request) {

	var authCode AuthCode

	var query string = "SELECT * FROM authCode LIMIT 1"
	row := config.DB.QueryRow(query)
	if err := row.Scan(&authCode.Auth_Code, &authCode.Expires_in); err != nil {

		if err == sql.ErrNoRows {
			http.Error(w, "Token does not exist", http.StatusNotFound)
			return
		}

		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(authCode)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot marshal authCode")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(jsonBody)
}

func GetPermanentUserToken(w http.ResponseWriter, r *http.Request) {

	var authToken AuthToken

	var query string = "SELECT * FROM authToken LIMIT 1"
	row := config.DB.QueryRow(query)
	if err := row.Scan(&authToken.Auth_token, &authToken.Id_user); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Token does not exist", http.StatusNotFound)
			return
		}

		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(authToken)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot marshal authCode")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(jsonBody)
}

func DeletePermanentUserToken(w http.ResponseWriter, r *http.Request) {

	var query string = "DELETE from authToken"
	_, err := config.DB.Exec(query)
	if err != nil {
		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
