package tools

import (
	"strconv"
	"time"
)

func MustInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func ProtoTime(t time.Time) *time.Time {
	return &t
}
