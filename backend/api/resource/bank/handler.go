package bank

import (
	"encoding/json"
	"net/http"

	"financialApp/config"
)

func GetAccounts(w http.ResponseWriter, r *http.Request) {

	var accounts []BankAccount

	var query string = "SELECT * FROM bankAccount ORDER BY original_name"
	rows, err := config.DB.Query(query)
	if err != nil {
		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var account BankAccount
		if err := rows.Scan(&account.Account_id, &account.User_id, &account.Number, &account.Original_name, &account.Balance, &account.Last_update, &account.Iban, &account.Currency, &account.Account_type, &account.Usage); err != nil {
			config.Logger.Error().Err(err).Msg("Cannot scan row")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		config.Logger.Error().Err(err).Msg("Error in rows")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(accounts)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot marshal accounts")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(jsonBody)
}
