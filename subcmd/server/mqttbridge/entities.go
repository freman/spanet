package mqttbridge

import (
	"fmt"
	"strconv"
	"time"

	"github.com/freman/spanet/pkg/hamqtt"
	"github.com/freman/spanet/pkg/spanet"
	"github.com/freman/spanet/subcmd/server/middleware/safespa"
)

// entities returns every Home Assistant entity this bridge exposes for spa.
// installed is a snapshot of which pumps are fitted and what states each
// supports, taken once at startup - the spa doesn't report that changing,
// and Home Assistant discovery entities aren't meant to be redefined on
// every poll.
func entities(spa *safespa.SafeSpa, installed []spanet.Pump) []entity {
	var es []entity

	es = append(es, temperatureEntities()...)
	es = append(es, powerEntities()...)
	es = append(es, diagnosticEntities()...)
	es = append(es, statusEntities()...)
	es = append(es, lightEntities()...)
	es = append(es, blowerEntities()...)
	es = append(es, scheduleEntities()...)
	es = append(es, modeEntities()...)
	es = append(es, pumpEntities(installed)...)
	es = append(es, sleepTimerEntities()...)

	for i := range es {
		wireHandler(spa, &es[i])
	}

	return es
}

// wireHandler rebinds an entity's Handler (currently a func(*spanet.Spanet)
// error built by the *Entities funcs below, stashed via handlerFn) to run
// through safeSpa.Do, so every command handler serializes onto the shared
// connection the same way HTTP requests do.
func wireHandler(spa *safespa.SafeSpa, e *entity) {
	if e.raw == nil {
		return
	}

	raw := e.raw
	e.Handler = func(payload string) error {
		return spa.Do(func(s *spanet.Spanet) error {
			return raw(s, payload)
		})
	}
}

// temperatureEntities covers the temperature sensors that aren't part of
// the water_heater entity (see waterheater.go), which already exposes
// SetTemperature and WaterTemperature as its target/current temperature.
func temperatureEntities() []entity {
	return []entity{
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "heater_temperature",
				Name:      "Heater Temperature",
				Config:    measurement("temperature", "°C"),
			},
			State: func(s spanet.Status) (string, bool) { return fmtFloat1(s.HeaterTemperature), true },
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "case_temperature",
				Name:      "Case Temperature",
				Config:    diagnostic(measurement("temperature", "°C")),
			},
			State: func(s spanet.Status) (string, bool) { return fmtFloat1(s.CaseTemperature), true },
		},
	}
}

func powerEntities() []entity {
	return []entity{
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "power_volts",
				Name:      "Voltage",
				Config:    diagnostic(measurement("voltage", "V")),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(s.Power.Volts), true },
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "power_amps",
				Name:      "Current",
				Config:    diagnostic(measurement("current", "A")),
			},
			State: func(s spanet.Status) (string, bool) { return fmtFloat1(s.Power.Amps), true },
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "power_apparent",
				Name:      "Apparent Power",
				// Volts * Amps is apparent power (VA), not true real power (W) -
				// the spa doesn't report a power factor, so we can't compute
				// real watts. Labelled and typed accordingly so it doesn't get
				// mistaken for (or summed with) a real energy sensor.
				Config: diagnostic(measurement("apparent_power", "VA")),
			},
			State: func(s spanet.Status) (string, bool) {
				return fmtFloat1(float64(s.Power.Volts) * s.Power.Amps), true
			},
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "heating_amps",
				Name:      "Heating Current",
				Config:    diagnostic(measurement("current", "A")),
			},
			State: func(s spanet.Status) (string, bool) { return fmtFloat1(s.Power.HeatingAmps), true },
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "current_limit",
				Name:      "Current Limit",
				Config:    diagnostic(measurement("current", "A")),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(s.Power.CurrentLimit), true },
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "load_shed",
				Name:      "Load Shed",
				Config:    diagnostic(map[string]any{}),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(s.Power.LoadShed), true },
		},
	}
}

