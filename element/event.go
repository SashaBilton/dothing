//elements contains all the Item level strcutures needed for DoThing
package element

import (
	"time"
)

//An Event is a named and date stamped process that has been applied to an Item
type Event struct {
	EventType string
	Stamp     time.Time
}

//Returns true if an event type exists in a collection of events
func Is(events []Event, ofType string) bool {
	for _, event := range events {
		if event.EventType == ofType {
			return true
		}
	}
	return false
}
