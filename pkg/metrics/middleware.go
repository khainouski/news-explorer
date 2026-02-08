package metrics

import (
	"net/http"
	"time"

	"github.com/khainouski/news-explorer/pkg/router"
)

func NewMiddleware(m *HTTPServer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			ww := router.WriterWrapper(w)
			next.ServeHTTP(ww, r)

			method := r.Method + " " + router.ExtractPath(r.Context())

			m.Duration(method, now)
			m.TotalInc(method, ww.Code())
		}

		return http.HandlerFunc(fn)
	}
}
