package spanet

import (
	"time"
)

func (s *Spanet) SetPowerSave(mode PowerSaveMode) (PowerSaveMode, error) {
	r, err := s.setMode("W63", byte(mode))
	return PowerSaveMode(r), err
}

func (s *Spanet) SetPeakStart(when time.Time) (time.Time, error) {
	return s.commandTime("W64", when)
}

func (s *Spanet) SetPeakEnd(when time.Time) (time.Time, error) {
	return s.commandTime("W65", when)
}
