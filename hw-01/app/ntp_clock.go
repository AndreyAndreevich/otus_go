package app

import (
	"github.com/beevik/ntp"
	"time"
)

type NtpClock struct {
	host string
}

func NewNtpClock(host string) *NtpClock {
	return &NtpClock{
		host: host,
	}
}

func (clock *NtpClock) GetCurrentTime() (time.Time, error) {
	return ntp.Time(clock.host)
}