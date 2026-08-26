package domain

import "time"

type RetentionPolicy struct {
	ID          string        `json:"id"`
	TenantID    string        `json:"tenant_id"`
	ShareTTL    time.Duration `json:"share_ttl"`
	EvidenceTTL time.Duration `json:"evidence_ttl"`
	DeleteAfter time.Time     `json:"delete_after"`
}

func (p RetentionPolicy) Expired(now time.Time) bool {
	return !p.DeleteAfter.IsZero() && now.After(p.DeleteAfter)
}
func (p RetentionPolicy) ShareExpired(created, now time.Time) bool {
	return p.ShareTTL > 0 && now.Sub(created) > p.ShareTTL
}
