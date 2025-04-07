package transaction

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"financialApp/config"
)

func GetTransactions(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/transaction/" {
		http.NotFound(w, r)
		return
	}

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

	rows, err := config.DB.Query("SELECT * FROM tx ORDER BY tx_datetime DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.Id, &tx.Bank_id, &tx.Datetime, &tx.Value, &tx.Transaction_type, &tx.Original_wording); err != nil {
			log.Fatal(err)
		}
		txs = append(txs, tx)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	jsonBody, err := json.Marshal(txs)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonBody)

}
