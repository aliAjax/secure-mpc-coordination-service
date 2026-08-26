package observability

import "sync/atomic"

type Metrics struct {
	Requests        atomic.Uint64
	Errors          atomic.Uint64
	Shares          atomic.Uint64
	Reconstructions atomic.Uint64
}

func (m *Metrics) IncRequest()        { m.Requests.Add(1) }
func (m *Metrics) IncError()          { m.Errors.Add(1) }
func (m *Metrics) IncShare()          { m.Shares.Add(1) }
func (m *Metrics) IncReconstruction() { m.Reconstructions.Add(1) }
func (m *Metrics) Prometheus() string {
	return "mpc_requests_total " + itoa(m.Requests.Load()) + "\nmpc_errors_total " + itoa(m.Errors.Load()) + "\nmpc_shares_total " + itoa(m.Shares.Load()) + "\nmpc_reconstructions_total " + itoa(m.Reconstructions.Load()) + "\n"
}
func itoa(v uint64) string {
	if v == 0 {
		return "0"
	}
	b := make([]byte, 0, 20)
	for v > 0 {
		b = append([]byte{byte('0' + v%10)}, b...)
		v /= 10
	}
	return string(b)
}
