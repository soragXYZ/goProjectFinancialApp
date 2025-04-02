package router

import (
	"fmt"
	"log"
	"net/http"

	"financialApp/api/resource/auth"
	"financialApp/api/resource/health"
)

func New() {

	http.HandleFunc("/health/", health.HealthCheck)

	// should be changed in GET POST DELETE
	http.HandleFunc("/auth/createPermanentUserToken/", auth.CreatePermanentUserToken)
	http.HandleFunc("/auth/getPermanentUserToken/", auth.GetPermanentUserToken)

	fmt.Println("Server running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
