package model

import (
	"time"
)

type Timestamp int64

func Now() Timestamp {
	return FromTime(time.Now())
}

func FromTime(t time.Time) Timestamp {
	return Timestamp(t.Unix())
}

func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}
