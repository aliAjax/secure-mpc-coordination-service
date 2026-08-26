package domain

import "fmt"

func TransitionComputation(c *Computation, next ComputationStatus) error {
	allowed := map[ComputationStatus]map[ComputationStatus]bool{
		StatusDraft:          {StatusCommitted: true, StatusAborted: true, StatusExpired: true, StatusRunning: true},
		StatusCommitted:      {StatusRunning: true, StatusAborted: true, StatusExpired: true},
		StatusRunning:        {StatusWaitingShares: true, StatusReconstructing: true, StatusAborted: true, StatusExpired: true},
		StatusWaitingShares:  {StatusAborted: true, StatusExpired: true},
		StatusReconstructing: {StatusSucceeded: true, StatusAborted: true},
		StatusSucceeded:      {}, StatusAborted: {}, StatusExpired: {},
	}
	if !allowed[c.Status][next] {
		return fmt.Errorf("%w: computation %s -> %s", ErrConflict, c.Status, next)
	}
	c.Status = next
	c.Version++
	return nil
}

func TransitionRound(r *Round, next RoundStatus) error {
	allowed := map[RoundStatus]map[RoundStatus]bool{
		RoundOpen:       {RoundCollecting: true, RoundExpired: true, RoundAborted: true},
		RoundCollecting: {RoundExpired: true, RoundAborted: true},
		RoundComplete:   {}, RoundExpired: {}, RoundAborted: {},
	}
	if !allowed[r.Status][next] {
		return fmt.Errorf("%w: round %s -> %s", ErrConflict, r.Status, next)
	}
	r.Status = next
	r.Version++
	return nil
}
