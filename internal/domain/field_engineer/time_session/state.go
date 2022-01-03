package tsession

import (
	"fmt"
)

// State values
var (
	StateTravel     = State{"travel"}
	StateWork       = State{"work"}
	StateBreak      = State{"break"}
	StateTravelBack = State{"travel back"}
	StateClosed     = State{"closed"}
)

var timeSessionStateValues = []State{
	StateTravel,
	StateWork,
	StateBreak,
	StateTravelBack,
	StateClosed,
}

// State of the time session is enum.
// swagger:strfmt string
type State struct {
	a string
}

// NewStateFromString creates new instance from string value
func NewStateFromString(stateStr string) (State, error) {
	for _, state := range timeSessionStateValues {
		if state.String() == stateStr {
			return state, nil
		}
	}
	return State{}, fmt.Errorf("unknown '%s' time session state", stateStr)
}

// IsZero returns true if State has zero value.
// Every type in Go have zero value. In that case it's `State{}`.
// It's always a good idea to check if provided value is not zero!
func (s State) IsZero() bool {
	return s == State{}
}

func (s State) String() string {
	return s.a
}
