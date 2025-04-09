package transaction

// Models taken from https://docs.powens.com/api-reference/products/data-aggregation/bank-transactions#data-model

// https://docs.powens.com/api-reference/products/data-aggregation/bank-transactions#transaction-object
type Transaction struct {
	Id               int     `json:"id"`
	Bank_id          int     `json:"id_account"`
	User_id          int     `json:"-"` // absent in base data, field added for simplicity
	Date             string  `json:"date"`
	Value            float32 `json:"value"`
	Transaction_type string  `json:"type"`
	Original_wording string  `json:"original_wording"`
}
