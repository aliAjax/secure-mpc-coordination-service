package transport

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type countingBody struct {
	io.Reader
	closed bool
}

func (c *countingBody) Close() error { c.closed = true; return nil }

func TestRecoveryWritesClean500(t *testing.T) {
	h := recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestRecoverySkipsAfterCommit(t *testing.T) {
	h := recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("partial"))
		panic("boom")
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Body.String() != "partial" {
		t.Fatalf("expected only partial body, got %q", rec.Body.String())
	}
}

func TestDecodeClosesBody(t *testing.T) {
	body := &countingBody{Reader: strings.NewReader(`{"tenant_id":"acme"}`)}
	req := httptest.NewRequest("POST", "/", body)
	var v struct {
		TenantID string `json:"tenant_id"`
	}
	decode(httptest.NewRecorder(), req, &v)
	if !body.closed {
		t.Fatal("expected request body to be closed")
	}
}
