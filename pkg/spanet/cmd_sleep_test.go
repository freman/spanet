package spanet_test

import (
	"strconv"
	"testing"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/stretchr/testify/assert"
)

func TestSetSleepTimer(t *testing.T) {
	client, done := MockSpa4Command(t, 16, []byte{'W', '6', '7'}, func(b []byte) []byte {
		n, err := strconv.ParseInt(string(b), 10, 64)
		assert.NoError(t, err)
		assert.True(t, spanet.SleepTimerState(n).IsASleepTimerState())

		return b
	})
	defer done()

	newState, err := client.SetSleepTimerState(1, spanet.SleepTimerStateOff)
	assert.NoError(t, err)
	assert.Equal(t, spanet.SleepTimerStateOff, newState)
}
