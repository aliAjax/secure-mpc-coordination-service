package repository

import (
	"sort"
	"sync"
	"time"
)

type TimeIndex[T any] struct {
	mu    sync.RWMutex
	items []indexed[T]
}
type indexed[T any] struct {
	at    time.Time
	value T
}

func NewTimeIndex[T any]() *TimeIndex[T] { return &TimeIndex[T]{items: []indexed[T]{}} }
func (i *TimeIndex[T]) Add(at time.Time, v T) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.items = append(i.items, indexed[T]{at, v})
	sort.Slice(i.items, func(a, b int) bool { return i.items[a].at.Before(i.items[b].at) })
}
func (i *TimeIndex[T]) Before(until time.Time) []T {
	i.mu.RLock()
	defer i.mu.RUnlock()
	out := []T{}
	for _, v := range i.items {
		if v.at.Before(until) {
			out = append(out, v.value)
		}
	}
	return out
}
func (i *TimeIndex[T]) Len() int { i.mu.RLock(); defer i.mu.RUnlock(); return len(i.items) }
