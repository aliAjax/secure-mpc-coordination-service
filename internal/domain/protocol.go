package domain

import (
	"context"
	"fmt"
	"time"
)

type ProtocolSpec struct {
	Name            string
	Version         string
	MinParticipants int
	MaxParticipants int
	Resource        ResourceBudget
}
type ResourceBudget struct {
	CPUSeconds   int64 `json:"cpu_seconds"`
	MemoryBytes  int64 `json:"memory_bytes"`
	NetworkBytes int64 `json:"network_bytes"`
}
type ProtocolRegistry struct{ specs map[string]ProtocolSpec }

func NewProtocolRegistry() *ProtocolRegistry {
	return &ProtocolRegistry{specs: map[string]ProtocolSpec{"sum": {Name: "sum", Version: "1", MinParticipants: 2, MaxParticipants: 100, Resource: ResourceBudget{CPUSeconds: 60, MemoryBytes: 1 << 30, NetworkBytes: 1 << 30}}, "average": {Name: "average", Version: "1", MinParticipants: 2, MaxParticipants: 100, Resource: ResourceBudget{CPUSeconds: 120, MemoryBytes: 1 << 30, NetworkBytes: 1 << 30}}}}
}
func (r *ProtocolRegistry) Register(s ProtocolSpec) error {
	if s.Name == "" || s.Version == "" || s.MinParticipants < 2 || s.MaxParticipants < s.MinParticipants {
		return ErrInvalid
	}
	r.specs[s.Name+":"+s.Version] = s
	return nil
}
func (r *ProtocolRegistry) Resolve(name, version string) (ProtocolSpec, error) {
	s, ok := r.specs[name+":"+version]
	if !ok {
		return ProtocolSpec{}, fmt.Errorf("%w: protocol %s/%s", ErrNotFound, name, version)
	}
	return s, nil
}

type ProtocolEngine interface {
	Validate(context.Context, ProtocolSpec, Computation) error
	Execute(context.Context, ProtocolSpec, []Share) (*Output, error)
}
type SumEngine struct{}

func (SumEngine) Validate(_ context.Context, s ProtocolSpec, c Computation) error {
	if c.ParticipantCount < s.MinParticipants || c.ParticipantCount > s.MaxParticipants {
		return ErrInvalid
	}
	return nil
}
func (SumEngine) Execute(ctx context.Context, _ ProtocolSpec, shares []Share) (*Output, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if len(shares) == 0 {
		return nil, ErrThreshold
	}
	return &Output{Value: shares[0].Value, ParticipantCount: len(shares), ProducedAt: time.Now().UTC()}, nil
}
