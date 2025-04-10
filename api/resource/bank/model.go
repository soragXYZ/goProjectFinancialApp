package bank

import (
	"financialApp/api/resource/invest"
	"financialApp/api/resource/miscellaneous"
	"financialApp/api/resource/transaction"
	"time"
)

// Models taken from https://docs.powens.com/api-reference/products/data-aggregation/bank-accounts#data-model

// Time sent by Powens API is not RFC3339
// so we store it as string
// https://stackoverflow.com/questions/25087960/json-unmarshal-time-that-isnt-in-rfc-3339-format

// https://docs.powens.com/api-reference/products/data-aggregation/bank-accounts#loan-object
type Loan struct {
	Total_amount         float32   `json:"total_amount"`
	Available_amount     float32   `json:"available_amount"`
	Used_amount          float32   `json:"used_amount"`
	Subscription_date    time.Time `json:"subscription_date"`
	Maturity_date        time.Time `json:"maturity_date"`
	Start_repayment_date time.Time `json:"start_repayment_date"`
	Deferred             bool      `json:"deferred"`
	Next_payment_amount  float32   `json:"next_payment_amount"`
	Next_payment_date    time.Time `json:"next_payment_date"`
	Rate                 float32   `json:"rate"`
	Nb_payments_left     uint      `json:"nb_payments_left"`
	Nb_payments_done     uint      `json:"nb_payments_done"`
	Nb_payments_total    uint      `json:"nb_payments_total"`
	Last_payment_amount  float32   `json:"last_payment_amount"`
	Last_payment_date    time.Time `json:"last_payment_date"`
	Account_label        string    `json:"account_label"`
	Insurance_label      string    `json:"insurance_label"`
	Insurance_amount     string    `json:"insurance_amount"`
	Insurance_rate       float32   `json:"insurance_rate"`
	Duration             uint      `json:"duration"`
	Loan_type            string    `json:"loan_type"`
}

// https://docs.powens.com/api-reference/products/data-aggregation/bank-accounts#bankaccount-object
type BankAccountWebhook struct {
	Bank_id       int                       `json:"id"`
	User_id       int                       `json:"id_user"`
	Number        string                    `json:"number"`
	Original_name string                    `json:"original_name"`
	Balance       float32                   `json:"balance"`
	Last_update   string                    `json:"last_update"`
	Iban          string                    `json:"iban"`
	Currency      miscellaneous.Currency    `json:"currency"`
	Account_type  string                    `json:"type"`
	Error         string                    `json:"error"` // not needed ?
	Usage         string                    `json:"usage"`
	Loan          Loan                      `json:"loan"` // not needed ?
	Investments   []invest.Invest           `json:"investments"`
	Transactions  []transaction.Transaction `json:"transactions"`
}
