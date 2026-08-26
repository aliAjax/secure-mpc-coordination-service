package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
)

type traceKey struct{}

func WithTrace(ctx context.Context) context.Context {
	if existing := TraceID(ctx); existing != "" {
		return ctx
	}
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return context.WithValue(ctx, traceKey{}, hex.EncodeToString(b))
}
func TraceID(ctx context.Context) string {
	if v, ok := ctx.Value(traceKey{}).(string); ok {
		return v
	}
	return ""
}
func LogContext(l *slog.Logger, ctx context.Context, msg string, args ...any) {
	tid := TraceID(ctx)
	if tid != "" {
		args = append([]any{"trace_id", tid}, args...)
	}
	l.Info(msg, args...)
}
