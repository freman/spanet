package server

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/freman/spanet/subcmd/server/middleware/safespa"
)

func TestSetPowerSave(t *testing.T) {
	ch := make(chan spanet.PowerSaveMode, 1)

	client, done := MockSpa4Command(t, 16, []byte{'W', '6', '3'}, func(b []byte) []byte {
		n, err := strconv.ParseInt(string(b), 10, 64)
		assert.NoError(t, err)
		assert.Equal(t, spanet.PowerSaveMode(n), <-ch)

		return b
	})
	defer done()

	svc := service{safespa.New(safespa.WithSpanet(client))}
	setPowerSave := svc.handleSimplePost("SetPowerSave", "Mode")

	for _, v := range spanet.PowerSaveModeValues() {
		t.Run(fmt.Sprintf("Send %s", v.String()), func(t *testing.T) {
			ch <- v

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf(`{"Mode":"%s"}`, v.String())))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := echo.New().NewContext(req, httptest.NewRecorder())

			if !assert.NoError(t, setPowerSave(c)) {
				assert.Fail(t, "Failed to even call the service function")
			}
		})
	}
}

func MockSpa4Command(t *testing.T, bufSz int, cmd []byte, checkFn func(b []byte) []byte) (*spanet.Spanet, func()) {
	spa, test := net.Pipe()
	go func() {
		b := make([]byte, bufSz)
		for {
			sz, err := spa.Read(b)
			if !assert.NoError(t, err) {
				assert.FailNow(t, "We just can't continue like this")
			}

			if sz == 1 && b[0] == '\n' {
				continue
			}

			assert.Equal(t, cmd, b[0:3])
			if res := checkFn(b[4 : sz-1]); res != nil {
				spa.Write(append(res, '\n'))
			}
		}
	}()

	return spanet.New(test), func() {
		if !assert.NoError(t, spa.Close()) {
			assert.FailNow(t, "failed to close spa conn")
		}

		if !assert.NoError(t, test.Close()) {
			assert.FailNow(t, "failed to close test conn")
		}
	}
}
