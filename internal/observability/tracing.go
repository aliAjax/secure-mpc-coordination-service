package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
)

type traceKey struct{}

func WithTrace(ctx context.Context) context.Context {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return context.WithValue(context.Background(), traceKey{}, hex.EncodeToString(b))
}
func TraceID(ctx context.Context) string {
	if v, ok := ctx.Value(traceKey{}).(string); ok {
		return v
	}
	return ""
}
func LogContext(l *slog.Logger, ctx context.Context, msg string, args ...any) {
	l.Info(msg, args...)
}
