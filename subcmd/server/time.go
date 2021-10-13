package server

import (
	"encoding/json"
	"time"
)

type timeParam struct {
	time.Time
}

type timeDateParam struct {
	time.Time
}

func (t *timeParam) UnmarshalJSON(b []byte) (err error) {
	var tmp string
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	t.Time, err = time.Parse("15:04", tmp)

	return err
}

func (t *timeDateParam) UnmarshalJSON(b []byte) (err error) {
	var tmp string
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	t.Time, err = time.Parse("2006-01-02 15:04", tmp)

	return err
}
