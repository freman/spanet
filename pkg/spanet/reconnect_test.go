package spanet_test

import (
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/freman/spanet/pkg/spanet"
)

// mockSpaServer answers a single command type on spa until the connection
// is closed (by the test, or by the client tearing down).
func mockSpaServer(t *testing.T, spa net.Conn, cmd []byte, checkFn func(b []byte) []byte) {
	t.Helper()

	b := make([]byte, 16)
	for {
		sz, err := spa.Read(b)
		if err != nil {
			return
		}

		if sz == 1 && b[0] == '\n' {
			continue
		}

		assert.Equal(t, cmd, b[0:3])
		if res := checkFn(b[4 : sz-1]); res != nil {
			_, _ = spa.Write(append(res, '\n'))
		}
	}
}

// TestReconnectAfterConnectionDrop verifies that a Spanet built with
// NewWithDialer transparently redials and retries a command once when the
// underlying connection has died, rather than surfacing the error.
func TestReconnectAfterConnectionDrop(t *testing.T) {
	var (
		mu             sync.Mutex
		dialCount      int
		currentSpaSide net.Conn
	)

	cmd := []byte{'W', '6', '3'}

	dial := func() (net.Conn, error) {
		spaSide, clientSide := net.Pipe()

		mu.Lock()
		dialCount++
		currentSpaSide = spaSide
		mu.Unlock()

		go mockSpaServer(t, spaSide, cmd, func(b []byte) []byte { return b })

		return clientSide, nil
	}

	client, err := spanet.NewWithDialer(dial)
	require.NoError(t, err)

	defer func() {
		mu.Lock()
		_ = currentSpaSide.Close()
		mu.Unlock()
	}()

	_, err = client.SetPowerSave(spanet.PowerSaveModeHigh)
	require.NoError(t, err)

	mu.Lock()
	assert.Equal(t, 1, dialCount)
	// Sever the connection out from under the client, simulating the spa
	// dropping it, and confirm the next command redials rather than failing.
	require.NoError(t, currentSpaSide.Close())
	mu.Unlock()

	_, err = client.SetPowerSave(spanet.PowerSaveModeHigh)
	require.NoError(t, err)

	mu.Lock()
	assert.Equal(t, 2, dialCount)
	mu.Unlock()
}
