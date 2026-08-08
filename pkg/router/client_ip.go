package router

import (
	"net/http"
	"strings"
)

// ClientIP reads X-Forwarded-For - trusted because Traefik is the sole entry point into this app,
// so nothing external can spoof it. Falls back to RemoteAddr (local dev, no proxy in front).
func ClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff == "" {
		return r.RemoteAddr
	}

	ip, _, _ := strings.Cut(xff, ",")

	return strings.TrimSpace(ip)
}
