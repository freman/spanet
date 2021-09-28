package spanet

import (
	"fmt"
	"strconv"
	"time"
)

func (s *Spanet) SetPowerSave(mode PowerSaveMode) (PowerSaveMode, error) {
	r, err := s.command(fmt.Sprintf("W63:%d", mode))
	if err != nil {
		return 0, err
	}

	_ = r
	rs := ""

	tmp, err := strconv.ParseInt(rs, 10, 64)
	if err != nil {
		return 0, err
	}

	return PowerSaveMode(tmp), nil
}

func (s *Spanet) SetPeakStart(when time.Time) (time.Time, error) {
	return s.commandTime("W64", when)
}

func (s *Spanet) SetPeakEnd(when time.Time) (time.Time, error) {
	return s.commandTime("W65", when)
}
