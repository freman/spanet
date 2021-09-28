package safespa

import (
	"net"
	"net/http"
	"sync"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/labstack/echo/v4"
)

type SafeSpa struct {
	addr string
	mu   sync.Mutex
	*spanet.Spanet
}

func New(addr string) *SafeSpa {
	return &SafeSpa{
		addr: addr,
	}
}

func (s *SafeSpa) Mutex(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s.mu.Lock()
		defer s.mu.Unlock()

		if s.Spanet == nil {
			c, err := net.Dial("tcp", s.addr)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadGateway, err)
			}

			s.Spanet = spanet.New(c)
		}

		return next(c)
	}
}
