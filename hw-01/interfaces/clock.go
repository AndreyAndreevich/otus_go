package interfaces

import "time"

// Clock is interface wor get current time from somewhere
type Clock interface {
	GetCurrentTime() (time.Time, error)
}
