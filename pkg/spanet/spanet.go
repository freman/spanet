package spanet

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type Spanet struct {
	c net.Conn
}

func New(c net.Conn) *Spanet {
	// Improve reliability by always starting on a new line
	time.Sleep(100 * time.Millisecond)
	c.Write([]byte{'\n'})
	time.Sleep(100 * time.Millisecond)
	return &Spanet{c}
}

func (s *Spanet) command(command string) (io.Reader, error) {
	if _, err := s.c.Write(append([]byte(command), '\n')); err != nil {
		return nil, err
	}
	return s.c, nil
}

func (s *Spanet) commandExpect(command, expect string) (string, error) {
	r, err := s.command(command)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 1024)
	c, err := r.Read(buf)
	if err != nil {
		return "", err
	}

	rs := string(buf[:c])
	if !strings.Contains(rs, expect) {
		return "", ErrUnexpectedResponse{expect, rs}
	}

	return rs, nil
}

func (s *Spanet) commandTime(cmd string, when time.Time) (time.Time, error) {
	r, err := s.command(fmt.Sprintf("%s:%04d", cmd, when.Hour()*256+when.Minute()))
	if err != nil {
		return time.Time{}, err
	}
	_ = r
	rs := ""
	tmp, err := strconv.ParseInt(rs, 10, 64)
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

	_, err := s.commandExpect(fmt.Sprintf("%s:"+format[0], cmd, value), fmt.Sprintf(format[0], value))

	return value, err
}

func (s *Spanet) Close() error {
	return s.c.Close()
}
