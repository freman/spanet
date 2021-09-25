package spanet

import "time"

// SetYear on the controller
// Valid: 1970-20237
// Command: S01:xxxx
func (s *Spanet) SetYear(year int) (int, error) {
	return s.commandInt("S01", year, 1970, 2037, "year", "%04d")
}

// SetMonth on the controller
// Valid: 1-12
// Command: S02:xx
func (s *Spanet) SetMonth(month int) (int, error) {
	return s.commandInt("S02", month, 1, 12, "month", "%02d")
}

// SetDay on the controller
// Valid: 1-13
// Command: S03:xx
func (s *Spanet) SetDay(day int) (int, error) {
	return s.commandInt("S03", day, 1, 31, "day", "%02d")
}

// SetHour on the controller
// Valid: 0-23
// Command: S04:xx
func (s *Spanet) SetHour(hour int) (int, error) {
	return s.commandInt("S04", hour, 0, 23, "hour", "%02d")
}

// SetMinute on the controller
// Valid: 0-59
// Command: S05:xx
func (s *Spanet) SetMinute(minute int) (int, error) {
	return s.commandInt("S05", minute, 0, 59, "minute", "%02d")
}

// SetDateTime on the controller
// This wrapper function will call all the appropriate functions
// on the controller in the correct order to set the date and time
// to a valid value.
func (s *Spanet) SetDateTime(when time.Time) error {
	calls := []func(int) (int, error){
		s.SetYear, s.SetMonth, s.SetDay, s.SetHour, s.SetMinute,
	}
	values := []int{
		when.Year(), int(when.Month()), when.Day(), when.Hour(), when.Minute(),
	}

	for i, v := range values {
		if _, err := calls[i](v); err != nil {
			return err
		}
	}

	return nil
}
