package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/freman/spanet/subcmd/server/middleware/safespa"
)

func TestServerSetPowerSave(t *testing.T) {
	ch := make(chan spanet.PowerSaveMode, 1)

	client, done := MockSpa4Command(t, 16, []byte{'W', '6', '3'}, func(b []byte) []byte {
		n, err := strconv.ParseInt(string(b), 10, 64)
		assert.NoError(t, err)
		assert.Equal(t, <-ch, spanet.PowerSaveMode(n))

		return b
	})
	defer done()

	e := echo.New()
	defineRoutes(e, safespa.New(safespa.WithSpanet(client)))

	server := httptest.NewServer(e)
	defer server.Close()

	for _, v := range spanet.PowerSaveModeValues() {
		t.Run(fmt.Sprintf("Send %s", v.String()), func(t *testing.T) {
			ch <- v

			doRequest(t, server, http.MethodPost, "spa/powersave/mode", struct{ Mode spanet.PowerSaveMode }{v})
		})
	}
}

func doRequest(t *testing.T, server *httptest.Server, method, uriPath string, body interface{}) {
	client := server.Client()

	bodyBuf, err := json.Marshal(body)
	require.NoError(t, err, "Failed to marshal the body for the test")

	bodyReader := bytes.NewReader(bodyBuf)

	req := httptest.NewRequest(method, server.URL+"/"+uriPath, bodyReader)
	req.RequestURI = ""

	if body != nil {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	_, err = client.Do(req)
	require.NoError(t, err, "Failed to even call the service function")
}
