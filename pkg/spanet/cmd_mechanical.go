package spanet

import (
	"errors"
	"fmt"
)

func (s *Spanet) ControlPump(pump int, state PumpState) error {
	if pump < 1 || pump > 5 {
		return ErrValueOutOfRange{1, 5, pump, "pump"}
	}

	if pump > 1 && state == PumpStateAuto {
		return errors.New("only pump one supports auto")
	}

	cmd := fmt.Sprintf("S%d", 21+pump)

	_, err := s.commandExpect(fmt.Sprintf("%s:%d", cmd, state), cmd+"-OK")
	if err != nil {
		return err
	}

	return nil
}

func (s *Spanet) ControlBlower(mode BlowerMode) error {
	_, err := s.commandExpect(fmt.Sprintf("S28:%d", mode), "S28-OK")
	return err
}

func (s *Spanet) SetBlowerVariableSpeed(speed int) (int, error) {
	return s.commandInt("S31", speed, 1, 5, "speed")
}
