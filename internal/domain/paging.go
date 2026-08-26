package domain

import "sort"

type Cursor struct {
	LastID string
	Limit  int
}

func NormalizeCursor(c Cursor) Cursor {
	if c.Limit <= 0 || c.Limit > 500 {
		c.Limit = 100
	}
	return c
}
func PageStrings(values []string, c Cursor) ([]string, Cursor) {
	c = NormalizeCursor(c)
	sorted := make([]string, len(values))
	copy(sorted, values)
	sort.Strings(sorted)
	start := 0
	for i, v := range sorted {
		if v == c.LastID {
			start = i + 1
			break
		}
	}
	end := start + c.Limit
	if end > len(sorted) {
		end = len(sorted)
	}
	next := Cursor{Limit: c.Limit}
	if end < len(sorted) {
		next.LastID = sorted[end]
	}
	return sorted[start:end], next
}
