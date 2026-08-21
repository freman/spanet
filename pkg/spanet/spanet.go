package spanet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// Spanet is not safe for concurrent use; callers are expected to serialize
// access (see subcmd/server/middleware/safespa).
type Spanet struct {
	c    net.Conn
	dial func() (net.Conn, error)
}

// New wraps an already-established connection. If it drops, commands will
// fail rather than transparently reconnect; use NewWithDialer for that.
func New(c net.Conn) *Spanet {
	s := &Spanet{c: c}
	s.prime()

	return s
}

// NewWithDialer establishes a connection via dial, and transparently
// redials (once per command) if the connection is later found to be dead.
func NewWithDialer(dial func() (net.Conn, error)) (*Spanet, error) {
	c, err := dial()
	if err != nil {
		return nil, err
	}

	s := &Spanet{c: c, dial: dial}
	s.prime()

	return s, nil
}

// prime improves reliability by always starting on a new line. Best effort:
// if this write fails the first real command will surface the same
// underlying error.
func (s *Spanet) prime() {
	time.Sleep(100 * time.Millisecond)
	_, _ = s.c.Write([]byte{'\n'})
	time.Sleep(100 * time.Millisecond)
}

// reconnect replaces a dead connection by redialing, if this Spanet was
// constructed with a dialer. It reports whether a new connection is ready.
func (s *Spanet) reconnect() bool {
	if s.dial == nil {
		return false
	}

	_ = s.c.Close()

	c, err := s.dial()
	if err != nil {
		return false
	}

	s.c = c
	s.prime()

	return true
}

func (s *Spanet) command(command string) (io.Reader, error) {
	if _, err := s.c.Write(append([]byte(command), '\n')); err != nil {
		return nil, err
	}

	return s.c, nil
}

func (s *Spanet) commandExpect(command, expect string) (string, error) {
	rs, err := s.tryCommandExpect(command, expect)
	if err == nil {
		return rs, nil
	}

	// A response that doesn't match what we expected is a protocol-level
	// problem, not a dead connection - retrying won't help.
	if _, isUnexpected := errors.AsType[ErrUnexpectedResponse](err); isUnexpected {
		return "", err
	}

	if !s.reconnect() {
		return "", err
	}

	return s.tryCommandExpect(command, expect)
}

func (s *Spanet) tryCommandExpect(command, expect string) (string, error) {
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
