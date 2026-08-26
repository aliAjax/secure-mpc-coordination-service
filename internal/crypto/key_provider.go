package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type KeyProvider interface {
	Key(id string) ([]byte, error)
	Rotate(id string) error
}
type SimulatedKeyProvider struct {
	mu   sync.RWMutex
	keys map[string][]byte
}

func NewKeyProvider() *SimulatedKeyProvider { return &SimulatedKeyProvider{keys: map[string][]byte{}} }
func (p *SimulatedKeyProvider) Key(id string) ([]byte, error) {
	p.mu.RLock()
	k := p.keys[id]
	p.mu.RUnlock()
	if k != nil {
		return append([]byte(nil), k...), nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if existing := p.keys[id]; existing != nil {
		return append([]byte(nil), existing...), nil
	}
	k = make([]byte, 32)
	if _, e := rand.Read(k); e != nil {
		return nil, e
	}
	p.keys[id] = k
	return append([]byte(nil), k...), nil
}
func (p *SimulatedKeyProvider) Rotate(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	k := make([]byte, 32)
	if _, e := rand.Read(k); e != nil {
		return e
	}
	p.keys[id] = k
	return nil
}
func KeyID(k []byte) string {
	if len(k) >= 8 {
		return hex.EncodeToString(k[:8])
	}
	return hex.EncodeToString(k)
}
