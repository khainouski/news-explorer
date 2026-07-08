package router

import (
	"net/http"
	"strings"
)

// ClientIP returns the real client address. Behind Traefik (the only way into this app - its own
// Service is ClusterIP-only, never exposed directly), r.RemoteAddr is just Traefik's own pod IP;
// the real client is in X-Forwarded-For, which Traefik sets on every proxied request. Leftmost
// entry is the original client - trusted here because Traefik is the sole entry point, so nothing
// external can inject a fake value upstream of it. Falls back to RemoteAddr when the header is
// absent (local dev, no proxy in front).
func ClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff == "" {
		return r.RemoteAddr
	}

	ip, _, _ := strings.Cut(xff, ",")

	return strings.TrimSpace(ip)
}
