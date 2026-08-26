package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

type ComputationStatus string

const (
	StatusDraft          ComputationStatus = "draft"
	StatusCommitted      ComputationStatus = "committed"
	StatusRunning        ComputationStatus = "running"
	StatusWaitingShares  ComputationStatus = "waiting_shares"
	StatusReconstructing ComputationStatus = "reconstructing"
	StatusSucceeded      ComputationStatus = "succeeded"
	StatusAborted        ComputationStatus = "aborted"
	StatusExpired        ComputationStatus = "expired"
)

type RoundStatus string

const (
	RoundOpen       RoundStatus = "open"
	RoundCollecting RoundStatus = "collecting"
	RoundComplete   RoundStatus = "complete"
	RoundExpired    RoundStatus = "expired"
	RoundAborted    RoundStatus = "aborted"
)

type Computation struct {
	ID               string            `json:"id"`
	TenantID         string            `json:"tenant_id"`
	Protocol         string            `json:"protocol"`
	ProtocolVersion  string            `json:"protocol_version"`
	Threshold        int               `json:"threshold"`
	ParticipantCount int               `json:"participant_count"`
	Status           ComputationStatus `json:"status"`
	InputCommitment  string            `json:"input_commitment"`
	Output           *Output           `json:"output,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Version          int64             `json:"version"`
}

type Participant struct {
	ID            string    `json:"id"`
	ComputationID string    `json:"computation_id"`
	Name          string    `json:"name"`
	Capability    string    `json:"capability"`
	PublicKey     string    `json:"public_key"`
	Active        bool      `json:"active"`
	RegisteredAt  time.Time `json:"registered_at"`
}

type Share struct {
	ParticipantID string    `json:"participant_id"`
	Index         int       `json:"index"`
	Value         string    `json:"value"`
	Commitment    string    `json:"commitment"`
	Signature     string    `json:"signature"`
	SubmittedAt   time.Time `json:"submitted_at"`
}

type Round struct {
	ID            string           `json:"id"`
	ComputationID string           `json:"computation_id"`
	Number        int              `json:"number"`
	Nonce         string           `json:"nonce"`
	Status        RoundStatus      `json:"status"`
	Deadline      time.Time        `json:"deadline"`
	LeaseOwner    string           `json:"lease_owner"`
	LeaseUntil    time.Time        `json:"lease_until"`
	Shares        map[string]Share `json:"shares"`
	Version       int64            `json:"version"`
}

type Output struct {
	Value            string    `json:"value"`
	Proof            string    `json:"proof"`
	ParticipantCount int       `json:"participant_count"`
	ProducedAt       time.Time `json:"produced_at"`
}

type Evidence struct {
	ID            string    `json:"id"`
	ComputationID string    `json:"computation_id"`
	Kind          string    `json:"kind"`
	Digest        string    `json:"digest"`
	Payload       string    `json:"payload"`
	CreatedAt     time.Time `json:"created_at"`
}

func NormalizeID(v string) string { return strings.TrimSpace(strings.ToLower(v)) }

func Digest(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

func (c Computation) CanStart() bool { return c.Status == StatusCommitted || c.Status == StatusDraft }
func (c Computation) CanAbort() bool {
	return c.Status != StatusSucceeded && c.Status != StatusAborted && c.Status != StatusExpired
}
