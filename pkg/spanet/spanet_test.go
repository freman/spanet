package spanet_test

import (
	"net"
	"testing"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/stretchr/testify/assert"
)

func MockSpa4Command(t *testing.T, bufSz int, cmd []byte, checkFn func(b []byte) []byte) (*spanet.Spanet, func()) {
	spa, test := net.Pipe()
	go func() {
		b := make([]byte, bufSz)
		for {
			sz, err := spa.Read(b)
			if err != nil {
				// The pipe was closed by the test's cleanup func; nothing left to do.
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
