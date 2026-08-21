package mqttbridge

import (
	"fmt"
	"strconv"
	"time"

	"github.com/freman/spanet/pkg/hamqtt"
	"github.com/freman/spanet/pkg/spanet"
)

// entity pairs a Home Assistant discovery entity with the logic needed to
// render its state from a status snapshot and handle incoming commands.
// State is nil for write-only entities (button). raw is nil for read-only
// entities (sensor, binary_sensor); entities() rebinds it into Entity.Handler
// so every command runs through safeSpa.Do.
type entity struct {
	hamqtt.Entity
	State func(status spanet.Status) (value string, ok bool)
	raw   func(spa *spanet.Spanet, payload string) error
}

func onOff(b bool) string {
	if b {
		return "ON"
	}

	return "OFF"
}

func parseOnOff(payload string) (bool, error) {
	switch payload {
	case "ON":
		return true, nil
	case "OFF":
		return false, nil
	default:
		return false, fmt.Errorf("expected ON or OFF, got %q", payload)
	}
}

// timeLayout matches the MQTT time platform's contract: it always sends
// command payloads in HH:MM:SS, and we publish state in the same format.
const timeLayout = "15:04:05"

func parseClock(payload string) (time.Time, error) {
	return time.Parse(timeLayout, payload)
}

// setterRaw builds an entity.raw command handler that parses the incoming
// MQTT payload with parse before calling set.
func setterRaw[T any](parse func(string) (T, error), set func(*spanet.Spanet, T) error) func(*spanet.Spanet, string) error {
	return func(spa *spanet.Spanet, payload string) error {
		v, err := parse(payload)
		if err != nil {
			return err
		}

		return set(spa, v)
	}
}

func numberConfig(minimum, maximum, step float64, unit string) map[string]any {
	cfg := map[string]any{
		"min":  minimum,
		"max":  maximum,
		"step": step,
	}
	if unit != "" {
		cfg["unit_of_measurement"] = unit
	}

	return cfg
}

func diagnostic(cfg map[string]any) map[string]any {
	cfg["entity_category"] = "diagnostic"

	return cfg
}

func measurement(deviceClass, unit string) map[string]any {
	cfg := map[string]any{"state_class": "measurement"}
	if deviceClass != "" {
		cfg["device_class"] = deviceClass
	}

	if unit != "" {
		cfg["unit_of_measurement"] = unit
	}

	return cfg
}

func fmtFloat1(f float64) string {
	return strconv.FormatFloat(f, 'f', 1, 64)
}
