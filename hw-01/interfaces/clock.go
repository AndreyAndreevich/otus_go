package interfaces

import "time"

type Clock interface {
	GetCurrentTime() (time.Time, error)
}

