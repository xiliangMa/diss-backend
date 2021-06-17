package utils

import (
	"testing"
	"time"
)

func Test_SubTime(t *testing.T) {
	now := time.Now()
	createTime, _ := time.Parse(time.RFC3339Nano, "2019-12-13T01:45:33.000Z0700")
	subM := now.Sub(createTime)
	t.Log(int(subM.Minutes()), "Hours")
}

func Test_GetTimeFromNow(t *testing.T) {
	now := time.Now().Format("2006-01-02T15:04:05Z")
	//nowStr, _ := time.Parse(time.RFC3339, now.String())
	timepoint := time.Now().Add(time.Hour * -1).Format("2006-01-02T15:04:05Z")

	t.Log("\n Now:", now, "\n", "Timepoint:", timepoint)
}
