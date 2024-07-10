package spanet

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (s *Spanet) SetTargetTemperature(target float64) (float64, error) {
	temp := int(target * 10)
	tmp, err := s.commandInt("W40", temp, 50, 410, "temperature")
	if err != nil {
		if merr, isa := err.(ErrValueOutOfRange); isa {
			merr.Max = merr.Max / 10
			merr.Min = merr.Min / 10
			return 0, merr
		}
		return 0, err
	}

	return float64(tmp) / 10.0, nil
}

func (s *Spanet) ToggleSanitise() error {
	_, err := s.commandExpect("W12", "W12")
	return err
}

func (s *Spanet) SetOperationMode(mode OperationMode) (OperationMode, error) {
	newMode, err := s.setMode("W66", byte(mode))
	return OperationMode(newMode), err
}

func (s *Spanet) SetFiltrationRunTime(hours int) (int, error) {
	return s.commandInt("W60", hours, 1, 24, "hours")
}

func (s *Spanet) SetFiltrationCycle(hours int) (int, error) {
	if !(hours >= 1 || hours <= 4 || hours == 6 || hours == 8 || hours == 12 || hours == 24) {
		return 0, errors.New("hours outside of permitted range 1, 2, 3, 4, 6, 8, 12, 24")
	}

	r, err := s.commandInt("W90", hours, 1, 24, "cycles")
	if err != nil {
		return 0, err
	}

	return int(r), nil
}

func (s *Spanet) SetAutoSanitiseTime(when time.Time) (time.Time, error) {
	return s.commandTime("W73", when)
}

func (s *Spanet) SetTimeout(minutes int) (int, error) {
	return s.commandInt("W74", minutes, 10, 60, "minutes")
}

func (s *Spanet) SetHeatPumpMode(mode HeatPumpMode) (HeatPumpMode, error) {
	newMode, err := s.setMode("W99", byte(mode))
	return HeatPumpMode(newMode), err
}

func (s *Spanet) SetSVElementBoost(enabled bool) (bool, error) {
	val := "0"
	if enabled {
		val = "1"
	}

	r, err := s.commandExpect(fmt.Sprintf("W98:%s", val), val)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(strings.TrimSpace(r))
}

func (s *Spanet) SetLockMode(mode LockMode) (LockMode, error) {
	newMode, err := s.setMode("S21", byte(mode))
	return LockMode(newMode), err
}

func (s *Spanet) GetStatus() (Status, error) {
	r, err := s.command("RF")
	if err != nil {
		return Status{}, err
	}

	return ParseStatus(r)
}
