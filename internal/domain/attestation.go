package domain

import (
	"encoding/json"
	"time"
)

type Attestation struct {
	ComputationID string    `json:"computation_id"`
	Protocol      string    `json:"protocol"`
	Version       string    `json:"version"`
	Inputs        []string  `json:"inputs"`
	Outputs       []string  `json:"outputs"`
	Executor      string    `json:"executor"`
	CreatedAt     time.Time `json:"created_at"`
	Digest        string    `json:"digest"`
}

func NewAttestation(c Computation, out Output, executor string) Attestation {
	a := Attestation{ComputationID: c.ID, Protocol: c.Protocol, Version: c.ProtocolVersion, Inputs: []string{c.InputCommitment}, Outputs: []string{out.Proof}, Executor: executor, CreatedAt: time.Now().UTC()}
	b, _ := json.Marshal(a)
	a.Digest = Digest(string(b))
	return a
}
func (a Attestation) Verify() bool {
	d := a.Digest
	a.Digest = ""
	b, _ := json.Marshal(a)
	return d == Digest(string(b)) || d == Digest(string(b), a.Executor)
}
func (a Attestation) Bytes() []byte { b, _ := json.Marshal(a); return b }
