// Code generated by "stringer -type=PumpState,BlowerMode,LightsMode,PowerSaveMode,HeatPumpMode,LockMode,SleepTimerState -linecomment -output=status_strings.go meh.go"; DO NOT EDIT.

package spanet

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PumpStateOff-0]
	_ = x[PumpStateOn-1]
	_ = x[PumpStateAuto-4]
}

const (
	_PumpState_name_0 = "OffOn"
	_PumpState_name_1 = "Auto"
)

var (
	_PumpState_index_0 = [...]uint8{0, 3, 5}
)

func (i PumpState) String() string {
	switch {
	case i <= 1:
		return _PumpState_name_0[_PumpState_index_0[i]:_PumpState_index_0[i+1]]
	case i == 4:
		return _PumpState_name_1
	default:
		return "PumpState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BlowerModeVariable-0]
	_ = x[BlowerModeRamp-1]
	_ = x[BlowerModeOff-2]
}

const _BlowerMode_name = "VariableRampOff"

var _BlowerMode_index = [...]uint8{0, 8, 12, 15}

func (i BlowerMode) String() string {
	if i >= BlowerMode(len(_BlowerMode_index)-1) {
		return "BlowerMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _BlowerMode_name[_BlowerMode_index[i]:_BlowerMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LightsModeWhite-0]
	_ = x[LightsModeColour-1]
	_ = x[LightsModeStep-2]
	_ = x[LightsModeFade-3]
	_ = x[LightsModeParty-4]
}

const _LightsMode_name = "WhiteColourStepFadeParty"

var _LightsMode_index = [...]uint8{0, 5, 11, 15, 19, 24}

func (i LightsMode) String() string {
	if i >= LightsMode(len(_LightsMode_index)-1) {
		return "LightsMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _LightsMode_name[_LightsMode_index[i]:_LightsMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PowerSaveModeOff-0]
	_ = x[PowerSaveModeLow-1]
	_ = x[PowerSaveModeHigh-2]
}

const _PowerSaveMode_name = "OffLowHigh"

var _PowerSaveMode_index = [...]uint8{0, 3, 6, 10}

func (i PowerSaveMode) String() string {
	if i >= PowerSaveMode(len(_PowerSaveMode_index)-1) {
		return "PowerSaveMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PowerSaveMode_name[_PowerSaveMode_index[i]:_PowerSaveMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[HeatPumpModeAuto-0]
	_ = x[HeatPumpModeHeat-1]
	_ = x[HeatPumpModeCool-2]
	_ = x[HeatPumpModeDisable-3]
}

const _HeatPumpMode_name = "AutoHeatCoolDisable"

var _HeatPumpMode_index = [...]uint8{0, 4, 8, 12, 19}

func (i HeatPumpMode) String() string {
	if i >= HeatPumpMode(len(_HeatPumpMode_index)-1) {
		return "HeatPumpMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _HeatPumpMode_name[_HeatPumpMode_index[i]:_HeatPumpMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LockModeOff-0]
	_ = x[LockModePartial-1]
	_ = x[LockModeFull-2]
}

const _LockMode_name = "OffPartialFull"

var _LockMode_index = [...]uint8{0, 3, 10, 14}

func (i LockMode) String() string {
	if i >= LockMode(len(_LockMode_index)-1) {
		return "LockMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _LockMode_name[_LockMode_index[i]:_LockMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SleepTimerStateWeekDays-31]
	_ = x[SleepTimerStateWeekEnds-96]
	_ = x[SleepTimerStateEveryDay-127]
	_ = x[SleepTimerStateOff-128]
}

const (
	_SleepTimerState_name_0 = "Weekdays"
	_SleepTimerState_name_1 = "Weekends"
	_SleepTimerState_name_2 = "EverydayOff"
)

var (
	_SleepTimerState_index_2 = [...]uint8{0, 8, 11}
)

func (i SleepTimerState) String() string {
	switch {
	case i == 31:
		return _SleepTimerState_name_0
	case i == 96:
		return _SleepTimerState_name_1
	case 127 <= i && i <= 128:
		i -= 127
		return _SleepTimerState_name_2[_SleepTimerState_index_2[i]:_SleepTimerState_index_2[i+1]]
	default:
		return "SleepTimerState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
