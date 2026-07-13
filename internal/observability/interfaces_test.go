package observability

// compile-time interface verification
var (
	_ snapshotter = (*Metrics)(nil)
)

// snapshotter is implemented by types that expose a Snapshot method.
type snapshotter interface {
	Snapshot() MetricsSnapshot
}
