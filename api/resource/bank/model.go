package bank

import (
	"financialApp/api/resource/invest"
	"financialApp/api/resource/loan"
	"financialApp/api/resource/miscellaneous"
	"financialApp/api/resource/transaction"
)

// Models taken from https://docs.powens.com/api-reference/products/data-aggregation/bank-accounts#data-model

// Time sent by Powens API is not RFC3339
// so we store it as string
// https://stackoverflow.com/questions/25087960/json-unmarshal-time-that-isnt-in-rfc-3339-format

// https://docs.powens.com/api-reference/products/data-aggregation/bank-accounts#bankaccount-object
type BankAccountWebhook struct {
	Account_id    int                       `json:"id"`
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
	Loan          loan.Loan                 `json:"loan"`
	Investments   []invest.Invest           `json:"investments"`
	Transactions  []transaction.Transaction `json:"transactions"`
}
