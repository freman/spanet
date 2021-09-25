package spanet

import (
	"errors"
	"fmt"
	"strconv"
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

func (s *Spanet) SetOperationalMode(mode OperationMode) (OperationMode, error) {
	var arg string
	switch mode {
	case OperationModeNormal:
		arg = "0"
	case OperationModeEconomy:
		arg = "1"
	case OperationModeAway:
		arg = "2"
	case OperationModeWeekdays:
		arg = "3"
	default:
		return "", errors.New("Unsupported operational mode")
	}

	r, err := s.commandExpect(fmt.Sprintf("W66:%s", arg), arg)
	if err != nil {
		return "", err
	}

	switch r {
	case "0":
		return OperationModeNormal, nil
	case "1":
		return OperationModeEconomy, nil
	case "2":
		return OperationModeAway, nil
	case "3":
		return OperationModeWeekdays, nil
	default:
		return "", errors.New("unexpected response")
	}
}

func (s *Spanet) SetFiltrationRunTime(hours int) (int, error) {
	return s.commandInt("W60", hours, 1, 24, "hours")
}

func (s *Spanet) SetFiltrationCycle(hours int) (int, error) {
	if !(hours >= 1 || hours <= 4 || hours == 6 || hours == 8 || hours == 12 || hours == 24) {
		return 0, errors.New("hours outside of permitted range 1, 2, 3, 4, 6, 8, 12, 24")
	}

	r, err := s.command(fmt.Sprintf("W90:%d", hours))
	if err != nil {
		return 0, err
	}

	tmp, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(tmp), nil
}

func (s *Spanet) SetAutoSanitiseTime(when time.Time) (time.Time, error) {
	return s.commandTime("W73", when)
}

func (s *Spanet) SetTimeout(minutes int) (int, error) {
	return s.commandInt("W74", minutes, 10, 60, "minutes")
}

func (s *Spanet) SetHeatPumpMode(mode HeatPumpMode) (HeatPumpMode, error) {
	r, err := s.command(fmt.Sprintf("W99:%d", mode))
	if err != nil {
		return 0, err
	}

	tmp, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return 0, err
	}

	return HeatPumpMode(tmp), nil
}

func (s *Spanet) SetSVElementBoost(enabled bool) (bool, error) {
	val := 0
	if enabled {
		val = 1
	}

	r, err := s.command(fmt.Sprintf("W98:%d", val))
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(r)
}

func (s *Spanet) SetLockMode(mode LockMode) (LockMode, error) {
	r, err := s.command(fmt.Sprintf("S21:%d", mode))
	if err != nil {
		return 0, err
	}

	tmp, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return 0, err
	}

	return LockMode(tmp), nil
}
