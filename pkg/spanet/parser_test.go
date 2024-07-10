package spanet_test

import (
	"strings"
	"testing"
	"time"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/stretchr/testify/assert"
)

func TestParseStatus(t *testing.T) {
	testRF := `RF:
,R2,0,235,48,70,1,18,25,6,9,7,2024,410,9999,1,0,0,199,0,6000,195498,1682,0,0,0,0,0,14191,60029,4852,0,:
,R3,15,1,4,4,4,SW V6 21 12 13,SV2,21240001,20000755,1,0,0,0,0,0,NA,7,0,484,In use,23,0,7,7,0,0,0,:
,R4,NORM,0,0,0,1,413,0,4,44,0,363892,666,1162,0,0,0,0,0,0,101,0,418,4,80,100,0,0,4,:
,R5,0,1,0,0,0,0,0,0,0,0,0,0,1,0,412,0,48,4,0,0,0,0,0,1,6,6,:
,R6,5,5,0,2,5,3,3,410,1,2,4096,2816,127,128,4126,5632,2846,1792,0,130,0,0,0,0,2,2,0,410,:
,R7,3072,0,1,0,1,0,1,27,8,2022,251,218,242,229,487,125,77,1,0,0,0,23,181,1,0,1,31,50,50,100,5,:
,R9,F1,255,0,0,0,0,0,0,0,0,0,0,:
,RA,F2,0,0,0,0,60160,238,255,0,0,0,0,:
,RB,F3,0,0,0,0,129,0,0,0,0,0,0,:
,RC,0,0,0,0,0,0,0,0,0,2,0,0,0,0,:
,RE,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,-4,13,30,8,5,1,:*
,RG,1,1,1,1,1,1,1-1-014,1-1-01,0-,0-,0-,0,0,0,3001,:*`

	expected := spanet.Status{
		SetTemperature:    41.0,
		WaterTemperature:  41.2,
		HeaterTemperature: 41.0,
		CaseTemperature:   48.0,
		WaterPresent:      true,
		Heating:           false,
		UVOzone:           false,
		Sanitise:          false,
		Auto:              true,
		Sleeping:          false,
		Pumps: []spanet.Pump{{
			State:     spanet.PumpStateAuto,
			Installed: true,
			SpeedType: 1,
			States: []spanet.PumpState{
				spanet.PumpStateOff,
				spanet.PumpStateOn,
				spanet.PumpStateAuto,
			},
			SwitchOn: true,
		}, {
			State:     spanet.PumpStateOff,
			Installed: true,
			SpeedType: 1,
			States: []spanet.PumpState{
				spanet.PumpStateOff,
				spanet.PumpStateOn,
			},
			SwitchOn: true,
		}, {
			State:     spanet.PumpStateOff,
			Installed: false,
			SpeedType: 0,
			States:    nil,
			SwitchOn:  true,
		}, {
			State:     spanet.PumpStateOff,
			Installed: false,
			SpeedType: 0,
			States:    nil,
			SwitchOn:  true,
		}, {
			State:     spanet.PumpStateOff,
			Installed: false,
			SpeedType: 0,
			States:    nil,
			SwitchOn:  true,
		}},
		Blower: spanet.Blower{
			Mode:          spanet.BlowerModeOff,
			VariableSpeed: 5,
		},
		Lights: spanet.Lights{
			On:         false,
			Mode:       spanet.LightsModeStep,
			Brightness: 5,
			Speed:      5,
			Colour:     0,
		},
		OperationMode:   spanet.OperationModeNormal,
		FiltrationHour:  3,
		FiltrationCycle: 3,
		SleepTimers: []spanet.SleepTimer{{
			State:      spanet.SleepTimerStateEveryDay,
			StartTime:  time.Date(0, 11, 30, 16, 30, 0, 0, time.UTC),
			FinishTime: time.Date(0, 11, 30, 11, 30, 0, 0, time.UTC),
		}, {
			State:      spanet.SleepTimerStateOff,
			StartTime:  time.Date(0, 11, 30, 22, 00, 0, 0, time.UTC),
			FinishTime: time.Date(0, 11, 30, 7, 00, 0, 0, time.UTC),
		}},
		PowerSave:      spanet.PowerSaveModeHigh,
		PeakStart:      time.Date(0, 11, 30, 16, 00, 0, 0, time.UTC),
		PeakEnd:        time.Date(0, 11, 30, 11, 00, 0, 0, time.UTC),
		AutoSanitise:   time.Date(0, 11, 30, 12, 00, 0, 0, time.UTC),
		Timeout:        spanet.Timeout{130 * time.Minute},
		HeatPumpMode:   spanet.HeatPumpModeHeat,
		SVElementBoost: false,
		TimeDate: spanet.TimeDate{
			Hour:   18,
			Minute: 25,
			Day:    9,
			Month:  7,
			Year:   2024,
		},
		Lock:                   spanet.LockModeOff,
		AwakeRemains:           0,
		FilterPumpTotalRunTime: 199,
		FilterPumpReq:          0,
		RuntimeHours:           168.2,
		Power: spanet.Power{
			Volts:        235,
			Amps:         0,
			CurrentLimit: 15,
			LoadShed:     7,
		},
	}

	status, err := spanet.ParseStatus(strings.NewReader(testRF))
	assert.NoError(t, err)
	assert.Equal(t, expected, status)
}
