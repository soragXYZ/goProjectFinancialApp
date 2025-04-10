package middleware

import (
	"net"
	"net/http"
	"slices"
	"time"

	"financialApp/config"
)

// See https://vishnubharathi.codes/blog/exploring-middlewares-in-go/
// https://gowebexamples.com/basic-middleware/
func Log(f http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var serverIP string
		if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
			serverIP = ipFromHostPort(addr.String())
		}

		var remoteIp string = ipFromHostPort(r.RemoteAddr)
		// Note: This is the originating ip, and thus not 100% the real IP (could be firewall or Load Balance for example)

		start := time.Now()
		f.ServeHTTP(w, r)
		timeTaken := time.Since(start)

		config.Logger.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("remote_ip", remoteIp).
			Str("server_ip", serverIP).
			Dur("msLatency", timeTaken).
			Msg("")
	})
}

func Whitelisted(f http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		remoteIp := ipFromHostPort(r.RemoteAddr)

		if !slices.Contains(config.Conf.Powens.WhitelistedIPs, remoteIp) {
			config.Logger.Warn().Msg("Unauthorized IP")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// If the IP is whitelisted, continue and process the request
		f.ServeHTTP(w, r)
	})
}

// Get the IP only and remove the port
func ipFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}
	if len(h) > 0 && h[0] == '[' { // IPv6 case, remove bracket []
		return h[1 : len(h)-1]
	}
	return h
}
