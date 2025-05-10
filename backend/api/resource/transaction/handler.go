package transaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"financialApp/config"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var tx Transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = config.DB.Exec(
		"INSERT INTO tx (tx_id, user_id, account_id, tx_date, tx_value, tx_type, original_wording) VALUES (?, ?, ?, ?, ?, ?, ?)",
		tx.Id, tx.User_id, tx.Account_id, tx.Date, tx.Value, tx.Transaction_type, tx.Original_wording)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ReadTransaction(w http.ResponseWriter, r *http.Request) {

	// To Do: Change from page-based to cursor-based pagination
	// https://www.merge.dev/blog/rest-api-pagination

	// Extract 'page' and 'limit' from query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 || limit > 50 {
		limit = 50 // Default to 50 txs per page
	}

	offset := (page - 1) * limit

	var txs []Transaction

	var query string = "SELECT * FROM tx ORDER BY tx_date DESC LIMIT ? OFFSET ?"
	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.Id, &tx.User_id, &tx.Account_id, &tx.Date, &tx.Value, &tx.Transaction_type, &tx.Original_wording, &tx.Pinned); err != nil {
			config.Logger.Error().Err(err).Msg("Cannot scan row")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		txs = append(txs, tx)
	}
	if err := rows.Err(); err != nil {
		config.Logger.Error().Err(err).Msg("Error in rows")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(txs)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot marshal txs")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonBody)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {

	txId := r.PathValue("id")
	var tx Transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var query string = "UPDATE tx SET tx_date=?, tx_value=?, tx_type=?, original_wording=?, pinned=? WHERE tx_id=?"
	_, err = config.DB.Exec(query, tx.Date, tx.Value, tx.Transaction_type, tx.Original_wording, tx.Pinned, txId)
	if err != nil {
		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {

	txId := r.PathValue("id")
	fmt.Print(txId)

	var query string = "DELETE from tx WHERE tx_id=?"
	_, err := config.DB.Exec(query, txId)
	if err != nil {
		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
