package domain

type Decision struct {
	Rule   string `json:"rule"`
	Passed bool   `json:"passed"`
	Detail string `json:"detail"`
}
type Explanation struct {
	ComputationID string     `json:"computation_id"`
	Decisions     []Decision `json:"decisions"`
	Safe          bool       `json:"safe"`
}

func Explain(c Computation, r Round) Explanation {
	e := Explanation{ComputationID: c.ID, Safe: true}
	checks := []Decision{{"threshold", c.Threshold >= 2, "threshold must be at least two"}, {"participants", c.ParticipantCount >= c.Threshold, "participant count covers threshold"}, {"deadline", !r.Deadline.IsZero(), "round has deadline"}, {"nonce", r.Nonce != "", "round replay nonce present"}}
	for _, d := range checks {
		if !d.Passed {
			e.Safe = false
		}
		e.Decisions = append(e.Decisions, d)
	}
	return e
}
