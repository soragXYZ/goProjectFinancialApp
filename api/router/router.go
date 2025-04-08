package router

import (
	"net/http"

	"financialApp/api/resource/auth"
	"financialApp/api/resource/miscellaneous"
	"financialApp/api/resource/transaction"
	"financialApp/api/resource/webhook"

	"financialApp/api/router/middleware"
)

func New() *http.ServeMux {

	// to do: dispatch routes in submodules
	// https://dev.to/kengowada/go-routing-101-handling-and-grouping-routes-with-nethttp-4k0e

	router := http.NewServeMux()

	router.HandleFunc("GET /health/", middleware.Middleware(miscellaneous.HealthCheck))

	router.HandleFunc("POST /webhook/connection_synced/", middleware.Middleware(webhook.ConnectionSynced))

	router.HandleFunc("GET /transaction/", middleware.Middleware(transaction.GetTransactions))

	router.HandleFunc("POST /auth/permanentUserToken/", middleware.Middleware(auth.CreatePermanentUserToken))
	router.HandleFunc("GET /auth/permanentUserToken/", middleware.Middleware(auth.GetPermanentUserToken))
	router.HandleFunc("DELETE /auth/permanentUserToken/", middleware.Middleware(auth.DeletePermanentUserToken))

	router.HandleFunc("POST /auth/temporaryUserToken/", middleware.Middleware(auth.CreateTemporaryUserToken))
	router.HandleFunc("GET /auth/temporaryUserToken/", middleware.Middleware(auth.GetTemporaryUserToken))

	return router
}
