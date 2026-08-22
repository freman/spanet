package spanet

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStatusServer answers any "RF" request with a minimal but valid RG
// terminator line, until the connection closes.
func mockStatusServer(spa net.Conn) {
	buf := make([]byte, 16)

	for {
		n, err := spa.Read(buf)
		if err != nil {
			return
		}

		if n == 1 && buf[0] == '\n' {
			continue
		}

		_, _ = spa.Write([]byte(",RG,0,0,0,0,0,0,0-0-0,0-0-0,0-0-0,0-0-0,0-0-0,0\n"))
	}
}

// TestGetStatusTimesOutThenReconnects covers the failure mode that shipped
// undetected: a peer the OS considers connected but that never answers a
// command used to hang net.Conn.Read forever, with no error and nothing
// for a caller (or a reconnect loop) to react to. It should now time out
// and successfully redial instead.
func TestGetStatusTimesOutThenReconnects(t *testing.T) {
	orig := ioTimeout
	ioTimeout = 150 * time.Millisecond

	defer func() { ioTimeout = orig }()

	dialCount := 0

	dial := func() (net.Conn, error) {
		dialCount++

		spaSide, clientSide := net.Pipe()

		if dialCount == 1 {
			// Accepts and silently discards everything - never answers,
			// simulating a wedged WiFly bridge that still lets a new TCP
			// connection through.
			go func() { _, _ = io.Copy(io.Discard, spaSide) }()
		} else {
			go mockStatusServer(spaSide)
		}

		return clientSide, nil
	}

	spa, err := NewWithDialer(dial)
	require.NoError(t, err)

	start := time.Now()
	_, err = spa.GetStatus()
	elapsed := time.Since(start)

	require.NoError(t, err, "GetStatus should succeed after reconnecting to a responsive peer")
	assert.Less(t, elapsed, 2*time.Second, "should time out and recover, not hang indefinitely")
	assert.Equal(t, 2, dialCount, "should have redialed exactly once after the first attempt timed out")
}
