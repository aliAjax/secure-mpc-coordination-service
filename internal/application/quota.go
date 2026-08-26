package application

import (
	"context"
	"github.com/example/027-mpc-coordinator/internal/domain"
	"sync"
	"time"
)

type QuotaManager struct {
	mu     sync.Mutex
	limits map[string]domain.ResourceBudget
	used   map[string]domain.ResourceBudget
}

func NewQuotaManager() *QuotaManager {
	return &QuotaManager{limits: map[string]domain.ResourceBudget{}, used: map[string]domain.ResourceBudget{}}
}
func (q *QuotaManager) SetLimit(tenant string, b domain.ResourceBudget) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.limits[tenant] = b
}
func (q *QuotaManager) Reserve(ctx context.Context, tenant string, b domain.ResourceBudget) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	limit := q.limits[tenant]
	used := q.used[tenant]
	if limit.CPUSeconds > 0 && used.CPUSeconds+b.CPUSeconds > limit.CPUSeconds {
		return domain.ErrConflict
	}
	if limit.MemoryBytes > 0 && used.MemoryBytes+b.MemoryBytes > limit.MemoryBytes {
		return domain.ErrConflict
	}
	if limit.NetworkBytes > 0 && used.NetworkBytes+b.NetworkBytes > limit.NetworkBytes {
		return domain.ErrConflict
	}
	used.CPUSeconds += b.CPUSeconds
	used.MemoryBytes += b.MemoryBytes
	used.NetworkBytes += b.NetworkBytes
	q.used[tenant] = used
	return nil
}
func (q *QuotaManager) Release(ctx context.Context, tenant string, b domain.ResourceBudget) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	u := q.used[tenant]
	u.CPUSeconds -= b.CPUSeconds
	u.MemoryBytes -= b.MemoryBytes
	u.NetworkBytes -= b.NetworkBytes
	if u.CPUSeconds < 0 {
		u.CPUSeconds = 0
	}
	if u.MemoryBytes < 0 {
		u.MemoryBytes = 0
	}
	if u.NetworkBytes < 0 {
		u.NetworkBytes = 0
	}
	q.used[tenant] = u
}
func (q *QuotaManager) Snapshot(tenant string) (domain.ResourceBudget, time.Time) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.used[tenant], time.Now().UTC()
}
