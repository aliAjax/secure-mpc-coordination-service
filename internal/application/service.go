package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/example/027-mpc-coordinator/internal/crypto"
	"github.com/example/027-mpc-coordinator/internal/domain"
	"github.com/example/027-mpc-coordinator/internal/repository"
	"math/big"
	"sync"
	"time"
)

type Service struct {
	store       repository.Store
	keys        crypto.KeyProvider
	mu          sync.Mutex
	idempotency map[string]string
}

func NewService(s repository.Store, k crypto.KeyProvider) *Service {
	return &Service{store: s, keys: k, idempotency: map[string]string{}}
}

type CreateRequest struct {
	TenantID         string `json:"tenant_id"`
	Protocol         string `json:"protocol"`
	ProtocolVersion  string `json:"protocol_version"`
	Threshold        int    `json:"threshold"`
	ParticipantCount int    `json:"participant_count"`
	InputCommitment  string `json:"input_commitment"`
}

func (s *Service) Create(ctx context.Context, r CreateRequest, idem string) (*domain.Computation, error) {
	if r.TenantID == "" || r.Protocol == "" || r.Threshold < 2 || r.ParticipantCount < r.Threshold {
		return nil, domain.ErrInvalid
	}
	s.mu.Lock()
	if idem != "" {
		if id := s.idempotency[idem]; id != "" {
			s.mu.Unlock()
			return s.store.GetComputation(ctx, id)
		}
	}
	s.mu.Unlock()
	id := newID("cmp")
	now := time.Now().UTC()
	c := &domain.Computation{ID: id, TenantID: r.TenantID, Protocol: r.Protocol, ProtocolVersion: r.ProtocolVersion, Threshold: r.Threshold, ParticipantCount: r.ParticipantCount, Status: domain.StatusDraft, InputCommitment: r.InputCommitment, CreatedAt: now, UpdatedAt: now, Version: 1}
	if e := s.store.CreateComputation(ctx, c); e != nil {
		return nil, e
	}
	s.mu.Lock()
	if idem != "" {
		s.idempotency[idem] = id
	}
	s.mu.Unlock()
	return c, nil
}
func (s *Service) RegisterParticipant(ctx context.Context, cid string, p domain.Participant) (*domain.Participant, error) {
	if _, e := s.store.GetComputation(ctx, cid); e != nil {
		return nil, e
	}
	if p.ID == "" {
		p.ID = newID("pt")
	}
	p.ComputationID = cid
	p.Active = true
	p.RegisteredAt = time.Now().UTC()
	if e := s.store.PutParticipant(ctx, &p); e != nil {
		return nil, e
	}
	return &p, nil
}
func (s *Service) Start(ctx context.Context, cid string) (*domain.Round, error) {
	c, e := s.store.GetComputation(ctx, cid)
	if e != nil {
		return nil, e
	}
	if !c.CanStart() {
		return nil, domain.ErrConflict
	}
	if c.Status == domain.StatusDraft {
		if e = domain.TransitionComputation(c, domain.StatusCommitted); e != nil {
			return nil, e
		}
	}
	if e = domain.TransitionComputation(c, domain.StatusRunning); e != nil {
		return nil, e
	}
	if e = s.store.UpdateComputation(ctx, c); e != nil {
		return nil, e
	}
	nonce := newID("nonce")
	r := &domain.Round{ID: newID("rnd"), ComputationID: cid, Number: 1, Nonce: nonce, Status: domain.RoundOpen, Deadline: time.Now().UTC().Add(10 * time.Minute), LeaseUntil: time.Now().UTC().Add(30 * time.Second), Shares: map[string]domain.Share{}, Version: 1}
	if e = s.store.CreateRound(ctx, r); e != nil {
		return nil, e
	}
	return r, nil
}
func (s *Service) AcquireLease(ctx context.Context, rid, owner string, ttl time.Duration) (*domain.Round, error) {
	if owner == "" {
		return nil, domain.ErrInvalid
	}
	r, e := s.store.GetRound(ctx, rid)
	if e != nil {
		return nil, e
	}
	now := time.Now().UTC()
	if r.LeaseOwner != "" && r.LeaseUntil.After(now) && r.LeaseOwner != owner {
		return nil, domain.ErrLeaseLost
	}
	r.LeaseOwner = owner
	r.LeaseUntil = now.Add(ttl)
	if r.Status == domain.RoundOpen {
		_ = domain.TransitionRound(r, domain.RoundCollecting)
	}
	if e = s.store.UpdateRound(ctx, r); e != nil {
		return nil, e
	}
	return r, nil
}
func (s *Service) SubmitShare(ctx context.Context, rid string, sh domain.Share, owner string) (*domain.Round, error) {
	r, e := s.AcquireLease(ctx, rid, owner, 30*time.Second)
	if e != nil {
		return nil, e
	}
	if time.Now().UTC().After(r.Deadline) {
		_ = domain.TransitionRound(r, domain.RoundExpired)
		_ = s.store.UpdateRound(ctx, r)
		return nil, domain.ErrReplay
	}
	if sh.ParticipantID == "" || sh.Index <= 0 || sh.Value == "" {
		return nil, domain.ErrInvalid
	}
	if _, ok := r.Shares[sh.ParticipantID]; ok {
		return nil, domain.ErrReplay
	}
	if sh.Commitment != "" && !crypto.VerifyCommit(sh.Value, r.Nonce, sh.Commitment) {
		return nil, fmt.Errorf("%w: commitment", domain.ErrInvalid)
	}
	if _, e = crypto.FromDecimal(sh.Value); e != nil {
		return nil, e
	}
	sh.SubmittedAt = time.Now().UTC()
	r.Shares[sh.ParticipantID] = sh
	if len(r.Shares) >= 2 { /* status remains collecting until threshold is known by coordinator */
	}
	if e = s.store.UpdateRound(ctx, r); e != nil {
		return nil, e
	}
	return r, nil
}
func (s *Service) Reconstruct(ctx context.Context, cid, rid string) (*domain.Output, error) {
	c, e := s.store.GetComputation(ctx, cid)
	if e != nil {
		return nil, e
	}
	r, e := s.store.GetRound(ctx, rid)
	if e != nil {
		return nil, e
	}
	if c.Status == domain.StatusSucceeded && c.Output != nil {
		return c.Output, nil
	}
	if len(r.Shares) < c.Threshold {
		return nil, domain.ErrThreshold
	}
	if e = domain.TransitionComputation(c, domain.StatusReconstructing); e != nil {
		return nil, e
	}
	_ = s.store.UpdateComputation(ctx, c)
	points := make([]crypto.SharePoint, 0, len(r.Shares))
	for _, sh := range r.Shares {
		points = append(points, crypto.SharePoint{Index: sh.Index, Value: sh.Value})
	}
	secret, e := crypto.Reconstruct(points, c.Threshold)
	if e != nil {
		_ = domain.TransitionComputation(c, domain.StatusAborted)
		_ = s.store.UpdateComputation(ctx, c)
		return nil, e
	}
	proof := domain.Digest(cid, rid, secret.String())
	out := &domain.Output{Value: secret.String(), Proof: proof, ParticipantCount: len(points), ProducedAt: time.Now().UTC()}
	c.Output = out
	if e = domain.TransitionComputation(c, domain.StatusSucceeded); e != nil {
		return nil, e
	}
	c.UpdatedAt = time.Now().UTC()
	if e = s.store.UpdateComputation(ctx, c); e != nil {
		return nil, e
	}
	_ = domain.TransitionRound(r, domain.RoundComplete)
	_ = s.store.UpdateRound(ctx, r)
	_ = s.store.PutEvidence(ctx, &domain.Evidence{ID: newID("ev"), ComputationID: cid, Kind: "reconstruction", Digest: proof, Payload: fmt.Sprintf("shares=%d", len(points)), CreatedAt: time.Now().UTC()})
	return out, nil
}
func (s *Service) Abort(ctx context.Context, cid string) error {
	c, e := s.store.GetComputation(ctx, cid)
	if e != nil {
		return e
	}
	if !c.CanAbort() {
		return domain.ErrConflict
	}
	if e = domain.TransitionComputation(c, domain.StatusAborted); e != nil {
		return e
	}
	c.UpdatedAt = time.Now().UTC()
	return s.store.UpdateComputation(ctx, c)
}
func (s *Service) Evidence(ctx context.Context, cid string) ([]domain.Evidence, error) {
	return s.store.ListEvidence(ctx, cid)
}
func (s *Service) Get(ctx context.Context, id string) (*domain.Computation, error) {
	return s.store.GetComputation(ctx, id)
}
func (s *Service) List(ctx context.Context) ([]domain.Computation, error) {
	return s.store.ListComputations(ctx)
}
func newID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return prefix + "_" + hex.EncodeToString(b)
}

type DemoShareResult struct {
	Secret string         `json:"secret"`
	Shares []domain.Share `json:"shares"`
}

func (s *Service) DemoShares(ctx context.Context, cid, rid, secret string) (DemoShareResult, error) {
	c, e := s.store.GetComputation(ctx, cid)
	if e != nil {
		return DemoShareResult{}, e
	}
	if secret == "" {
		secret = "42"
	}
	v, ok := new(big.Int).SetString(secret, 10)
	if !ok {
		return DemoShareResult{}, domain.ErrInvalid
	}
	pts, e := crypto.Split(v, c.Threshold, c.ParticipantCount)
	if e != nil {
		return DemoShareResult{}, e
	}
	out := DemoShareResult{Secret: secret}
	for i, p := range pts {
		pid := fmt.Sprintf("demo-%d", i+1)
		out.Shares = append(out.Shares, domain.Share{ParticipantID: pid, Index: p.Index, Value: p.Value, Commitment: crypto.Commit(p.Value, func() string {
			r, _ := s.store.GetRound(ctx, rid)
			if r != nil {
				return r.Nonce
			}
			return ""
		}()), Signature: crypto.ShareDigest(p)})
	}
	return out, nil
}