func diagnosticEntities() []entity {
	return []entity{
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "runtime_hours",
				Name:      "Runtime",
				Config:    diagnostic(measurement("duration", "h")),
			},
			State: func(s spanet.Status) (string, bool) { return fmtFloat1(s.RuntimeHours), true },
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "filter_pump_total_runtime",
				Name:      "Filter Pump Total Runtime",
				Config:    diagnostic(measurement("duration", "min")),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(s.FilterPumpTotalRunTime), true },
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "filter_pump_req",
				Name:      "Filter Pump Request",
				Config:    diagnostic(map[string]any{}),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(s.FilterPumpReq), true },
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "awake_remains",
				Name:      "Awake Time Remaining",
				Config:    diagnostic(measurement("duration", "min")),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(s.AwakeRemains), true },
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "spa_datetime",
				Name:      "Spa Clock",
				Config:    diagnostic(map[string]any{"device_class": "timestamp"}),
			},
			State: func(s spanet.Status) (string, bool) { return s.TimeDate.AsTime().Format(time.RFC3339), true },
		},
		{
			Entity: hamqtt.Entity{
				Component: "button",
				ObjectID:  "sync_datetime",
				Name:      "Sync Spa Clock",
				Config:    diagnostic(map[string]any{}),
			},
			raw: func(spa *spanet.Spanet, _ string) error {
				return spa.SetDateTime(time.Now())
			},
		},
		{
			Entity: hamqtt.Entity{
				Component: "sensor",
				ObjectID:  "last_update",
				Name:      "Last Successful Update",
				Config:    diagnostic(map[string]any{"device_class": "timestamp"}),
			},
			// Called exactly when a poll succeeds (see Bridge.poll), so
			// "now" at call time is the last successful update time.
			State: func(s spanet.Status) (string, bool) { return time.Now().Format(time.RFC3339), true },
		},
	}
}

func statusEntities() []entity {
	return []entity{
		{
			Entity: hamqtt.Entity{Component: "binary_sensor", ObjectID: "water_present", Name: "Water Present"},
			State:  func(s spanet.Status) (string, bool) { return onOff(s.WaterPresent), true },
		},
		{
			Entity: hamqtt.Entity{Component: "binary_sensor", ObjectID: "heating", Name: "Heating", Config: map[string]any{"device_class": "running"}},
			State:  func(s spanet.Status) (string, bool) { return onOff(s.Heating), true },
		},
		{
			Entity: hamqtt.Entity{Component: "binary_sensor", ObjectID: "uv_ozone", Name: "UV/Ozone", Config: map[string]any{"device_class": "running"}},
			State:  func(s spanet.Status) (string, bool) { return onOff(s.UVOzone), true },
		},
		{
			Entity: hamqtt.Entity{Component: "binary_sensor", ObjectID: "auto", Name: "Auto"},
			State:  func(s spanet.Status) (string, bool) { return onOff(s.Auto), true },
		},
		{
			Entity: hamqtt.Entity{Component: "binary_sensor", ObjectID: "sleeping", Name: "Sleeping"},
			State:  func(s spanet.Status) (string, bool) { return onOff(s.Sleeping), true },
		},
		{
			Entity: hamqtt.Entity{Component: "switch", ObjectID: "sanitise", Name: "Sanitise"},
			State:  func(s spanet.Status) (string, bool) { return onOff(s.Sanitise), true },
			raw: func(spa *spanet.Spanet, payload string) error {
				want, err := parseOnOff(payload)
				if err != nil {
					return err
				}

				status, err := spa.GetStatus()
				if err != nil {
					return err
				}

				if status.Sanitise == want {
					return nil
				}

				return spa.ToggleSanitise()
			},
		},
		{
			Entity: hamqtt.Entity{Component: "switch", ObjectID: "sv_element_boost", Name: "SV Element Boost"},
			State:  func(s spanet.Status) (string, bool) { return onOff(s.SVElementBoost), true },
			raw: func(spa *spanet.Spanet, payload string) error {
				want, err := parseOnOff(payload)
				if err != nil {
					return err
				}

				_, err = spa.SetSVElementBoost(want)

				return err
			},
		},
	}
}

