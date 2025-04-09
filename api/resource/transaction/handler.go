package transaction

import (
	"encoding/json"
	"net/http"
	"strconv"

	"financialApp/config"
)

func GetTransactions(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.Id, &tx.User_id, &tx.Bank_id, &tx.Date, &tx.Value, &tx.Transaction_type, &tx.Original_wording); err != nil {
			config.Logger.Error().Err(err).Msg("Cannot scan row")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		txs = append(txs, tx)
	}
	if err := rows.Err(); err != nil {
		config.Logger.Error().Err(err).Msg("Error in rows")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(txs)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot marshal txs")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(jsonBody)
}
