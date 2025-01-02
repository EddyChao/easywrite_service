package verification_code

import "time"

const (
	CodeTTL          = time.Minute * 10
	SingleDayMaximum = 30
	Interval         = time.Second * 60
	ValidPeriod      = time.Minute * 10
	TryMax           = 3
)
