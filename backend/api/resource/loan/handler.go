package loan

import (
	"encoding/json"
	"net/http"

	"financialApp/config"
)

func GetLoans(w http.ResponseWriter, r *http.Request) {

	var loans []Loan

	var query string = "SELECT * FROM loan"
	rows, err := config.DB.Query(query)
	if err != nil {
		config.Logger.Error().Err(err).Msg(query)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var account Loan
		if err := rows.Scan(&account.Loan_account_id, &account.Total_amount, &account.Available_amount, &account.Used_amount, &account.Subscription_date, &account.Maturity_date, &account.Start_repayment_date, &account.Deferred, &account.Next_payment_amount, &account.Next_payment_date, &account.Rate, &account.Nb_payments_left, &account.Nb_payments_done, &account.Nb_payments_total, &account.Last_payment_amount, &account.Last_payment_date, &account.Account_label, &account.Insurance_label, &account.Insurance_amount, &account.Insurance_rate, &account.Duration, &account.Loan_type); err != nil {
			config.Logger.Error().Err(err).Msg("Cannot scan row")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		loans = append(loans, account)
	}
	if err := rows.Err(); err != nil {
		config.Logger.Error().Err(err).Msg("Error in rows")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(loans)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Cannot marshal loans")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(jsonBody)
}
