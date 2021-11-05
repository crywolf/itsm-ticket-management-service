package incident

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// State values
var (
	StateNew        = State{"new"}
	StatePreOnHold  = State{"pre on hold"}
	StateInProgress = State{"in progress"}
	StateOnHold     = State{"on hold"}
	StateResolved   = State{"resolved"}
	StateClosed     = State{"closed"}
	StateCancelled  = State{"cancelled"}
)

var ticketStateValues = []State{
	StateNew,
	StatePreOnHold,
	StateInProgress,
	StateOnHold,
	StateResolved,
	StateClosed,
	StateCancelled,
}

// State of the ticket is enum.
// swagger:strfmt string
type State struct {
	v string
}

// NewStateFromString creates new instance from string value
func NewStateFromString(stateStr string) (State, error) {
	for _, state := range ticketStateValues {
		if state.String() == stateStr {
			return state, nil
		}
	}
	return State{}, errors.Errorf("unknown '%s' state", stateStr)
}

// IsZero returns true if State has zero value.
// Every type in Go have zero value. In that case it's `State{}`.
// It's always a good idea to check if provided value is not zero!
func (s State) IsZero() bool {
	return s == State{}
}

func (s State) String() string {
	return s.v
}

// MarshalJSON returns Entity as the JSON encoding of Entity
func (s State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
