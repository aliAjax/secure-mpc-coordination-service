package transport

import (
	"encoding/json"
	"errors"
	"github.com/example/027-mpc-coordinator/internal/application"
	"github.com/example/027-mpc-coordinator/internal/domain"
	"github.com/example/027-mpc-coordinator/internal/observability"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	svc *application.Service
	log *slog.Logger
}

func NewServer(s *application.Service, l *slog.Logger) *Server { return &Server{svc: s, log: l} }
func (s *Server) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", s.health)
	m.HandleFunc("/readyz", s.health)
	m.HandleFunc("/api/v1/computations", s.computations)
	m.HandleFunc("/api/v1/computations/", s.computation)
	m.HandleFunc("/api/v1/rounds/", s.round)
	m.HandleFunc("/api/v1/evidence/", s.evidence)
	m.HandleFunc("/metrics", s.metrics)
	return requestLog(m, s.log)
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	write(w, http.StatusOK, map[string]string{"status": "ok", "service": "mpc-coordinator"})
}
func (s *Server) metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("mpc_requests_total 0\nmpc_reconstructions_total 0\n"))
}
func (s *Server) computations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		out, e := s.svc.List(r.Context())
		if e != nil {
			fail(w, r, e)
			return
		}
		write(w, http.StatusOK, map[string]any{"items": out, "count": len(out)})
	case http.MethodPost:
		var req application.CreateRequest
		if !decode(w, r, &req) {
			return
		}
		c, e := s.svc.Create(r.Context(), req, r.Header.Get("Idempotency-Key"))
		if e != nil {
			fail(w, r, e)
			return
		}
		write(w, http.StatusCreated, c)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func (s *Server) computation(w http.ResponseWriter, r *http.Request) {
	id, action := pathParts(r.URL.Path, "/api/v1/computations/")
	if id == "" {
		fail(w, r, domain.ErrNotFound)
		return
	}
	switch action {
	case "":
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		c, e := s.svc.Get(r.Context(), id)
		if e != nil {
			fail(w, r, e)
			return
		}
		write(w, http.StatusOK, c)
	case "start":
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		rnd, e := s.svc.Start(r.Context(), id)
		if e != nil {
			fail(w, r, e)
			return
		}
		write(w, http.StatusCreated, rnd)
	case "abort":
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if e := s.svc.Abort(r.Context(), id); e != nil {
			fail(w, r, e)
			return
		}
		write(w, http.StatusOK, map[string]string{"status": "aborted"})
	case "demo-shares":
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			RoundID string `json:"round_id"`
			Secret  string `json:"secret"`
		}
		if !decode(w, r, &req) {
			return
		}
		res, e := s.svc.DemoShares(r.Context(), id, req.RoundID, req.Secret)
		if e != nil {
			fail(w, r, e)
			return
		}
		write(w, http.StatusOK, res)
	case "participants":
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req domain.Participant
		if !decode(w, r, &req) {
			return
		}
		p, e := s.svc.RegisterParticipant(r.Context(), id, req)
		if e != nil {
			fail(w, r, e)
			return
		}
		write(w, http.StatusCreated, p)
	case "reconstruct":
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			RoundID string `json:"round_id"`
		}
		if !decode(w, r, &req) {
			return
		}
		out, e := s.svc.Reconstruct(r.Context(), id, req.RoundID)
		if e != nil {
			fail(w, r, e)
			return
		}
		write(w, http.StatusOK, out)
	default:
		fail(w, r, domain.ErrNotFound)
	}
}
func (s *Server) round(w http.ResponseWriter, r *http.Request) {
	id, action := pathParts(r.URL.Path, "/api/v1/rounds/")
	if action != "shares" || r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ParticipantID string `json:"participant_id"`
		Index         int    `json:"index"`
		Value         string `json:"value"`
		Commitment    string `json:"commitment"`
		Signature     string `json:"signature"`
	}
	if !decode(w, r, &req) {
		return
	}
	owner := r.Header.Get("X-Lease-Owner")
	if owner == "" {
		owner = "http-api"
	}
	rnd, e := s.svc.SubmitShare(r.Context(), id, domain.Share{ParticipantID: req.ParticipantID, Index: req.Index, Value: req.Value, Commitment: req.Commitment, Signature: req.Signature}, owner)
	if e != nil {
		fail(w, r, e)
		return
	}
	write(w, http.StatusAccepted, rnd)
}
func (s *Server) evidence(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/evidence/")
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	items, e := s.svc.Evidence(r.Context(), id)
	if e != nil {
		fail(w, r, e)
		return
	}
	write(w, http.StatusOK, map[string]any{"items": items})
}
func (s *Server) reconstruct(w http.ResponseWriter, r *http.Request) {
	id, _ := pathParts(r.URL.Path, "/api/v1/computations/")
	out, e := s.svc.Reconstruct(r.Context(), id, "")
	if e != nil {
		fail(w, r, e)
		return
	}
	write(w, http.StatusOK, out)
}
func pathParts(path, prefix string) (string, string) {
	rest := strings.TrimPrefix(path, prefix)
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "", ""
	}
	a := ""
	if len(parts) > 1 {
		a = parts[1]
	}
	return parts[0], a
}
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	de := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	if e := de.Decode(v); e != nil {
		fail(w, r, domain.ErrInvalid)
		return false
	}
	return true
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, r *http.Request, e error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(e, domain.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(e, domain.ErrInvalid), errors.Is(e, domain.ErrThreshold):
		status = http.StatusBadRequest
	case errors.Is(e, domain.ErrConflict), errors.Is(e, domain.ErrReplay), errors.Is(e, domain.ErrLeaseLost):
		status = http.StatusConflict
	}
	write(w, status, map[string]any{"error": map[string]string{"code": http.StatusText(status), "message": e.Error()}, "trace_id": observability.TraceID(r.Context())})
}
func requestLog(next http.Handler, l *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := observability.WithTrace(r.Context())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		if l != nil {
			observability.LogContext(l, ctx, "http_request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
		}
	})
}
