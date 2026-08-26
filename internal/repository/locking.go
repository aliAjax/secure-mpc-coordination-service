package repository

import "sync"

type LockTable struct {
	mu    sync.Mutex
	locks map[string]string
}

func NewLockTable() *LockTable { return &LockTable{locks: map[string]string{}} }
func (l *LockTable) Acquire(resource, owner string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if v := l.locks[resource]; v != "" && v != owner {
		return false
	}
	l.locks[resource] = owner
	return true
}
func (l *LockTable) Release(resource, owner string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.locks[resource] == owner {
		delete(l.locks, resource)
	}
}
func (l *LockTable) Owner(resource string) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.locks[resource]
}
