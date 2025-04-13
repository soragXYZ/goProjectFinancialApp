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

		config.Logger.Trace().
			Int("account_id", account.Account_id).
			Str("account_name", account.Original_name).
			Str("last_update", account.Last_update).
			Int("user_id", account.User_id).
			Msg("Account Update")

		// Create bank account if it does not exists. Otherwise, update last_update value
		var query string = "INSERT INTO bankAccount (account_id, user_id, bank_number, original_name, balance, last_update, iban, currency, account_type, usage_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE last_update=?"
		_, err = config.DB.Exec(
			query, account.Account_id, account.User_id, account.Number, account.Original_name, account.Balance, account.Last_update, account.Iban, account.Currency.Id, account.Account_type, account.Usage, account.Last_update,
		)
		if err != nil {
			config.Logger.Error().Err(err).Msg(query)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Proceed with transactions
		if len(account.Transactions) != 0 {

			// Bulk insert txs
			query = "INSERT INTO tx (tx_id, user_id, account_id, tx_date, tx_value, tx_type, original_wording) VALUES "
			vals := []any{}
			for _, tx := range account.Transactions {

				config.Logger.Trace().
					Int("account_id", tx.Account_id).
					Str("date", tx.Date).
					Str("original_wording", tx.Original_wording).
					Int("tx_id", tx.Id).
					Float32("value", tx.Value).
					Msg("Tx update")

				query += "(?, ?, ?, ?, ?, ?, ?),"
				vals = append(vals, tx.Id, account.User_id, tx.Account_id, tx.Date, tx.Value, tx.Transaction_type, tx.Original_wording)
			}

			// remove last comma
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

			// Bulk insert invests
			query = "INSERT INTO invest (invest_id, account_id, type_id, invest_label, invest_code, invest_code_type, stock_symbol, quantity, unit_price, unit_value, valuation, diff, diff_percent, last_update) VALUES "
			vals := []any{}
			for _, invest := range account.Investments {

				config.Logger.Trace().
					Str("account_name", account.Original_name).
					Str("code", invest.Code).
					Int("invest_id", invest.Invest_id).
					Str("label", invest.Label).
					Float32("unit_price", invest.Unit_price).
					Float32("unit_value", invest.Unit_value).
					Float32("valuation", invest.Valuation).
					Msg("Invest update")

				query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
				vals = append(vals, invest.Invest_id, invest.Account_id, invest.Type_id, invest.Label, invest.Code, invest.Code_type, invest.Stock_symbol, invest.Quantity, invest.Unit_price, invest.Unit_value, invest.Valuation, invest.Diff, invest.Diff_percent, invest.Last_update)
			}

			// remove last comma
			query = query[0 : len(query)-1]

			// if duplicate entry, update the field by the new value
			query = query + "AS new(a, b, c, d, e, f, g, Nquantity, Nunit_price, Nunit_value, Nvaluation, Ndiff, Ndiff_percent, Nlast_update)"
			query = query + "ON DUPLICATE KEY UPDATE quantity=Nquantity, unit_price=Nunit_price, unit_value=Nunit_value, valuation=Nvaluation, diff=Ndiff, diff_percent=Ndiff_percent, last_update=Nlast_update"

			_, err := config.DB.Exec(query, vals...)
			if err != nil {
				config.Logger.Error().Err(err).Msg(query)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}
	}
}
