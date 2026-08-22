package safespa

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/freman/spanet/pkg/spanet"
)

// dialTimeout bounds establishing the TCP connection itself, separately
// from pkg/spanet's own read/write deadlines on an already-open one.
const dialTimeout = 10 * time.Second

// SafeSpa holds a single, persistent connection to the spa, guarded by a
// mutex so concurrent HTTP requests serialize onto it rather than each
// dialing their own (the spa's wifi bridge only tolerates one connection
// at a time). A dead connection is transparently redialed by pkg/spanet
// on next use - see spanet.NewWithDialer.
type SafeSpa struct {
	addr string
	mu   sync.Mutex
	*spanet.Spanet
}

func New(opt initOpt) *SafeSpa {
	var s SafeSpa
	opt(&s)

	return &s
}

// Do locks the connection, dialing it first if this is the first use (or
// the previous connection was torn down), and runs fn against it. Any
// caller that needs to talk to the spa - HTTP handlers via Mutex, or an
// MQTT bridge - should go through Do so they all serialize onto the same
// connection.
func (s *SafeSpa) Do(fn func(*spanet.Spanet) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Spanet == nil {
		spa, err := spanet.NewWithDialer(func() (net.Conn, error) {
			return net.DialTimeout("tcp", s.addr, dialTimeout)
		})
		if err != nil {
			return err
		}

		s.Spanet = spa
	}

	return fn(s.Spanet)
}

func (s *SafeSpa) Mutex(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		err := s.Do(func(*spanet.Spanet) error {
			return next(c)
		})
		if err != nil {
			if httpErr, isa := errors.AsType[*echo.HTTPError](err); isa {
				return httpErr
			}

			return echo.NewHTTPError(http.StatusBadGateway, err.Error()).Wrap(err)
		}

		return nil
	}
}

// Close closes the underlying connection, if one was ever established.
// It shadows the promoted *spanet.Spanet.Close, which would otherwise
// panic on a nil embedded Spanet.
func (s *SafeSpa) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Spanet == nil {
		return nil
	}

	return s.Spanet.Close()
}

type initOpt func(s *SafeSpa)

func WithAddr(addr string) func(s *SafeSpa) {
	return func(s *SafeSpa) {
		s.addr = addr
	}
}

func WithSpanet(spa *spanet.Spanet) func(s *SafeSpa) {
	return func(s *SafeSpa) {
		s.Spanet = spa
	}
}
