package subscription

type State int

const (
	StateNew              State = iota
	StateAddChart         State = iota
	StateAwaitAddChart    State = iota
	StateRemoveChart      State = iota
	StateAwaitRemoveChart State = iota
	StateOnUpdates        State = iota
)

func (s State) IsInputState() bool {
	switch s {
	case StateAwaitAddChart, StateAwaitRemoveChart:
		return true
	default:
		return false
	}
}
