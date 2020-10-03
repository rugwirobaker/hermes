package helmes

// Health report
type Health struct {
	GitRev     string  `json:"git_rev"`
	Uptime     float64 `json:"uptime"`
	Goroutines int     `json:"goroutines"`
}
