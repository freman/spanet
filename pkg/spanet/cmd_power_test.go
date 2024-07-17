package spanet_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/freman/spanet/pkg/spanet"
)

func TestSetPowerSave(t *testing.T) {
	client, done := MockSpa4Command(t, 16, []byte{'W', '6', '3'}, func(b []byte) []byte {
		n, err := strconv.ParseInt(string(b), 10, 64)
		assert.NoError(t, err)
		assert.True(t, spanet.PowerSaveMode(n).IsAPowerSaveMode())

		return b
	})
	defer done()

	newState, err := client.SetPowerSave(spanet.PowerSaveModeHigh)
	assert.NoError(t, err)
	assert.Equal(t, spanet.PowerSaveModeHigh, newState)
}
