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

	router.HandleFunc("GET /health/", middleware.Log(middleware.Whitelisted(miscellaneous.HealthCheck)))
	router.HandleFunc("/", middleware.Log(middleware.Whitelisted(miscellaneous.NotFound)))

	router.HandleFunc("POST /webhook/connection_synced/", middleware.Log(middleware.Whitelisted(webhook.ConnectionSynced)))

	router.HandleFunc("GET /transaction/", middleware.Log(middleware.Whitelisted(transaction.GetTransactions)))

	router.HandleFunc("POST /auth/permanentUserToken/", middleware.Log(middleware.Whitelisted(auth.CreatePermanentUserToken)))
	router.HandleFunc("GET /auth/permanentUserToken/", middleware.Log(middleware.Whitelisted(auth.GetPermanentUserToken)))
	router.HandleFunc("DELETE /auth/permanentUserToken/", middleware.Log(middleware.Whitelisted(auth.DeletePermanentUserToken)))

	router.HandleFunc("POST /auth/temporaryUserToken/", middleware.Log(middleware.Whitelisted(auth.CreateTemporaryUserToken)))
	router.HandleFunc("GET /auth/temporaryUserToken/", middleware.Log(middleware.Whitelisted(auth.GetTemporaryUserToken)))

	return router
}
