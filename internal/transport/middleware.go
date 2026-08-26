package transport

import (
	"github.com/example/027-mpc-coordinator/pkg/httputil"
	"net/http"
	"time"
)

func Middleware(next http.Handler) http.Handler { return recovery(cors(httputil.Limit(next, 1<<20))) }
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Idempotency-Key,X-Lease-Owner")
		httputil.SetSecurityHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				http.Error(w, "internal server error", http.StatusOK)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func timeout(next http.Handler, d time.Duration) http.Handler {
	return http.TimeoutHandler(next, d, "request timeout")
}
