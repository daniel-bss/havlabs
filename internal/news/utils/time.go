package utils

import (
	"time"
)

func TimeFromString(s string, format string) time.Time {
	date, _ := time.Parse(format, s)
	return date
}
