package cron

import (
	"testing"
	"time"
)

func TestNextInactive(t *testing.T) {
	fromTime := time.Date(2020, 1, 1, 12, 00, 0, 0, time.UTC)
	sched, err := ParseStandard("* * * 12 *")
	if err != nil {
		t.Error(err)
	}
	if tm := sched.NextInactive(fromTime); tm != time.Date(2020, 1, 1, 12, 0, 1, 0, time.UTC) {
		t.Errorf("got time: %+v", tm)
	}
	fromTime = time.Date(2020, 11, 30, 23, 59, 59, 0, time.UTC)
	if tm := sched.NextInactive(fromTime); tm != time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) {
		t.Errorf("got time: %+v", tm)
	}
	fromTime = time.Date(2020, 12, 12, 12, 12, 12, 0, time.UTC)
	if tm := sched.NextInactive(fromTime); tm != time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) {
		t.Errorf("got time: %+v", tm)
	}
	fromTime = time.Date(2020, 12, 24, 12, 0, 0, 0, time.UTC)
	if sched, err = ParseStandard("* * 24 12 *"); err != nil {
		t.Error(err)
	}
	if tm := sched.NextInactive(fromTime); tm != time.Date(2020, 12, 25, 0, 0, 0, 0, time.UTC) {
		t.Errorf("got time: %+v", tm)
	}
	fromTime = time.Date(2020, 1, 31, 23, 59, 59, 0, time.UTC)
	if sched, err = ParseStandard("* * 1-30 1,3-12 *"); err != nil {
		t.Error(err)
	}
	if tm := sched.NextInactive(fromTime); tm != time.Date(2020, 02, 1, 0, 0, 0, 0, time.UTC) {
		t.Errorf("got time: %+v", tm)
	}
	fromTime = time.Date(2020, 1, 31, 23, 59, 59, 0, time.UTC)
	if sched, err = ParseStandard("* * 1-30 * *"); err != nil {
		t.Error(err)
	}
	if tm := sched.NextInactive(fromTime); tm != time.Date(2020, 03, 31, 0, 0, 0, 0, time.UTC) {
		t.Errorf("got time: %+v", tm)
	}
	if sched, err = ParseStandard("* * * * *"); err != nil {
		t.Error(err)
	}
	eTm := time.Time{}
	if tm := sched.NextInactive(fromTime); tm != eTm {
		t.Errorf("got time: %+v", tm)
	}
}

// BenchmarkCronNext-8           	  742411	      1596 ns/op
func BenchmarkCronNext(b *testing.B) {
	now := time.Now()
	for n := 0; n < b.N; n++ {
		sched, _ := ParseStandard("* * * 12 *")
		sched.Next(now)
	}
}

// BenchmarkCronNextInactive-8   	  692031	      1599 ns/op
func BenchmarkCronNextInactive(b *testing.B) {
	now := time.Now()
	for n := 0; n < b.N; n++ {
		sched, _ := ParseStandard("* * * 12 *")
		sched.NextInactive(now)
	}
}

func TestNextInactive2(t *testing.T) {
	sched, err := ParseStandard("* * 24 12 *")
	if err != nil {
		t.Error(err)
	}
    exp := time.Date(2020, 12, 25, 0, 0, 0, 0, time.UTC)

	fromTime := time.Date(2020, 12, 24, 12, 59, 58, 0, time.UTC)
	if tm := sched.NextInactive(fromTime); tm != exp{
		t.Errorf("Expected %+v, received %+v", exp, tm)
	}

	fromTime = time.Date(2020, 12, 24, 12, 59, 59, 0, time.UTC)
	if tm := sched.NextInactive(fromTime); tm != exp{
		t.Errorf("Expected %+v, received %+v", exp, tm)
	}

	fromTime = time.Date(2020, 12, 24, 17, 0, 0, 0, time.UTC)
	if tm := sched.NextInactive(fromTime); tm != exp{
		t.Errorf("Expected %+v, received %+v", exp, tm)
	}
}



