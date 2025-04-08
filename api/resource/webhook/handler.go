package webhook

import (
	"encoding/json"
	"net/http"

	"financialApp/config"
)

func ConnectionSynced(w http.ResponseWriter, r *http.Request) {

	var conn Conn_synced

	err := json.NewDecoder(r.Body).Decode(&conn)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot decode r.Body")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	for _, account := range conn.Connection.Accounts {
		// Create bank account if does not exists. Otherwise, update last_update value

		var query string = "INSERT INTO bankAccount (bank_id, user_id, bank_number, original_name, balance, last_update, iban, account_type, usage_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE last_update=?"
		_, err = config.DB.Exec(
			query, account.Bank_id, account.User_id, account.Number, account.Original_name, account.Balance, account.Last_update, account.Iban, account.Account_type, account.Usage, account.Last_update,
		)
		if err != nil {
			config.Logger.Error().Err(err).Msg(query)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		for _, tx := range account.Transactions {
			// toDo: should bulk INSERT and verify duplicate on insert ?
			query = "INSERT INTO tx (tx_id, bank_id, tx_datetime, tx_value, tx_type, original_wording) VALUES (?, ?, ?, ?, ?, ?)"
			_, err = config.DB.Exec(query, tx.Id, tx.Bank_id, tx.Datetime, tx.Value, tx.Transaction_type, tx.Original_wording)
			if err != nil {
				config.Logger.Error().Err(err).Msg(query)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}
	}
}
