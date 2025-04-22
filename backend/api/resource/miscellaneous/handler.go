package miscellaneous

import (
	"net/http"

	"financialApp/config"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Healthy"))
}

func Version(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("0.0.1"))
}

func NotFound(w http.ResponseWriter, r *http.Request) {

	config.Logger.Warn().Msg("Someone used a wrong path")
	http.Error(w, "Page does not exist", http.StatusNotFound)
}
