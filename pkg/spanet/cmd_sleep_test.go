package spanet_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/stretchr/testify/assert"
)

func TestSetSleepTimer(t *testing.T) {
	spa, test := net.Pipe()
	go func() {
		b := make([]byte, 16)
		for {
			sz, err := spa.Read(b)
			if !assert.NoError(t, err) {
				assert.FailNow(t, "We just can't continue like this")
			}

			if sz == 1 && b[0] == '\n' {
				continue
			}

			assert.Equal(t, []byte{'W', '6', '7', ':'}, b[0:4])
			n, err := strconv.ParseInt(string(b[4:sz-1]), 10, 64)
			assert.NoError(t, err)
			assert.True(t, spanet.SleepTimerState(n).IsASleepTimerState())

			spa.Write(b[4:sz])
		}
	}()
	client := spanet.New(test)
	newState, err := client.SetSleepTimerState(1, spanet.SleepTimerStateOff)
	assert.NoError(t, err)
	assert.Equal(t, spanet.SleepTimerStateOff, newState)
}
