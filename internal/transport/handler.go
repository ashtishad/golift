package transport

import (
	"log/slog"
	"net/http"

	"github.com/ashtishad/golift/internal/domain"
)

func ProxyRequestHandler(serverPool domain.ServerPooler, l *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetServer := serverPool.SelectServer()
		if targetServer == nil {
			l.Error("target server unavailable", "srv", targetServer.GetURL())
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		// Update the request's host to match the target URL.
		targetURL := targetServer.GetURL()
		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = targetURL.Host

		// Add X-Forwarded-For and X-Forwarded-Proto
		r.Header.Set("X-Forwarded-For", r.RemoteAddr)
		if r.TLS != nil {
			r.Header.Set("X-Forwarded-Proto", "https")
		} else {
			r.Header.Set("X-Forwarded-Proto", "http")
		}

		// Serve the request using reverseProxy of server instance.
		targetServer.Serve(w, r)
	}
}
