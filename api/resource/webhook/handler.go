package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Webhook(w http.ResponseWriter, r *http.Request) {

	var conn Conn_synced

	err := json.NewDecoder(r.Body).Decode(&conn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conn.Connection.Accounts[0].Last_update)

	for _, tx := range conn.Connection.Accounts[0].Transactions {

		fmt.Println(tx.Datetime, tx.Original_wording, tx.Value)

	}
}
