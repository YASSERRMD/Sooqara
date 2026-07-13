package store

// State represents a job's lifecycle position.
type State string

const (
	StateQueued      State = "queued"
	StateAnalysing   State = "analysing"
	StateCopywriting State = "copywriting"
	StateImaging     State = "imaging"
	StateVideoing    State = "videoing"
	StateAssembling  State = "assembling"
	StateDone        State = "done"
	StateFailed      State = "failed"
	StateCancelled   State = "cancelled"
)

// ValidTransitions maps each state to the states it can transition to.
var ValidTransitions = map[State][]State{
	StateQueued:      {StateAnalysing, StateFailed, StateCancelled},
	StateAnalysing:   {StateCopywriting, StateFailed, StateCancelled},
	StateCopywriting: {StateImaging, StateFailed, StateCancelled},
	StateImaging:     {StateVideoing, StateFailed, StateCancelled},
	StateVideoing:    {StateAssembling, StateFailed, StateCancelled},
	StateAssembling:  {StateDone, StateFailed, StateCancelled},
	StateDone:        {},
	StateFailed:      {},
	StateCancelled:   {},
}

// NextState returns the state after a stage completes successfully.
func NextState(s State) State {
	switch s {
	case StateQueued:
		return StateAnalysing
	case StateAnalysing:
		return StateCopywriting
	case StateCopywriting:
		return StateImaging
	case StateImaging:
		return StateVideoing
	case StateVideoing:
		return StateAssembling
	case StateAssembling:
		return StateDone
	default:
		return s
	}
}
