package middleware

import (
	"net/http"
	"time"

	"financialApp/config"
)

// See https://vishnubharathi.codes/blog/exploring-middlewares-in-go/
// https://gowebexamples.com/basic-middleware/
func Middleware(f http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		f.ServeHTTP(w, r)
		config.Logger.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Dur("msLatency", time.Since(start)).
			Msg("")
	})
}
