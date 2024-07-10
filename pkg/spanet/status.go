package spanet

//go:generate enumer -sql=false -linecomment -type=PumpState,BlowerMode,LightsMode,PowerSaveMode,HeatPumpMode,LockMode,SleepTimerState -output=status_strings.go $GOFILE

import (
	"encoding/json"
	"time"
)

type PumpState byte
type BlowerMode byte
type LightsMode byte
type OperationMode byte
type SleepTimerState byte
type PowerSaveMode byte
type HeatPumpMode byte
type LockMode byte
type Timeout struct {
	time.Duration
}

type Status struct {
	SetTemperature    float64 `rf:"R6,9"`
	WaterTemperature  float64 `rf:"R5,16"`
	HeaterTemperature float64 `rf:"R2,12"`
	CaseTemperature   float64 `rf:"R2,4"`
	WaterPresent      bool    `rf:"R2:15"`
	Heating           bool    `rf:"R5,13"`
	UVOzone           bool    `rf:"R5,12"`
	Sanitise          bool    `rf:"R5,17"`
	Auto              bool    `rf:"R5,14"`
	Sleeping          bool    `rf:"R5,11"`
	Pumps             []Pump
	Blower            Blower
	Lights            Lights

	OperationMode   OperationMode `rf:"R4,2"`
	FiltrationHour  byte          `rf:"R6,7"`
	FiltrationCycle byte          `rf:"R6,8"` // 1, 2, 3, 4, 6, 8, 12 and 24
	SleepTimers     []SleepTimer

	PowerSave      PowerSaveMode `rf:"R6,11"`
	PeakStart      time.Time     `rf:"R5,12"`
	PeakEnd        time.Time     `rf:"R6,13"`
	AutoSanitise   time.Time     `rf:"R7,2"`
	Timeout        Timeout       `rf:"R6,21"` // minutes
	HeatPumpMode   HeatPumpMode  `rf:"R7,27"`
	SVElementBoost bool          `rf:"R7,26"`

	TimeDate TimeDate

	Lock LockMode `rf:"RG,13"`

	AwakeRemains           int     `rf:"R2,17"` // unsure
	FilterPumpTotalRunTime int     `rf:"R2,18"` // minutes
	FilterPumpReq          int     `rf:"R2,19"` //unsure
	RuntimeHours           float64 `rf:"R2,22"` // * 10

	Power Power
}

type Pump struct {
	State     PumpState   `rf:"R5,19|R5,20|R5,21|R5,22|R5,23"`
	Installed bool        `rf:"RG,8-0|RG,9-0|RG,10-0|RG,11-0|RG,12-0"`
	SpeedType byte        `rf:"RG,8-1|RG,9-1|RG,10-1|RG,11-1|RG,12-1"`
	States    []PumpState `rf:"RG,8-2|RG,9-2|RG,10-2|RG,11-2|RG,12-2"` //Data: First part (1- or 0-) indicates whether the pump is installed/fitted. If so (1- means it is), the second part (1- above) indicates it's speed type. The third part (014 above) represents it's possible states (0 OFF, 1 ON, 4 AUTO)
	SwitchOn  bool        `rf:"RG,3|RG,4|RG,5|RG,6"`
}

type Blower struct {
	Mode          BlowerMode `rf:"RC,11"`
	VariableSpeed byte       `rf:"R6,2"`
}

type Lights struct {
	On         bool       `rf:"R5,15"`
	Mode       LightsMode `rf:"R6,5"`
	Brightness byte       `rf:"R5,3"`
	Speed      byte       `rf:"R6,6"`
	Colour     byte       `rf:"R6,5"`
}

type SleepTimer struct {
	State      SleepTimerState `rf:"R6,14|R6,15"`
	StartTime  time.Time       `rf:"R6,16|R6,17"`
	FinishTime time.Time       `rf:"R6,18|R6,19"`
}

type TimeDate struct {
	Hour   byte `rf:"R2,7"`
	Minute byte `rf:"R2,8"`
	Day    byte `rf:"R2,10"`
	Month  byte `rf:"R2,11"`
	Year   uint `rf:"R2,12"`
}

type Power struct {
	Volts        int     `rf:"R2,3"`
	Amps         int     `rf:"R2,2"`
	HeatingAmps  float64 `rf:"R3,23"`
	CurrentLimit int     `rf:"R3,2"`
	LoadShed     int     `rf:"R3,18"`
}

const (
	PumpStateOff  PumpState = 0 // Off
	PumpStateOn   PumpState = 1 // On
	PumpStateAuto PumpState = 4 // Auto
)

const (
	BlowerModeVariable BlowerMode = 0 // Variable
	BlowerModeRamp     BlowerMode = 1 // Ramp
	BlowerModeOff      BlowerMode = 2 // Off
)

const (
	LightsModeWhite  LightsMode = 0 // White
	LightsModeColour LightsMode = 1 // Colour
	LightsModeStep   LightsMode = 2 // Step
	LightsModeFade   LightsMode = 3 // Fade
	LightsModeParty  LightsMode = 4 // Party
)

const (
	OperationModeNormal   OperationMode = 0 // NORM
	OperationModeEconomy  OperationMode = 1 // ECON
	OperationModeAway     OperationMode = 2 // AWAY
	OperationModeWeekdays OperationMode = 3 // WEEK
)

const (
	PowerSaveModeOff  PowerSaveMode = 0 // Off
	PowerSaveModeLow  PowerSaveMode = 1 // Low
	PowerSaveModeHigh PowerSaveMode = 2 // High
)

const (
	HeatPumpModeAuto    HeatPumpMode = 0 // Auto
	HeatPumpModeHeat    HeatPumpMode = 1 // Heat
	HeatPumpModeCool    HeatPumpMode = 2 // Cool
	HeatPumpModeDisable HeatPumpMode = 3 // Disable
)

const (
	LockModeOff     LockMode = 0 // Off
	LockModePartial LockMode = 1 // Partial
	LockModeFull    LockMode = 2 // Full
)

const (
	SleepTimerStateWeekDays SleepTimerState = 31  // Weekdays
	SleepTimerStateWeekEnds SleepTimerState = 96  // Weekends
	SleepTimerStateEveryDay SleepTimerState = 127 // Everyday
	SleepTimerStateOff      SleepTimerState = 128 // Off
)

func (t TimeDate) AsTime() time.Time {
	return time.Date(
		int(t.Year),
		time.Month(t.Month),
		int(t.Day),
		int(t.Hour),
		int(t.Minute),
		0,
		0,
		time.UTC,
	)
}

func (t TimeDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.AsTime())
}

func (t *TimeDate) UnmarshalJSON(b []byte) error {
	var tmp time.Time
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	t.Year = uint(tmp.Year())
	t.Month = byte(tmp.Month())
	t.Day = byte(tmp.Day())
	t.Hour = byte(tmp.Hour())
	t.Minute = byte(tmp.Minute())

	return nil
}

func (t Timeout) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Duration.Round(time.Minute).Minutes())
}

func (t *Timeout) UnmarshalJSON(b []byte) error {
	var tmp uint
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	t.Duration = time.Duration(tmp) * time.Minute

	return nil
}
