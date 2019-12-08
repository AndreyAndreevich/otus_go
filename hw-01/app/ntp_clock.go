package app

import (
	"time"

	"github.com/beevik/ntp"
)

// NtpClock gets realtime clock
type NtpClock struct {
	host string
}

// NewNtpClock creates new NtpClock
func NewNtpClock(host string) *NtpClock {
	return &NtpClock{
		host: host,
	}
}

// GetCurrentTime gets current timr from ntp server
func (clock *NtpClock) GetCurrentTime() (time.Time, error) {
	return ntp.Time(clock.host)
}
