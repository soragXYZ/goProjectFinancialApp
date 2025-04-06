package transaction

// Models taken from https://docs.powens.com/api-reference/products/data-aggregation/bank-transactions#data-model

type Transaction struct {
	Id               int     `json:"id"`
	Id_account       int     `json:"id_account"`
	Datetime         string  `json:"datetime"`
	Value            float32 `json:"value"`
	Transaction_type string  `json:"type"`
	Original_wording string  `json:"original_wording"`
}
