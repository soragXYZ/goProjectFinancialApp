package webhook

import (
	"encoding/json"
	"net/http"

	"financialApp/config"
)

func ConnectionSynced(w http.ResponseWriter, r *http.Request) {

	// // Display the JSON body in plain text, only for debug
	// buf, _ := io.ReadAll(r.Body)
	// rdr1 := io.NopCloser(bytes.NewBuffer(buf))
	// rdr2 := io.NopCloser(bytes.NewBuffer(buf))
	// io.Copy(os.Stdout, rdr1)
	// r.Body = rdr2

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

		// Proceed with transactions
		if len(account.Transactions) != 0 {

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

		// Proceed with invests
		if len(account.Investments) != 0 {

			for _, invest := range account.Investments {
				config.Logger.Trace().
					Int("invest ID", invest.Invest_id).
					Str("Label", invest.Label).
					Str("code", invest.Code).
					Str("code_type", invest.Code_type).
					Float32("unit_price", invest.Unit_price).
					Float32("unit_value", invest.Unit_value).
					Float32("valuation", invest.Valuation).
					Msg("")
			}
		}
	}
}
