package miscellaneous

// https://docs.powens.com/api-reference/products/data-aggregation/currencies#currency-object
type Currency struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	Precision uint   `json:"precision"`
}

type User struct {
	Id     int    `json:"id"`
	Signin string `json:"signin"`
}
