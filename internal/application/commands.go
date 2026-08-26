package application

import (
	"context"
	"fmt"
	"github.com/example/027-mpc-coordinator/internal/domain"
	"time"
)

type CommandResult struct {
	Accepted bool
	Message  string
	At       time.Time
}

func (s *Service) ValidateComputation(ctx context.Context, id string) error {
	c, e := s.store.GetComputation(ctx, id)
	if e != nil {
		return e
	}
	if err := domain.ValidateTenant(c.TenantID); err != nil {
		return err
	}
	return domain.ValidateProtocol(c.Protocol, c.ProtocolVersion)
}
func (s *Service) MarkCommitted(ctx context.Context, id string) (CommandResult, error) {
	c, e := s.store.GetComputation(ctx, id)
	if e != nil {
		return CommandResult{}, e
	}
	if e = domain.TransitionComputation(c, domain.StatusCommitted); e != nil {
		return CommandResult{}, e
	}
	c.UpdatedAt = time.Now().UTC()
	if e = s.store.UpdateComputation(ctx, c); e != nil {
		return CommandResult{}, e
	}
	return CommandResult{true, "committed", time.Now().UTC()}, nil
}
func (s *Service) RenewRound(ctx context.Context, rid, owner string, ttl time.Duration) (CommandResult, error) {
	r, e := s.AcquireLease(ctx, rid, owner, ttl)
	if e != nil {
		return CommandResult{}, e
	}
	return CommandResult{true, fmt.Sprintf("lease_until=%s", r.LeaseUntil.Format(time.RFC3339)), time.Now().UTC()}, nil
}
func (s *Service) VerifyOutput(ctx context.Context, id string, out domain.Output) bool {
	c, e := s.store.GetComputation(ctx, id)
	if e != nil || c.Output == nil {
		return false
	}
	return c.Output.Proof == out.Proof && c.Output.Value == out.Value
}
