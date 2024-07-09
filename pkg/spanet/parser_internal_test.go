package spanet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpa256ToTime(t *testing.T) {
	tests := []struct {
		i int
		t time.Time
	}{
		{5120, time.Date(1, 0, 0, 20, 0, 0, 0, time.UTC)},
		{3375, time.Date(1, 0, 0, 13, 47, 0, 0, time.UTC)},
		{1300, time.Date(1, 0, 0, 5, 20, 0, 0, time.UTC)},
	}

	for _, v := range tests {
		v := v
		t.Run(v.t.String(), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, v.t, spa256toTime(v.i))
		})
	}
}
