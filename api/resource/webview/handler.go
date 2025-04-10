package webview

import (
	"database/sql"
	"net/http"

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

	// The URL is builded as indicated here
	// https://docs.powens.com/api-reference/overview/webview#add-connection-flow
	var powensURL string = config.Conf.Powens.WebviewUrl + config.Conf.Other.Language + "/connect"
	var domain string = "domain=" + config.Conf.Powens.Domain
	var client_id string = "client_id=" + config.Conf.Powens.ClientId
	var redirect_uri string = "redirect_uri=" + config.Conf.Powens.RedirectUrl
	var connector_capabilities string = "connector_capabilities=bank,bankwealth"

	var connexionUrl string = powensURL + "?" + domain + "&" + client_id + "&" + redirect_uri + "&" + connector_capabilities + "&code=" + authCode.Auth_Code

	w.Write([]byte(connexionUrl))
}

func GetManageLink(w http.ResponseWriter, r *http.Request) {

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

	// The URL is builded as indicated here
	// https://docs.powens.com/api-reference/overview/webview#manage-connections
	var powensURL string = config.Conf.Powens.WebviewUrl + config.Conf.Other.Language + "/manage"
	var domain string = "domain=" + config.Conf.Powens.Domain
	var client_id string = "client_id=" + config.Conf.Powens.ClientId
	var redirect_uri string = "redirect_uri=" + config.Conf.Powens.RedirectUrl
	var connector_capabilities string = "connector_capabilities=bank,bankwealth"

	var connexionUrl string = powensURL + "?" + domain + "&" + client_id + "&" + redirect_uri + "&" + connector_capabilities + "&code=" + authCode.Auth_Code

	w.Write([]byte(connexionUrl))
}
