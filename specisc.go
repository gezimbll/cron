package cron

import (
	"time"
)

func DaysInMonth(year int, month time.Month) float64 {
	return float64(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1).Day())
}

func (s *SpecSchedule) IsNextAnytime() (isNextAnytime bool) {
	isNextAnytime = s.Month&starBit != 0 &&
		s.Dom&starBit != 0 &&
		s.Hour&starBit != 0 &&
		s.Minute&starBit != 0
	if s.options&SecondOptional > 0 {
		isNextAnytime = isNextAnytime && s.Second&starBit != 0
	}
	if s.options&DowOptional == 0 {
		isNextAnytime = isNextAnytime && s.Dow&starBit != 0
	}
	return
}

// NextInactive returns the next time this schedule is deactivated, greater than the given
// time.  If the schedule is always active in the next 5 years, return the zero time.
func (s *SpecSchedule) NextInactive(t time.Time) time.Time {
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
	if s.IsNextAnytime() {
		return time.Time{}
	}
	origLocation := t.Location()
	loc := s.Location
	if loc == time.Local {
		loc = t.Location()
	}
	if s.Location != time.Local {
		t = t.In(s.Location)
	}
	// Start at the earliest possible time (the upcoming second).
	t = t.Add(time.Second - time.Duration(t.Nanosecond()))

	// First no-match wins

	// Check seconds if they are activated in options
	if s.options&SecondOptional > 0 &&
		s.Second&starBit == 0 {
		tChk := t
		max := int(seconds.max)
		for i := 0; i <= max; i++ {
			if 1<<uint(tChk.Second())&s.Second == 0 { // found it
				return tChk.In(origLocation)
			}
			if i == 0 {
				tChk = tChk.Truncate(time.Second) // round to second
			}
			if i != max {
				tChk = tChk.Add(1 * time.Second)
			}
		}
	}

	// Check minutes
	if s.Minute&starBit == 0 {
		tChk := t
		max := int(minutes.max)
		for i := 0; i <= max; i++ {
			if 1<<uint(tChk.Minute())&s.Minute == 0 {
				return tChk.In(origLocation)
			}
			if i == 0 {
				tChk = tChk.Truncate(time.Minute) // round to minute
			}
			if i != max {
				tChk = tChk.Add(1 * time.Minute)
			}
		}
	}

	// Check hours
	if s.Hour&starBit == 0 {
		tChk := t
		max := int(hours.max)
		for i := 0; i <= max; i++ {
			if 1<<uint(tChk.Hour())&s.Hour == 0 {
				return tChk.In(origLocation)
			}
			if i == 0 {
				tChk = time.Date(tChk.Year(), tChk.Month(), tChk.Day(), tChk.Hour(), 0, 0, 0, loc) // Round to hour
			}
			if i != max {
				tChk = tChk.Add(1 * time.Hour)
			}
		}

	}

	// Now get a day in that month.
	//
	// NOTE: This causes issues for daylight savings regimes where midnight does
	// not exist.  For example: Sao Paulo has DST that transforms midnight on
	// 11/3 into 1am. Handle that by noticing when the Hour ends up != 0.
	if s.Dom&starBit == 0 || s.Dow&starBit == 0 {
		tChk := t
		max := int(dom.max) // cover 31 days every second month
		if DaysInMonth(tChk.Year(), tChk.Month()) != 31 {
			max = max * 2
		}
		for i := 0; i <= max; i++ {
			if !dayMatches(s, tChk) {
				return tChk.In(origLocation)
			}
			if tChk.Day() == 1 && // month has priority over day as shortcut
				s.Month&starBit == 0 && 1<<uint(tChk.Month())&s.Month == 0 {
				return tChk.In(origLocation)
			}

			if i == 0 {
				tChk = time.Date(tChk.Year(), tChk.Month(), tChk.Day(), 0, 0, 0, 0, loc) // Round to hour
			}
			if i != max {
				tChk = tChk.AddDate(0, 0, 1)
				if tChk.Hour() != 0 {
					if tChk.Hour() > 12 { // Notice if the hour is no longer midnight due to DST.
						// Add an hour if it's 23, subtract an hour if it's 1.
						tChk = tChk.Add(time.Duration(24-tChk.Hour()) * time.Hour)
					} else {
						tChk = tChk.Add(time.Duration(-tChk.Hour()) * time.Hour)
					}
				}
			}
		}
	}

	// Find the first non applicable month.
	if s.Month&starBit == 0 {
		tChk := t
		max := int(months.max)
		for i := 0; i <= max; i++ {
			if 1<<uint(tChk.Month())&s.Month == 0 {
				return tChk.In(origLocation)
			}
			if i == 0 {
				tChk = time.Date(tChk.Year(), tChk.Month(), 1, 0, 0, 0, 0, loc)
			}
			if i != max {
				tChk = tChk.AddDate(0, 1, 0)
			}
		}

	}

	return time.Time{}
}
