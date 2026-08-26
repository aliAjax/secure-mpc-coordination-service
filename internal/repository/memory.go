package repository

import (
	"context"
	"encoding/json"
	"github.com/example/027-mpc-coordinator/internal/domain"
	"os"
	"sync"
)

type MemoryStore struct {
	mu           sync.RWMutex
	computations map[string]*domain.Computation
	participants map[string]*domain.Participant
	rounds       map[string]*domain.Round
	evidence     map[string][]*domain.Evidence
	file         string
}

func NewMemoryStore(file string) *MemoryStore {
	m := &MemoryStore{computations: map[string]*domain.Computation{}, participants: map[string]*domain.Participant{}, rounds: map[string]*domain.Round{}, evidence: map[string][]*domain.Evidence{}, file: file}
	_ = m.load()
	return m
}
func (m *MemoryStore) CreateComputation(_ context.Context, c *domain.Computation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.computations[c.ID]; ok {
		return domain.ErrConflict
	}
	cp := *c
	m.computations[c.ID] = &cp
	return m.persistLocked()
}
func (m *MemoryStore) GetComputation(_ context.Context, id string) (*domain.Computation, error) {
	c, ok := m.computations[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return c, nil
}
func (m *MemoryStore) UpdateComputation(_ context.Context, c *domain.Computation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.computations[c.ID]; !ok {
		return domain.ErrNotFound
	}
	cp := *c
	m.computations[c.ID] = &cp
	return m.persistLocked()
}
func (m *MemoryStore) ListComputations(_ context.Context) ([]domain.Computation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]domain.Computation, 0, len(m.computations))
	for _, c := range m.computations {
		out = append(out, *c)
	}
	return out, nil
}
func (m *MemoryStore) PutParticipant(_ context.Context, p *domain.Participant) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *p
	m.participants[p.ID] = &cp
	return m.persistLocked()
}
func (m *MemoryStore) ListParticipants(_ context.Context, cid string) ([]domain.Participant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Participant{}
	for _, p := range m.participants {
		if p.ComputationID == cid {
			out = append(out, *p)
		}
	}
	return out, nil
}
func (m *MemoryStore) CreateRound(_ context.Context, r *domain.Round) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.rounds[r.ID]; ok {
		return domain.ErrConflict
	}
	cp := *r
	cp.Shares = cloneShares(r.Shares)
	m.rounds[r.ID] = &cp
	return m.persistLocked()
}
func (m *MemoryStore) GetRound(_ context.Context, id string) (*domain.Round, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rounds[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return r, nil
}
func (m *MemoryStore) UpdateRound(_ context.Context, r *domain.Round) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.rounds[r.ID]; !ok {
		return domain.ErrNotFound
	}
	cp := *r
	cp.Shares = cloneShares(r.Shares)
	m.rounds[r.ID] = &cp
	return m.persistLocked()
}
func cloneShares(in map[string]domain.Share) map[string]domain.Share {
	out := map[string]domain.Share{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
func (m *MemoryStore) PutEvidence(_ context.Context, e *domain.Evidence) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *e
	m.evidence[e.ComputationID] = append(m.evidence[e.ComputationID], &cp)
	return m.persistLocked()
}
func (m *MemoryStore) ListEvidence(_ context.Context, id string) ([]domain.Evidence, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Evidence{}
	for _, e := range m.evidence[id] {
		out = append(out, *e)
	}
	return out, nil
}

type snapshot struct {
	Computations map[string]*domain.Computation `json:"computations"`
	Participants map[string]*domain.Participant `json:"participants"`
	Rounds       map[string]*domain.Round       `json:"rounds"`
	Evidence     map[string][]*domain.Evidence  `json:"evidence"`
}

func (m *MemoryStore) persistLocked() error {
	if m.file == "" {
		return nil
	}
	s := snapshot{m.computations, m.participants, m.rounds, m.evidence}
	b, e := json.MarshalIndent(s, "", "  ")
	if e != nil {
		return e
	}
	tmp := m.file + ".tmp"
	if e = os.WriteFile(tmp, b, 0600); e != nil {
		return e
	}
	return os.Rename(tmp, m.file)
}
func (m *MemoryStore) load() error {
	if m.file == "" {
		return nil
	}
	b, e := os.ReadFile(m.file)
	if e != nil {
		return nil
	}
	var s snapshot
	if json.Unmarshal(b, &s) != nil {
		return nil
	}
	if s.Computations != nil {
		m.computations = s.Computations
	}
	if s.Participants != nil {
		m.participants = s.Participants
	}
	if s.Rounds != nil {
		m.rounds = s.Rounds
	}
	if s.Evidence != nil {
		m.evidence = s.Evidence
	}
	return nil
}
