package spanet

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Spanet struct {
}

func (s *Spanet) command(command string) (string, error) {
	return "", errors.New("not implemented")
}

func (s *Spanet) commandExpect(command, expect string) (string, error) {
	r, err := s.command(command)
	if err != nil {
		return "", err
	}

	if r != expect {
		return "", ErrUnexpectedResponse{expect, r}
	}

	return r, errors.New("not implemented")
}

func (s *Spanet) commandTime(cmd string, when time.Time) (time.Time, error) {
	r, err := s.command(fmt.Sprintf("%s:%04d", cmd, when.Hour()*256+when.Minute()))
	if err != nil {
		return time.Time{}, err
	}

	tmp, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return spa256toTime(int(tmp)), nil
}

func (s *Spanet) commandInt(cmd string, value, min, max int, name string, format ...string) (int, error) {
	if value < min || value > max {
		return 0, ErrValueOutOfRange{min, max, value, name}
	}

	format = append(format, "%d")

	r, err := s.command(fmt.Sprintf("%s:"+format[0], cmd, value))
	if err != nil {
		return 0, err
	}

	tmp, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(tmp), nil
}
