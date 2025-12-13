package middleware

import (
	"net/http"
	"strings"
)

// NormalizeHostnameMiddleware ensures consistent hostname usage
// Redirects 127.0.0.1 to localhost to prevent cookie/sessionStorage issues
func NormalizeHostnameMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		
		// Check if request is using 127.0.0.1 instead of localhost
		if strings.HasPrefix(host, "127.0.0.1:") {
			// Replace 127.0.0.1 with localhost
			newHost := strings.Replace(host, "127.0.0.1", "localhost", 1)
			newURL := "http://" + newHost + r.URL.RequestURI()
			
			// Redirect to localhost version
			http.Redirect(w, r, newURL, http.StatusMovedPermanently)
			return
		}
		
		// Continue with normal processing
		next.ServeHTTP(w, r)
	})
}

