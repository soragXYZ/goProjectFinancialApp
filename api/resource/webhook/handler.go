package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"financialApp/config"
)

func ConnectionSynced(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/webhook/connection_synced/" {
		http.NotFound(w, r)
		return
	}

	var conn Conn_synced

	err := json.NewDecoder(r.Body).Decode(&conn)
	if err != nil {
		log.Fatal(err)
	}

	for _, account := range conn.Connection.Accounts {
		// Create bank account if not exists. Otherwise, update last_update value
		_, err = config.DB.Exec(
			"INSERT INTO bankAccount (bank_id, user_id, bank_number, original_name, balance, last_update, iban, account_type, usage_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE last_update=?",
			account.Bank_id, account.User_id, account.Number, account.Original_name, account.Balance, account.Last_update, account.Iban, account.Account_type, account.Usage, account.Last_update,
		)
		if err != nil {
			log.Fatal("INSERT INTO bankAccount (bank_id, user_id, bank_number, original_name, balance, last_update, iban, account_type, usage_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE last_update=?: %v", err)
		}

		for _, tx := range account.Transactions {

			fmt.Println(tx.Id, tx.Bank_id, tx.Datetime, tx.Value, tx.Transaction_type, tx.Original_wording)
			// toDo: should bulk INSERT and verify duplicate on insert ?
			_, err = config.DB.Exec(
				"INSERT INTO tx (tx_id, bank_id, tx_datetime, tx_value, tx_type, original_wording) VALUES (?, ?, ?, ?, ?, ?)",
				tx.Id, tx.Bank_id, tx.Datetime, tx.Value, tx.Transaction_type, tx.Original_wording,
			)
			if err != nil {
				log.Fatal("INSERT INTO tx (tx_id, bank_id, tx_datetime, tx_value, tx_type, original_wording) VALUES (?, ?, ?, ?, ?, ?): %v", err)
			}
		}
	}
}
