package auth

// Models taken from https://docs.powens.com/api-reference/overview/authentication#data-model

type AuthTokenInitRequest struct {
	Client_id     string `json:"client_id"`
	Client_secret string `json:"client_secret"`
}

type AuthToken struct {
	Auth_token string `json:"auth_token"`
	Token_type string `json:"type"`
	Id_user    int    `json:"id_user"`
	Expires_in int    `json:"expires_in"`
}

// type AuthTokenType struct {
// 	singleAccess
// }

type AuthCode struct {
	code       string
	code_type  string
	access     string
	expires_in int
}

type AuthTokenExchangeRequest struct {
	grant_type    string
	client_id     string
	client_secret string
	code          string
}

type AuthTokenExchange struct {
	access_token string
	token_type   string
}

type AuthServiceTokenRequest struct {
	grant_type    string
	client_id     string
	client_secret string
	scope         []string
}

type AuthServiceToken struct {
	token string
	scope string
}

type AuthRenewRequest struct {
	grant_type      string
	client_id       string
	client_secret   string
	id_user         int
	revoke_previous bool
}
