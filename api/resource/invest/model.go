package invest

// Models taken from https://docs.powens.com/api-reference/products/wealth-aggregation/investments#data-model

// https://docs.powens.com/api-reference/products/wealth-aggregation/investments#investment-object
type Invest struct {
	Invest_id    int     `json:"id"`
	Account_id   int     `json:"id_account"`
	User_id      int     `json:"-"` // absent in base data, field added for simplicity
	Type_id      int     `json:"id_type"`
	Label        string  `json:"label"`
	Code         string  `json:"code"`
	Code_type    string  `json:"code_type"`
	Stock_symbol string  `json:"stock_symbol"`
	Quantity     float32 `json:"quantity"`
	Unit_price   float32 `json:"unitprice"`
	Unit_value   float32 `json:"unitvalue"`
	Valuation    float32 `json:"valuation"`
	Diff         float32 `json:"diff"`
	Diff_percent float32 `json:"diff_percent"`
	Last_update  string  `json:"last_update"`
}
