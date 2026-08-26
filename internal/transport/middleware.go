package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/example/027-mpc-coordinator/pkg/httputil"
)

func Middleware(next http.Handler) http.Handler {
	// recovery sits inside Limit so MaxBytesReader receives the server's
	// unwrapped *response and its connection-close signal (requestTooLarge)
	// is not lost behind the recorder wrapper.
	return cors(httputil.Limit(recovery(next), 1<<20))
}
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
		rw := &recorder{ResponseWriter: w}
		defer func() {
			// Always close the request body so oversized or unread
			// bodies don't outlive the handler. For bodies already
			// drained by decode this is a harmless no-op. The body is
			// not closed again by the server because the *response
			// has already been signaled to close (closeAfterReply
			// via requestTooLarge) for oversized requests.
			_ = r.Body.Close()
			if v := recover(); v != nil {
				slog.Error("panic recovered", "error", v, "path", r.URL.Path, "stack", string(debug.Stack()))
				// If the handler already committed the response, writing
				// again would only garble the half-written body.
				if !rw.committed {
					rw.Header().Set("Content-Type", "application/json")
					rw.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(rw).Encode(map[string]any{"error": map[string]string{"code": http.StatusText(http.StatusInternalServerError), "message": "internal server error"}})
				}
			}
		}()
		next.ServeHTTP(rw, r)
	})
}
func timeout(next http.Handler, d time.Duration) http.Handler {
	return http.TimeoutHandler(next, d, "request timeout")
}

type recorder struct {
	http.ResponseWriter
	committed bool
}

func (r *recorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }
func (r *recorder) WriteHeader(code int) {
	if r.committed {
		return
	}
	r.committed = true
	r.ResponseWriter.WriteHeader(code)
}
func (r *recorder) Write(b []byte) (int, error) {
	if !r.committed {
		r.committed = true
	}
	return r.ResponseWriter.Write(b)
}
