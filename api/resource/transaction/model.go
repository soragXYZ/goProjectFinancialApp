package transaction

import (
	"financialApp/api/resource/miscellaneous"
	"time"
)

// Models taken from https://docs.powens.com/api-reference/products/data-aggregation/bank-transactions#data-model

type AccountSchemeName struct {
	iban                     string
	bban                     string
	sort_code_account_number string
	cpan                     string
	tpan                     string
}

type Counterparty struct {
	label               string
	account_scheme_name AccountSchemeName
}
type Transaction struct {
	id                   int
	id_account           int
	application_date     time.Time
	date                 time.Time
	datetime             time.Time
	vdate                time.Time
	vdatetime            time.Time
	rdate                time.Time
	rdatetime            time.Time
	value                float32
	gross_value          float32
	transaction_type     string
	original_wording     string
	simplified_wording   string
	wording              string
	categories           []miscellaneous.Category
	date_scraped         time.Time
	coming               bool
	active               bool
	id_cluster           int
	comment              string
	last_update          time.Time
	deleted              time.Time
	original_value       float32
	original_gross_value float32
	original_currency    miscellaneous.Currency
	comission            float32
	commission_currency  miscellaneous.Currency
	card                 string
	counterparty         Counterparty
}

type TransactionsList struct {
	transactions    []Transaction
	first_date      time.Time
	last_date       time.Time
	result_min_date time.Time
	result_max_date time.Time
	_links          miscellaneous.PaginationLinks
}

type TransactionUpdateRequest struct {
	wording          string
	application_date time.Time
	categories       []miscellaneous.Category
	comment          string
	active           string
}
