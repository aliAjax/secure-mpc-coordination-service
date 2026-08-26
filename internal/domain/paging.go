package domain

import "sort"

type Cursor struct {
	LastID string
	Limit  int
}

func NormalizeCursor(c Cursor) Cursor {
	if c.Limit < 0 || c.Limit > 500 {
		c.Limit = 100
	}
	return c
}
func PageStrings(values []string, c Cursor) ([]string, Cursor) {
	c = NormalizeCursor(c)
	sort.Strings(values)
	start := 0
	for i, v := range values {
		if v == c.LastID {
			start = i + 1
			break
		}
	}
	end := start + c.Limit
	if end > len(values) {
		end = len(values)
	}
	next := Cursor{Limit: c.Limit}
	if end < len(values) {
		next.LastID = values[end]
	}
	return values[start:end], next
}
