package spanet

import (
	"fmt"
	"strconv"
)

// SetLightsMode on the controller
// Command: S07:x
func (s *Spanet) SetLightsMode(mode LightsMode) (LightsMode, error) {
	r, err := s.command(fmt.Sprintf("S07:%d", mode))
	if err != nil {
		return 0, err
	}

	tmp, err := strconv.ParseInt(r[0:0], 10, 64)
	if err != nil {
		return 0, err
	}

	return LightsMode(tmp), nil
}

// SetLightsBrightness on the controller
// Valid: 1-5
// Command: S08:x
func (s *Spanet) SetLightsBrightness(brightness int) (int, error) {
	return s.commandInt("S08", brightness, 1, 5, "brightness")
}

// SetLightsEffectSpeed on the controller
// Valid: 1-5
// Command: S09:x
func (s *Spanet) SetLightsEffectSpeed(speed int) (int, error) {
	return s.commandInt("S09", speed, 1, 5, "speed")
}

// SetLightsColour on the controller
// Valid: 1-30
// Command: S10:x
func (s *Spanet) SetLightsColour(colour int) (int, error) {
	return s.commandInt("S10", colour, 1, 30, "colour")
}

// SetLightsOff on the controller
// Command: S11
func (s *Spanet) SetLightsOff() error {
	_, err := s.commandExpect("S11", "S11")
	return err
}

// ToggleLights on the controller
// Command: W14
func (s *Spanet) ToggleLights() error {
	_, err := s.commandExpect("W14", "W14")
	return err
}
