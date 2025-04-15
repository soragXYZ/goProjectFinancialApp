package webview

// Models taken from https://docs.powens.com/api-reference/overview/authentication#data-model

// https://docs.powens.com/api-reference/overview/authentication#authcode-object
type authCode struct {
	Auth_Code   string `json:"code"`
	Error       error  `json:"-"`
	ErrorString string `json:"-"`
}
