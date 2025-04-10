package webhook

import (
	"encoding/json"
	"net/http"

	"financialApp/config"
)

func ConnectionSynced(w http.ResponseWriter, r *http.Request) {

	// Display the JSON body in plain text, for debug
	// io.Copy(os.Stdout, r.Body)

	var conn Conn_synced

	// Error here in decode sometimes, one wrong field in conn ?
	err := json.NewDecoder(r.Body).Decode(&conn)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot decode r.Body")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	for _, account := range conn.Connection.Accounts {

		// Create bank account if does not exists. Otherwise, update last_update value
		var query string = "INSERT INTO bankAccount (bank_id, user_id, bank_number, original_name, balance, last_update, iban, currency, account_type, usage_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE last_update=?"
		_, err = config.DB.Exec(
			query, account.Bank_id, account.User_id, account.Number, account.Original_name, account.Balance, account.Last_update, account.Iban, account.Currency.Id, account.Account_type, account.Usage, account.Last_update,
		)
		if err != nil {
			config.Logger.Error().Err(err).Msg(query)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// If no transactions, process the next account
		if len(account.Transactions) == 0 {
			continue
		}

		// Bulk insert txs
		query = "INSERT INTO tx (tx_id, user_id, bank_id, tx_date, tx_value, tx_type, original_wording) VALUES "
		vals := []any{}
		for _, tx := range account.Transactions {
			query += "(?, ?, ?, ?, ?, ?, ?),"
			vals = append(vals, tx.Id, account.User_id, tx.Bank_id, tx.Date, tx.Value, tx.Transaction_type, tx.Original_wording)
		}
		query = query[0 : len(query)-1]

		_, err := config.DB.Exec(query, vals...)
		if err != nil {
			config.Logger.Error().Err(err).Msg(query)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
