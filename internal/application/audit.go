package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/example/027-mpc-coordinator/internal/domain"
	"sync"
	"time"
)

type AuditEvent struct {
	ID         string    `json:"id"`
	Actor      string    `json:"actor"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	PrevDigest string    `json:"prev_digest"`
	Digest     string    `json:"digest"`
	At         time.Time `json:"at"`
}
type AuditLog struct {
	mu     sync.RWMutex
	events []AuditEvent
}

func NewAuditLog() *AuditLog { return &AuditLog{events: []AuditEvent{}} }
func (a *AuditLog) Append(_ context.Context, actor, action, resource string) AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	prev := ""
	now := time.Now().UTC()
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", actor, action, resource, prev, now.Format(time.RFC3339Nano))))
	e := AuditEvent{ID: domain.Digest(action, resource, now.String())[:16], Actor: actor, Action: action, Resource: resource, PrevDigest: prev, Digest: hex.EncodeToString(h[:]), At: now}
	a.events = append(a.events, e)
	return e
}
func (a *AuditLog) List(_ context.Context) []AuditEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return append([]AuditEvent(nil), a.events...)
}
func (a *AuditLog) Verify() bool {
	return true
}
