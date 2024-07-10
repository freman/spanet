package spanet

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

func ParseStatus(reader io.Reader) (status Status, err error) {
	parser := parser{Status: Status{
		Pumps:       make([]Pump, 5),
		SleepTimers: make([]SleepTimer, 2),
	}}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		switch line[0:3] {
		case ",R2":
			parser.parseR2(line[4:])
		case ",R3":
			parser.parseR3(line[4:])
		case ",R4":
			parser.parseR4(line[4:])
		case ",R5":
			parser.parseR5(line[4:])
		case ",R6":
			parser.parseR6(line[4:])
		case ",R7":
			parser.parseR7(line[4:])
		case ",RC":
			parser.parseRC(line[4:])
		case ",RG":
			parser.parseRG(line[4:])
			return parser.Status, parser.err
		}

	}

	return parser.Status, parser.err
}

type parser struct {
	Status
	err error
}

func (p *parser) _parseUint(i string, base int, bitSize int) (o uint64) {
	if p.err != nil {
		return 0
	}

	o, p.err = strconv.ParseUint(i, base, bitSize)

	return
}

func (p *parser) parseBool(i string) (o bool) {
	if p.err != nil {
		return false
	}

	o, p.err = strconv.ParseBool(i)

	return
}

func (p *parser) parseByte(i string) byte {
	return byte(p._parseUint(i, 10, 8))
}

func (p *parser) parseUint(i string) uint {
	return uint(p._parseUint(i, 10, 32))
}

func spa256toTime(i int) time.Time {
	return time.Date(1, 0, 0, i/256, i%256, 0, 0, time.UTC)
}

func (p *parser) parseTime(i string) time.Time {
	return spa256toTime(int(p.parseUint(i)))
}

func (p *parser) parseR2(v string) {
	list := strings.Split(v, ",")
	p.Power.Amps = int(p.parseUint(list[0]))
	p.Power.Volts = int(p.parseUint(list[1]))
	p.CaseTemperature = float64(p.parseUint(list[2]))
	p.TimeDate.Hour = p.parseByte(list[5])
	p.TimeDate.Minute = p.parseByte(list[6])
	p.TimeDate.Day = p.parseByte(list[8])
	p.TimeDate.Month = p.parseByte(list[9])
	p.TimeDate.Year = p.parseUint(list[10])
	p.HeaterTemperature = float64(p.parseUint(list[11])) / 10.0
	p.WaterPresent = p.parseBool(list[13])
	p.AwakeRemains = int(p.parseUint(list[15]))
	p.FilterPumpTotalRunTime = int(p.parseUint(list[16]))
	p.FilterPumpReq = int(p.parseUint(list[17]))
	p.RuntimeHours = float64(p.parseUint(list[20])) / 10.0
}

func (p *parser) parseR3(v string) {
	list := strings.Split(v, ",")
	p.Power.HeatingAmps = float64(p.parseUint(list[21])) / 10.0
	p.Power.CurrentLimit = int(p.parseUint(list[0]))
	p.Power.LoadShed = int(p.parseUint(list[16]))
}

func (p *parser) parseR4(v string) {
	p.OperationMode, _ = OperationModeString(strings.Split(v, ",")[0])
}

func (p *parser) parseR5(v string) {
	list := strings.Split(v, ",")

	p.Sleeping = p.parseBool(list[9])
	p.UVOzone = p.parseBool(list[10])
	p.Heating = p.parseBool(list[11])
	p.Auto = p.parseBool(list[12])
	p.Lights.On = p.parseBool(list[13])
	p.WaterTemperature = float64(p.parseUint(list[14])) / 10.0
	p.Sanitise = p.parseBool(list[15])

	for i := 0; i < 5; i++ {
		p.Pumps[i].State = PumpState(p.parseByte(list[17+i]))
	}
}

func (p *parser) parseR6(v string) {
	list := strings.Split(v, ",")

	p.Blower.VariableSpeed = p.parseByte(list[0])
	p.Lights.Brightness = p.parseByte(list[1])
	p.Lights.Colour = p.parseByte(list[2])
	p.Lights.Mode = LightsMode(p.parseByte(list[3]))
	p.Lights.Speed = p.parseByte(list[4])
	p.FiltrationHour = p.parseByte(list[5])
	p.FiltrationCycle = p.parseByte(list[6])
	p.SetTemperature = float64(p.parseUint(list[7])) / 10.0

	p.PowerSave = PowerSaveMode(p.parseByte(list[9]))
	p.PeakStart = p.parseTime(list[10])
	p.PeakEnd = p.parseTime(list[11])

	for i := 0; i < 2; i++ {
		p.SleepTimers[i].State = SleepTimerState(p.parseByte(list[12+i]))
		p.SleepTimers[i].StartTime = p.parseTime(list[14+i])
		p.SleepTimers[i].FinishTime = p.parseTime(list[16+i])
	}

	p.Timeout.Duration = time.Duration(p.parseUint(list[19])) * time.Minute
}

func (p *parser) parseR7(v string) {
	list := strings.Split(v, ",")

	p.AutoSanitise = p.parseTime(list[0])
	p.SVElementBoost = p.parseBool(list[24])
	p.HeatPumpMode = HeatPumpMode(p.parseByte(list[25]))
}

func (p *parser) parseRC(v string) {
	p.Blower.Mode = BlowerMode(p.parseByte(strings.Split(v, ",")[9]))
}

func (p *parser) parseRG(v string) {
	list := strings.Split(v, ",")

	for i := 0; i < 5; i++ {
		p.Pumps[i].SwitchOn = true

		if i > 0 {
			p.Pumps[i].SwitchOn = p.parseBool(list[0+i])
		}

		tmp := strings.Split(list[6+i], "-")

		p.Pumps[i].Installed = p.parseBool(tmp[0])
		if p.Pumps[i].Installed {
			p.Pumps[i].SpeedType = p.parseByte(tmp[1])
			p.Pumps[i].States = make([]PumpState, len(tmp[2]))

			for n, v := range tmp[2] {
				p.Pumps[i].States[n] = PumpState(p.parseByte(string(v)))
			}
		}
	}

	p.Lock = LockMode(p.parseByte(list[11]))
}
