package cron

import (
	"fmt"
	"time"
)

// NextNotActive returns the next time this schedule is deactivated, greater than the given
// time.  If the schedule is always active in the next 5 years, return the zero time.
func (s *SpecSchedule) NextNotActive(t time.Time) (tOut time.Time) {
	// General approach
	//
	// This implementation is the mirror of the original Next implementation. Only the comparison operator is changed.
	// For Month, Day, Hour, Minute, Second:
	// Check if the time value matches.  If not, continue to the next field.
	// If the field matches the schedule, then increment the field until it doesn't.
	// While incrementing the field, a wrap-around brings it back to the beginning
	// of the field list (since it is necessary to re-verify previous field
	// values)

	// Convert the given time into the schedule's timezone, if one is specified.
	// Save the original timezone so we can convert back after we find a time.
	// Note that schedules without a time zone specified (time.Local) are treated
	// as local to the time provided.
	origLocation := t.Location()
	loc := s.Location
	if loc == time.Local {
		loc = t.Location()
	}
	if s.Location != time.Local {
		t = t.In(s.Location)
	}
	// Start at the earliest possible time (the upcoming second).
	t = t.Add(1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)

	// This flag indicates whether a field has been incremented.
	var added bool

	// If no time is found within five years, return zero.
	yearLimit := t.Year() + 5
	notMatching := make(map[ParseOption]int)

WRAP:
	if t.Year() > yearLimit {
		if tOut.IsZero() {
			return
		}
		return tOut.In(origLocation)
	}
	// Find the first non applicable month.
	// If it's this month, then do nothing.

	for 1<<uint(t.Month())&s.Month != 0 {
		fmt.Printf("checking month: %+v, flag: %d against s.Month: %+v with result: %+v\n", t.Month(), 1<<uint(t.Month()), s.Month, 1<<uint(t.Month())&s.Month)
		// If we have to add a month, reset the other parts to 0.
		if !added {
			added = true
			// Otherwise, set the date at the beginning (since the current time is irrelevant).
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
		}
		t = t.AddDate(0, 1, 0)
		// Wrapped around.
		if t.Month() == time.January {
			goto WRAP
		}
	}
	if _, has := notMatching[Month]; !has {
		notMatching[Month] = int(t.Month())
		notMatching[Year] = t.Year()
	}

	// Now get a day in that month.
	//
	// NOTE: This causes issues for daylight savings regimes where midnight does
	// not exist.  For example: Sao Paulo has DST that transforms midnight on
	// 11/3 into 1am. Handle that by noticing when the Hour ends up != 0.
	for dayMatches(s, t) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		}
		t = t.AddDate(0, 0, 1)
		// Notice if the hour is no longer midnight due to DST.
		// Add an hour if it's 23, subtract an hour if it's 1.
		if t.Hour() != 0 {
			if t.Hour() > 12 {
				t = t.Add(time.Duration(24-t.Hour()) * time.Hour)
			} else {
				t = t.Add(time.Duration(-t.Hour()) * time.Hour)
			}
		}

		if t.Day() == 1 {
			goto WRAP
		}
	}
	if _, has := notMatching[Dom]; !has {
		notMatching[Dom] = int(t.Month())
	}

	for 1<<uint(t.Hour())&s.Hour != 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
		}
		t = t.Add(1 * time.Hour)

		if t.Hour() == 0 {
			goto WRAP
		}
	}
	tOut = t

	for 1<<uint(t.Minute())&s.Minute != 0 {
		if !added {
			added = true
			t = t.Truncate(time.Minute)
		}
		t = t.Add(1 * time.Minute)

		if t.Minute() == 0 {
			goto WRAP
		}
	}
	tOut = t

	for 1<<uint(t.Second())&s.Second != 0 {
		if !added {
			added = true
			t = t.Truncate(time.Second)
		}
		t = t.Add(1 * time.Second)

		if t.Second() == 0 {
			goto WRAP
		}
	}
	tOut = t

	return tOut.In(origLocation)
}
