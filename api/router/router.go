package router

import (
	"net/http"

	"financialApp/api/resource/auth"
	"financialApp/api/resource/miscellaneous"
	"financialApp/api/resource/transaction"
	"financialApp/api/resource/webhook"
)

func New() *http.ServeMux {

	// to do: dispatch routes in submodules
	// https://dev.to/kengowada/go-routing-101-handling-and-grouping-routes-with-nethttp-4k0e

	router := http.NewServeMux()
	router.HandleFunc("GET /health/", miscellaneous.HealthCheck)

	router.HandleFunc("POST /webhook/connection_synced/", webhook.ConnectionSynced)

	router.HandleFunc("GET /transaction/", transaction.GetTransactions)

	router.HandleFunc("POST /auth/permanentUserToken/", auth.CreatePermanentUserToken)
	router.HandleFunc("GET /auth/permanentUserToken/", auth.GetPermanentUserToken)
	router.HandleFunc("DELETE /auth/permanentUserToken/", auth.DeletePermanentUserToken)

	router.HandleFunc("POST /auth/temporaryUserToken/", auth.CreateTemporaryUserToken)
	router.HandleFunc("GET /auth/temporaryUserToken/", auth.GetTemporaryUserToken)

	return router
}
