package api

type ErrorResponse struct {
	Error   ErrorBody `json:"error"`
	TraceID string    `json:"trace_id,omitempty"`
}
type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}
type Page[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
	Count      int    `json:"count"`
}
type VerifyRequest struct {
	Value      string `json:"value"`
	Nonce      string `json:"nonce"`
	Commitment string `json:"commitment"`
}
type VerifyResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}
