package bank

import (
	"financialApp/api/resource/miscellaneous"
	"time"
)

// Models taken from https://docs.powens.com/api-reference/products/data-aggregation/bank-accounts#data-model

type AccountType struct {
	id             uint
	name           string
	id_parent      uint
	is_invest      bool
	display_name   string
	display_name_p string
}

type Loan struct {
	total_amount         float32
	available_amount     float32
	used_amount          float32
	subscription_date    time.Time
	maturity_date        time.Time
	start_repayment_date time.Time
	deferred             bool
	next_payment_amount  float32
	next_payment_date    time.Time
	rate                 float32
	nb_payments_left     uint
	nb_payments_done     uint
	nb_payments_total    uint
	last_payment_amount  float32
	last_payment_date    time.Time
	account_label        string
	insurance_label      string
	insurance_amount     string
	insurance_rate       float32
	duration             uint
	loan_type            string
}

type BankAccount struct {
	id            int
	id_connection int
	id_user       int
	id_source     int
	id_parent     int
	number        string
	original_name string
	balance       float32
	coming        float32
	display       bool
	last_update   time.Time
	deleted       time.Time
	disabled      time.Time
	iban          string
	currency      miscellaneous.Currency
	account_type  AccountType
	id_type       int
	bookmarked    int
	name          string
	error         string
	usage         string
	company_name  string
	loan          Loan
}

type Balance struct {
	currency string
	amount   float32
}

type BankAccountsList struct {
	balance         float32
	accounts        []BankAccount
	balances        []Balance
	coming_balances []Balance
}

type BankAccountUpdateRequest struct {
	display    bool
	name       string
	disabled   bool
	bookmarked bool
	usage      string
}
