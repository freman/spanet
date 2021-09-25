package spanet

/*
RF:
,R2,18,250,51,70,4,[TimeDate.Hour],[TimeDate.Minute],55,[TimeDate.Day],[TimeDate.Month],[TimeDate.Year],376,9999,1,0,490,207,34,6000,602,23,20,0,0,0,0,44,35,45,:
,R3,32,1,4,4,4,SW V5 17 05 31,SV3,18480001,20000826,1,0,0,0,0,0,NA,7,0,470,Filtering,4,0,7,7,0,0,:
,R4,[OperationMode],0,0,0,1,0,3547,4,20,4500,7413,567,1686,0,8388608,0,0,5,0,98,0,10084,4,80,100,0,0,4,:
,R5,0,1,0,1,0,0,0,0,0,[Sleeping],[UVOzone],[Heating],[Auto],[Lights.On],[WaterTemperature],[Sanitise],3,[Pumps.1.State],[Pumps.2.State],[Pumps.3.State],[Pumps.4.State],[Pumps.5.State],0,1,2,6,:
,R6,[Blower.VariableSpeed],[Lights.Brightness],[Lights.Colour],[Lights.Mode],[Lights.Speed],[FiltrationHour],[FiltrationCycles],[SetTemperature],1,[PowerSave],[PeakStart],[PeakEnd],[SleepTimers.1.State],[SleepTimers.2.State],[SleepTimers.1.StartTime],[SleepTimers.2.StartTime],[SleepTimers.1.FinishTime],[SleepTimers.2.FinishTime],0,[Timeout],0,0,0,0,2,3,0,:
,R7,[AutoSanitise],0,1,1,1,0,1,0,0,0,253,191,253,240,483,125,77,1,0,0,0,23,200,1,[SVElementBoost],[HeatPumpMode],31,32,35,100,5,:
,R9,F1,255,0,0,0,0,0,0,0,0,0,0,:
,RA,F2,0,0,0,0,0,0,255,0,0,0,0,:
,RB,F3,0,0,0,0,0,0,0,0,0,0,0,:
,RC,0,1,1,0,0,0,0,0,0,[Blower.Mode],0,0,1,0,:
,RE,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,-4,13,30,8,5,1,0,0,0,0,0,:*
,RG,1,[Pumps.2.SwitchOn],[Pumps.3.SwitchOn],[Pumps.4.SwitchOn],[Pumps.5.SwitchOn],1,[Pumps.1.Installed]-[Pumps.1.SpeedType]-[Pumps.1.SpeedTypes],[Pumps.2.Installed]-[Pumps.2.SpeedType]-[Pumps.2.SpeedTypes],[Pumps.3.Installed]-[Pumps.3.SpeedType]-[Pumps.3.SpeedTypes],[Pumps.4.Installed]-[Pumps.4.SpeedType]-[Pumps.4.SpeedTypes],[Pumps.5.Installed]-[Pumps.5.SpeedType]-[Pumps.5.SpeedTypes],[Lock],:*
*/

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

func Parse(reader io.Reader) (status Status, err error) {
	parser := parser{Status: Status{
		Pumps:       make([]Pump, 5, 5),
		SleepTimers: make([]SleepTimer, 2, 2),
	}}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		switch line[0:3] {
		case ",R2":
			parser.parseR2(line[4:])
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
	return time.Date(0, 0, 0, i/256, i%256, 0, 0, time.Local)
}

func (p *parser) parseTime(i string) time.Time {
	return spa256toTime(int(p.parseUint(i)))
}

func (p *parser) parseR2(v string) {
	list := strings.Split(v, ",")
	p.TimeDate.Hour = p.parseByte(list[5])
	p.TimeDate.Minute = p.parseByte(list[6])
	p.TimeDate.Day = p.parseByte(list[8])
	p.TimeDate.Month = p.parseByte(list[9])
	p.TimeDate.Year = p.parseUint(list[10])
}

func (p *parser) parseR4(v string) {
	p.OperationMode = OperationMode(strings.Split(v, ",")[0])
}

func (p *parser) parseR5(v string) {
	list := strings.Split(v, ",")

	p.Sleeping = p.parseBool(list[9])
	p.UVOzone = p.parseBool(list[10])
	p.Heating = p.parseBool(list[11])
	p.Auto = p.parseBool(list[12])
	p.Lights.On = p.parseBool(list[13])
	p.WaterTemperature = Temperature(p.parseUint(list[14]))
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
	p.FiltrationCycles = p.parseByte(list[6])
	p.SetTemperature = Temperature(p.parseUint(list[7]))

	p.PowerSave = PowerSaveMode(p.parseByte(list[9]))

	p.PeakStart = p.parseTime(list[10])
	p.PeakEnd = p.parseTime(list[11])

	for i := 0; i < 2; i++ {
		p.SleepTimers[i].State = SleepTimerState(p.parseByte(list[12+i]))
		p.SleepTimers[i].StartTime = p.parseTime(list[14+i])
		p.SleepTimers[i].FinishTime = p.parseTime(list[16+i])
	}

	p.Timeout = time.Duration(p.parseUint(list[19])) * time.Minute
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
			p.Pumps[i].States = make([]PumpState, len(tmp[2]), len(tmp[2]))

			for n, v := range tmp[2] {
				p.Pumps[i].States[n] = PumpState(p.parseByte(string(v)))
			}
		}
	}

	p.Lock = LockMode(p.parseByte(list[11]))
}
