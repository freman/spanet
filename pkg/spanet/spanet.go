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
	r, err := s.commandInt(cmd, when.Hour()*256+when.Minute(), 0, 6204, "time", "%d")
	if err != nil {
		return time.Time{}, err
	}

	return spa256toTime(int(r)), nil
}

func (s *Spanet) commandInt(cmd string, value, min, max int, name string, format ...string) (int, error) {
	if value < min || value > max {
		return 0, ErrValueOutOfRange{min, max, value, name}
	}

	format = append(format, "%d")

	// Amusingly while the unit expects 0 prefixed months, it'll return the month without it
	r, err := s.commandExpect(fmt.Sprintf("%s:"+format[0], cmd, value), strings.TrimLeft(fmt.Sprintf(format[0], value), "0"))
	if err != nil {
		return 0, err
	}

	// Weirdly, sometimes S01 returns $year\nS01\n so we'll just grab the first returned string
	firstChunk := strings.TrimLeft(strings.TrimSpace(strings.Split(r, "\n")[0]), "0")

	i, err := strconv.ParseInt(firstChunk, 10, 64)
	return int(i), err
}

func (s *Spanet) setMode(command string, mode byte) (byte, error) {
	arg := strconv.Itoa(int(mode))

	r, err := s.commandExpect(fmt.Sprintf("%s:%s", command, arg), arg)
	if err != nil {
		return 0, err
	}

	// Weirdly, sometimes S21 returns $mode\nS21\n so we'll just grab the first returned string
	firstChunk := strings.TrimSpace(strings.Split(r, "\n")[0])
	newMode, err := strconv.ParseInt(firstChunk, 10, 64)
	return byte(newMode), err
}

func (s *Spanet) Close() error {
	return s.c.Close()
}
