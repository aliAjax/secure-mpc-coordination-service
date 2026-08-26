package observability

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestTraceIDPropagation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	child := WithTrace(parent)
	cancel()
	if child.Err() == nil {
		t.Fatal("expected child context to inherit parent cancellation")
	}
}

type captureHandler struct {
	keys []string
}

func (h *captureHandler) Enabled(context.Context, slog.Level) bool { return true }
func (h *captureHandler) Handle(_ context.Context, r slog.Record) error {
	r.Attrs(func(a slog.Attr) bool {
		h.keys = append(h.keys, a.Key)
		return true
	})
	return nil
}
func (h *captureHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *captureHandler) WithGroup(string) slog.Handler     { return h }

func TestLogContextTraceID(t *testing.T) {
	ctx := WithTrace(context.Background())
	h := &captureHandler{}
	l := slog.New(h)
	LogContext(l, ctx, "hello")
	found := false
	for _, k := range h.keys {
		if k == "trace_id" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected trace_id in log context")
	}
	_ = io.Discard
}
