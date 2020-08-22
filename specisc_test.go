package cron

import (
	"testing"
	"time"
)

/*
func TestNextNotActive(t *testing.T) {
	fromTime := time.Date(2020, 1, 1, 12, 00, 0, 0, time.UTC)
	sched, err := ParseStandard("* * * 12 *")
	if err != nil {
		t.Error(err)
	}
	if tm := sched.NextNotActive(fromTime); tm != time.Date(2020, 12, 1, 12, 0, 1, 0, time.UTC) {
		t.Errorf("got time: %+v", tm)
	}

}
*/

func BenchmarkCronNext(b *testing.B) {
	now := time.Now()
	for n := 0; n < b.N; n++ {
		sched, _ := ParseStandard("* * * 12 *")
		sched.Next(now)
	}
}
