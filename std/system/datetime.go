package system

import "time"

func newDateTime() *DateTime {
	return &DateTime{}
}

type DateTime struct {
}

func (dt *DateTime) Format(format string) string {
	return time.Now().Format(format)
}

func (dt *DateTime) GetTimestamp() int64 {
	return time.Now().Unix()
}
