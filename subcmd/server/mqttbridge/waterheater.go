package mqttbridge

import (
	"strconv"

	"github.com/freman/spanet/pkg/hamqtt"
	"github.com/freman/spanet/pkg/spanet"
	"github.com/freman/spanet/subcmd/server/middleware/safespa"
)

const waterHeaterObjectID = "water_heater"

// registerWaterHeater wires up a proper Home Assistant water_heater entity
// for the spa's temperature control (current + target temperature).
//
// water_heater doesn't fit the single state/command topic shape every
// other entity in this package uses, so it's built directly on
// hamqtt.Client's lower-level primitives rather than through entities().
// Its state is published from Bridge.poll via publishWaterHeaterState,
// alongside the rest.
//
// water_heater's "modes" are a fixed Home Assistant vocabulary
// (eco/electric/gas/heat_pump/...) that doesn't map cleanly onto the spa's
// own OperationMode (NORM/ECON/AWAY/WEEK), which is already exposed as its
// own select entity - so mode support is intentionally left out here
// rather than force a misleading mapping onto it.
func registerWaterHeater(mqc *hamqtt.Client, spa *safespa.SafeSpa) error {
	commandTopic := mqc.Topic(waterHeaterObjectID, "target_temperature/set")

	payload := map[string]any{
		"name":                      "Spa",
		"current_temperature_topic": mqc.Topic(waterHeaterObjectID, "current_temperature"),
		"temperature_state_topic":   mqc.Topic(waterHeaterObjectID, "target_temperature"),
		"temperature_command_topic": commandTopic,
		"min_temp":                  5,
		"max_temp":                  41,
		"precision":                 0.5,
	}

	if err := mqc.PublishDiscovery("water_heater", waterHeaterObjectID, payload); err != nil {
		return err
	}

	return mqc.SubscribeRaw(commandTopic, func(payload string) error {
		v, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			return err
		}

		return spa.Do(func(s *spanet.Spanet) error {
			_, err := s.SetTargetTemperature(v)

			return err
		})
	})
}

func publishWaterHeaterState(mqc *hamqtt.Client, status spanet.Status) {
	mqc.PublishRaw(mqc.Topic(waterHeaterObjectID, "current_temperature"), fmtFloat1(status.WaterTemperature))
	mqc.PublishRaw(mqc.Topic(waterHeaterObjectID, "target_temperature"), fmtFloat1(status.SetTemperature))
}
