package spanet_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/stretchr/testify/assert"
)

func TestSetAutoSanitiseTime(t *testing.T) {
	client, done := MockSpa4Command(t, 16, []byte{'W', '7', '3'}, func(b []byte) []byte {
		n, err := strconv.ParseInt(string(b), 10, 64)
		assert.NoError(t, err)
		assert.Equal(t, 2836, int(n))

		return b
	})
	defer done()

	expect := time.Date(1, 0, 0, 11, 20, 0, 0, time.UTC)
	newState, err := client.SetAutoSanitiseTime(expect)
	assert.NoError(t, err)
	assert.Equal(t, expect, newState)
}

func TestSetLockMode(t *testing.T) {
	client, done := MockSpa4Command(t, 16, []byte{'S', '2', '1'}, func(b []byte) []byte {
		n, err := strconv.ParseInt(string(b), 10, 64)
		assert.NoError(t, err)
		assert.Equal(t, 0, int(n))

		return b
	})
	defer done()

	newState, err := client.SetLockMode(spanet.LockModeOff)
	assert.NoError(t, err)
	assert.Equal(t, spanet.LockModeOff, newState)
}

func TestBuggedSetLockMode(t *testing.T) {
	client, done := MockSpa4Command(t, 16, []byte{'S', '2', '1'}, func(b []byte) []byte {
		n, err := strconv.ParseInt(string(b), 10, 64)
		assert.NoError(t, err)
		assert.Equal(t, 0, int(n))

		return append(b, '\n', 'S', '2', '1')
	})
	defer done()

	newState, err := client.SetLockMode(spanet.LockModeOff)
	assert.NoError(t, err)
	assert.Equal(t, spanet.LockModeOff, newState)
}