func lightEntities() []entity {
	return []entity{
		{
			Entity: hamqtt.Entity{Component: "switch", ObjectID: "lights", Name: "Lights"},
			State:  func(s spanet.Status) (string, bool) { return onOff(s.Lights.On), true },
			raw: func(spa *spanet.Spanet, payload string) error {
				want, err := parseOnOff(payload)
				if err != nil {
					return err
				}

				status, err := spa.GetStatus()
				if err != nil {
					return err
				}

				if status.Lights.On == want {
					return nil
				}

				return spa.ToggleLights()
			},
		},
		{
			Entity: hamqtt.Entity{
				Component: "select",
				ObjectID:  "lights_mode",
				Name:      "Light Mode",
				Config:    map[string]any{"options": spanet.LightsModeStrings()},
			},
			State: func(s spanet.Status) (string, bool) { return s.Lights.Mode.String(), true },
			raw: setterRaw(spanet.LightsModeString, func(spa *spanet.Spanet, v spanet.LightsMode) error {
				_, err := spa.SetLightsMode(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{
				Component: "number",
				ObjectID:  "lights_brightness",
				Name:      "Light Brightness",
				Config:    numberConfig(1, 5, 1, ""),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(int(s.Lights.Brightness)), true },
			raw: setterRaw(strconv.Atoi, func(spa *spanet.Spanet, v int) error {
				_, err := spa.SetLightsBrightness(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{
				Component: "number",
				ObjectID:  "lights_effect_speed",
				Name:      "Light Effect Speed",
				Config:    numberConfig(1, 5, 1, ""),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(int(s.Lights.Speed)), true },
			raw: setterRaw(strconv.Atoi, func(spa *spanet.Spanet, v int) error {
				_, err := spa.SetLightsEffectSpeed(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{
				Component: "number",
				ObjectID:  "lights_colour",
				Name:      "Light Colour",
				Config:    numberConfig(1, 30, 1, ""),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(int(s.Lights.Colour)), true },
			raw: setterRaw(strconv.Atoi, func(spa *spanet.Spanet, v int) error {
				_, err := spa.SetLightsColour(v)

				return err
			}),
		},
	}
}

func blowerEntities() []entity {
	return []entity{
		{
			Entity: hamqtt.Entity{
				Component: "select",
				ObjectID:  "blower_mode",
				Name:      "Blower Mode",
				Config:    map[string]any{"options": spanet.BlowerModeStrings()},
			},
			State: func(s spanet.Status) (string, bool) { return s.Blower.Mode.String(), true },
			raw: setterRaw(spanet.BlowerModeString, func(spa *spanet.Spanet, v spanet.BlowerMode) error {
				return spa.ControlBlower(v)
			}),
		},
		{
			Entity: hamqtt.Entity{
				Component: "number",
				ObjectID:  "blower_speed",
				Name:      "Blower Speed",
				Config:    numberConfig(1, 5, 1, ""),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(int(s.Blower.VariableSpeed)), true },
			raw: setterRaw(strconv.Atoi, func(spa *spanet.Spanet, v int) error {
				_, err := spa.SetBlowerVariableSpeed(v)

				return err
			}),
		},
	}
}

func scheduleEntities() []entity {
	return []entity{
		{
			Entity: hamqtt.Entity{Component: "time", ObjectID: "auto_sanitise_time", Name: "Auto Sanitise Time"},
			State:  func(s spanet.Status) (string, bool) { return s.AutoSanitise.Format(timeLayout), true },
			raw: setterRaw(parseClock, func(spa *spanet.Spanet, v time.Time) error {
				_, err := spa.SetAutoSanitiseTime(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{Component: "time", ObjectID: "peak_start", Name: "Peak Period Start"},
			State:  func(s spanet.Status) (string, bool) { return s.PeakStart.Format(timeLayout), true },
			raw: setterRaw(parseClock, func(spa *spanet.Spanet, v time.Time) error {
				_, err := spa.SetPeakStart(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{Component: "time", ObjectID: "peak_end", Name: "Peak Period End"},
			State:  func(s spanet.Status) (string, bool) { return s.PeakEnd.Format(timeLayout), true },
			raw: setterRaw(parseClock, func(spa *spanet.Spanet, v time.Time) error {
				_, err := spa.SetPeakEnd(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{
				Component: "number",
				ObjectID:  "filtration_runtime",
				Name:      "Filtration Run Time",
				Config:    numberConfig(1, 24, 1, "h"),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(int(s.FiltrationHour)), true },
			raw: setterRaw(strconv.Atoi, func(spa *spanet.Spanet, v int) error {
				_, err := spa.SetFiltrationRunTime(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{
				Component: "select",
				ObjectID:  "filtration_cycle",
				Name:      "Filtration Cycle",
				Config:    map[string]any{"options": []string{"1", "2", "3", "4", "6", "8", "12", "24"}},
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(int(s.FiltrationCycle)), true },
			raw: setterRaw(strconv.Atoi, func(spa *spanet.Spanet, v int) error {
				_, err := spa.SetFiltrationCycle(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{
				Component: "number",
				ObjectID:  "timeout_minutes",
				Name:      "Inactivity Timeout",
				Config:    numberConfig(10, 60, 1, "min"),
			},
			State: func(s spanet.Status) (string, bool) { return strconv.Itoa(int(s.Timeout.Minutes())), true },
			raw: setterRaw(strconv.Atoi, func(spa *spanet.Spanet, v int) error {
				_, err := spa.SetTimeout(v)

				return err
			}),
		},
	}
}

func modeEntities() []entity {
	return []entity{
		{
			Entity: hamqtt.Entity{
				Component: "select",
				ObjectID:  "operation_mode",
				Name:      "Operation Mode",
				Config:    map[string]any{"options": spanet.OperationModeStrings()},
			},
			State: func(s spanet.Status) (string, bool) { return s.OperationMode.String(), true },
			raw: setterRaw(spanet.OperationModeString, func(spa *spanet.Spanet, v spanet.OperationMode) error {
				_, err := spa.SetOperationMode(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{
				Component: "select",
				ObjectID:  "power_save_mode",
				Name:      "Power Save Mode",
				Config:    map[string]any{"options": spanet.PowerSaveModeStrings()},
			},
			State: func(s spanet.Status) (string, bool) { return s.PowerSave.String(), true },
			raw: setterRaw(spanet.PowerSaveModeString, func(spa *spanet.Spanet, v spanet.PowerSaveMode) error {
				_, err := spa.SetPowerSave(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{
				Component: "select",
				ObjectID:  "heat_pump_mode",
				Name:      "Heat Pump Mode",
				Config:    map[string]any{"options": spanet.HeatPumpModeStrings()},
			},
			State: func(s spanet.Status) (string, bool) { return s.HeatPumpMode.String(), true },
			raw: setterRaw(spanet.HeatPumpModeString, func(spa *spanet.Spanet, v spanet.HeatPumpMode) error {
				_, err := spa.SetHeatPumpMode(v)

				return err
			}),
		},
		{
			Entity: hamqtt.Entity{
				Component: "select",
				ObjectID:  "lock_mode",
				Name:      "Lock Mode",
				Config:    map[string]any{"options": spanet.LockModeStrings()},
			},
			State: func(s spanet.Status) (string, bool) { return s.Lock.String(), true },
			raw: setterRaw(spanet.LockModeString, func(spa *spanet.Spanet, v spanet.LockMode) error {
				_, err := spa.SetLockMode(v)

				return err
			}),
		},
	}
}

// pumpEntities returns one select entity per installed pump, using each
// pump's own reported States as the option list (some pumps only support
// Off/On, others also Auto).
func pumpEntities(installed []spanet.Pump) []entity {
	var es []entity

	for i, pump := range installed {
		if !pump.Installed {
			continue
		}

		pumpNum := i + 1 // ControlPump is 1-indexed
		index := i       // capture for closures below

		options := make([]string, len(pump.States))
		for j, st := range pump.States {
			options[j] = st.String()
		}

		es = append(es, entity{
			Entity: hamqtt.Entity{
				Component: "select",
				ObjectID:  fmt.Sprintf("pump_%d", pumpNum),
				Name:      fmt.Sprintf("Pump %d", pumpNum),
				Config:    map[string]any{"options": options},
			},
			State: func(s spanet.Status) (string, bool) {
				if index >= len(s.Pumps) {
					return "", false
				}

				return s.Pumps[index].State.String(), true
			},
			raw: setterRaw(spanet.PumpStateString, func(spa *spanet.Spanet, v spanet.PumpState) error {
				return spa.ControlPump(pumpNum, v)
			}),
		})
	}

	return es
}

func sleepTimerEntities() []entity {
	var es []entity

	for timer := 1; timer <= 2; timer++ {
		timer, index := timer, timer-1

		es = append(es,
			entity{
				Entity: hamqtt.Entity{
					Component: "select",
					ObjectID:  fmt.Sprintf("sleeptimer_%d_state", timer),
					Name:      fmt.Sprintf("Sleep Timer %d", timer),
					Config:    map[string]any{"options": spanet.SleepTimerStateStrings()},
				},
				State: func(s spanet.Status) (string, bool) {
					if index >= len(s.SleepTimers) {
						return "", false
					}

					return s.SleepTimers[index].State.String(), true
				},
				raw: setterRaw(spanet.SleepTimerStateString, func(spa *spanet.Spanet, v spanet.SleepTimerState) error {
					_, err := spa.SetSleepTimerState(timer, v)

					return err
				}),
			},
			entity{
				Entity: hamqtt.Entity{Component: "time", ObjectID: fmt.Sprintf("sleeptimer_%d_start", timer), Name: fmt.Sprintf("Sleep Timer %d Start", timer)},
				State: func(s spanet.Status) (string, bool) {
					if index >= len(s.SleepTimers) {
						return "", false
					}

					return s.SleepTimers[index].StartTime.Format(timeLayout), true
				},
				raw: setterRaw(parseClock, func(spa *spanet.Spanet, v time.Time) error {
					_, err := spa.SetSleepTimerStart(timer, v)

					return err
				}),
			},
			entity{
				Entity: hamqtt.Entity{Component: "time", ObjectID: fmt.Sprintf("sleeptimer_%d_end", timer), Name: fmt.Sprintf("Sleep Timer %d End", timer)},
				State: func(s spanet.Status) (string, bool) {
					if index >= len(s.SleepTimers) {
						return "", false
					}

					return s.SleepTimers[index].FinishTime.Format(timeLayout), true
				},
				raw: setterRaw(parseClock, func(spa *spanet.Spanet, v time.Time) error {
					_, err := spa.SetSleepTimerEnd(timer, v)

					return err
				}),
			},
		)
	}

	return es
}
