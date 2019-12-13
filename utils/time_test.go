package utils

import (
	"testing"
	"time"
)

func Test_SubTime(t *testing.T) {
	now := time.Now()
	createTime, _ := time.Parse(time.RFC3339Nano, "2019-12-13T01:45:33.000Z")
	subM := now.Sub(createTime)
	t.Log(int(subM.Minutes()), "Hours")
}