package health

import (
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Healthy"))
}
