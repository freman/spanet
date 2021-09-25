package spanet

import (
	"testing"
	"time"
)

func TestSpa256ToTime(t *testing.T) {
	tests := []struct {
		i int
		t time.Time
	}{
		{5120, time.Date(0, 0, 0, 20, 0, 0, 0, time.Local)},
		{3375, time.Date(0, 0, 0, 13, 47, 0, 0, time.Local)},
		{1300, time.Date(0, 0, 0, 5, 20, 0, 0, time.Local)},
	}

	for _, v := range tests {
		v := v
		t.Run(v.t.String(), func(t *testing.T) {
			t.Parallel()

			if !v.t.Equal(spa256toTime(v.i)) {
				t.Error("math failure")
			}
		})
	}
}
