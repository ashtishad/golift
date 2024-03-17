package transport

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ashtishad/golift/internal/domain"
)

func ProxyRequestHandler(serverPool domain.ServerPooler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("sewhrnwkj")
		fmt.Println("hwgwejkrhmjk")
		targetServer := serverPool.SelectServer()
		if targetServer == nil {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		// Update the request's host to match the target URL.
		targetURL := targetServer.GetURL()
		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = targetURL.Host

		// Serve the request using reverseProxy of server instance.
		targetServer.Serve(w, r)
	}
}
