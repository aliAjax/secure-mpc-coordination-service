package domain

import "strings"

type Capability struct {
	Protocols       []string        `json:"protocols"`
	MaxMemory       int64           `json:"max_memory"`
	MaxParticipants int             `json:"max_participants"`
	Features        map[string]bool `json:"features"`
}

func (c Capability) Supports(protocol string) bool {
	for _, p := range c.Protocols {
		if strings.EqualFold(p, protocol) {
			return true
		}
	}
	return false
}
func (c Capability) CanHandle(n int, memory int64) bool {
	return n > 0 && n <= c.MaxParticipants && (c.MaxMemory == 0 || memory <= c.MaxMemory)
}
func Negotiate(all []Capability, protocol string) Capability {
	out := Capability{Features: map[string]bool{}}
	for _, c := range all {
		if !c.Supports(protocol) {
			continue
		}
		out.MaxParticipants += c.MaxParticipants
		if c.MaxMemory == 0 || (out.MaxMemory > 0 && c.MaxMemory < out.MaxMemory) {
			out.MaxMemory = c.MaxMemory
		}
		out.Protocols = append(out.Protocols, protocol)
		for k, v := range c.Features {
			out.Features[k] = out.Features[k] || v
		}
	}
	return out
}
