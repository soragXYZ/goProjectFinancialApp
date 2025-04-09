package webview

import (
	"database/sql"
	"net/http"
	"os"

	"financialApp/api/resource/auth"
	"financialApp/config"
)

func GetCreationLink(w http.ResponseWriter, r *http.Request) {

	var authCode auth.AuthCode

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

	var powensURL string = "https://webview.powens.com/fr/connect"
	var domain string = "domain=" + os.Getenv("DOMAIN")
	var client_id string = "client_id=" + os.Getenv("CLIENT_ID")
	var redirect_uri string = "redirect_uri=https://example.com" // to be changed
	var connector_capabilities string = "connector_capabilities=bank,bankwealth"

	var connexionUrl string = powensURL + "?" + domain + "&" + client_id + "&" + redirect_uri + "&" + connector_capabilities + "&code=" + authCode.Auth_Code

	w.Write([]byte(connexionUrl))
}
