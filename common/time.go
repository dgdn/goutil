package common

import (
	"time"
)

var TimeLayout = "2006-01-02 15:04:05"

func TimeStrUnixNano(t int64) string {
	var ti = time.Unix(0, t*int64(time.Nanosecond))
	return ti.Format(TimeLayout)
}
