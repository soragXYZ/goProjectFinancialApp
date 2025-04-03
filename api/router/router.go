package router

import (
	"fmt"
	"log"
	"net/http"

	"financialApp/api/resource/auth"
	"financialApp/api/resource/health"
)

func New() {

	http.HandleFunc("GET /health/", health.HealthCheck)

	http.HandleFunc("POST /auth/permanentUserToken/", auth.CreatePermanentUserToken)
	http.HandleFunc("GET /auth/permanentUserToken/", auth.GetPermanentUserToken)
	http.HandleFunc("DELETE /auth/permanentUserToken/", auth.DeletePermanentUserToken)

	fmt.Println("Server running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
