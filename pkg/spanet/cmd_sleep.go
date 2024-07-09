package spanet

import (
	"fmt"
	"strconv"
	"time"
)

func (s *Spanet) setSleepTimer(timer int, offset int, value string) (int, error) {
	if timer < 1 || timer > 2 {
		return 0, ErrValueOutOfRange{1, 2, timer, "timer"}
	}

	// All 3 commands for setting the 2 timers (6 commands in total) are in the sequence
	// SetState, SetStart, SetEnd hence the funky math with the offset by 3
	cmd := fmt.Sprintf("W%d:%s", 67+offset+(timer-1)*3, value)

	if _, err := s.commandExpect(cmd, value); err != nil {
		return 0, err
	}

	tmp, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(tmp), nil
}

func (s *Spanet) setSleepTimerTime(timer, offset int, when time.Time) (time.Time, error) {
	value := fmt.Sprintf("%04d", when.Hour()*256+when.Minute())

	tmp, err := s.setSleepTimer(timer, offset, value)
	if err != nil {
		return time.Time{}, err
	}

	return spa256toTime(tmp), nil
}

func (s *Spanet) SetSleepTimerState(timer int, state SleepTimerState) (SleepTimerState, error) {
	tmp, err := s.setSleepTimer(timer, 0, fmt.Sprintf("%d", state))
	if err != nil {
		return 0, err
	}
	return SleepTimerState(tmp), nil
}

func (s *Spanet) SetSleepTimerStart(timer int, when time.Time) (time.Time, error) {
	return s.setSleepTimerTime(timer, 1, when)
}

func (s *Spanet) SetSleepTimerEnd(timer int, when time.Time) (time.Time, error) {
	return s.setSleepTimerTime(timer, 2, when)
}
