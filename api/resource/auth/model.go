package auth

// Models taken from https://docs.powens.com/api-reference/overview/authentication#data-model

// https://docs.powens.com/api-reference/overview/authentication#authtokeninitrequest-object
type AuthTokenInitRequest struct {
	Client_id     string `json:"client_id"`
	Client_secret string `json:"client_secret"`
}

// https://docs.powens.com/api-reference/overview/authentication#authtoken-object
type AuthToken struct {
	Auth_token string `json:"auth_token"`
	Id_user    int    `json:"id_user"`
}

// https://docs.powens.com/api-reference/overview/authentication#authcode-object
type AuthCode struct {
	Auth_Code  string `json:"code"`
	Expires_in int    `json:"expires_in"`
}
