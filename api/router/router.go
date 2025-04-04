package router

import (
	"fmt"
	"log"
	"net/http"

	"financialApp/api/resource/auth"
	"financialApp/api/resource/miscellaneous"
)

func New() {

	http.HandleFunc("GET /health/", miscellaneous.HealthCheck)

	http.HandleFunc("POST /auth/webhook/", auth.Webhook)

	http.HandleFunc("POST /auth/permanentUserToken/", auth.CreatePermanentUserToken)
	http.HandleFunc("GET /auth/permanentUserToken/", auth.GetPermanentUserToken)
	http.HandleFunc("DELETE /auth/permanentUserToken/", auth.DeletePermanentUserToken)

	http.HandleFunc("POST /auth/temporaryUserToken/", auth.CreateTemporaryUserToken)
	http.HandleFunc("GET /auth/temporaryUserToken/", auth.GetTemporaryUserToken)

	fmt.Println("Server running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
