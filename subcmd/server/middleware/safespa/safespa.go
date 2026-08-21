package safespa

import (
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/labstack/echo/v5"

	"github.com/freman/spanet/pkg/spanet"
)

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

func (s *SafeSpa) Mutex(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		s.mu.Lock()
		defer s.mu.Unlock()

		if s.Spanet == nil {
			conn, err := net.Dial("tcp", s.addr)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadGateway, err.Error()).Wrap(err)
			}

			s.Spanet = spanet.New(conn)
		}

		err := next(c)

		// Hack: until I can reliably detect the spa dropping the connection
		// every connection will be a new connection, we'll just pretend it's recycled.
		// Don't destroy the connection if there's no address to re-create it - handy for tests
		if s.addr != "" {
			if closeErr := s.Close(); closeErr != nil {
				slog.Warn("failed to close spa connection", "error", closeErr)
			}

			s.Spanet = nil
		}

		return err
	}
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
