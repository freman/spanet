package server

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/freman/spanet/subcmd/server/middleware/safespa"
)

type serverCmd struct {
	spa    string
	listen string
}

func (*serverCmd) Name() string     { return "server" }
func (*serverCmd) Synopsis() string { return "A JSON bridge to your spalink" }
func (*serverCmd) Usage() string {
	return `server -spa ip:port -listen ip:port
`
}
func (s *serverCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.spa, "spa", "", "Spa host:port")
	f.StringVar(&s.listen, "listen", ":8080", "Listen host:port")
}

func (s *serverCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	e := echo.New()
	e.Use(middleware.Logger())

	safeSpa := safespa.New(safespa.WithAddr(s.spa))

	defineRoutes(e, safeSpa)

	e.Start(s.listen)

	return subcommands.ExitSuccess
}

func defineRoutes(e *echo.Echo, safeSpa *safespa.SafeSpa) {
	svc := service{
		spa: safeSpa,
	}

	api := e.Group("/spa", safeSpa.Mutex)

	// cmd_lights.go
	api.POST("/lights", svc.handlePostLights)
	api.GET("/lights/modes", svc.handleGetList(spanet.LightsModeStrings()))
	api.POST("/lights/mode", svc.handleSimplePost("SetLightsMode"))
	api.POST("/lights/brightness", svc.handleSimplePost("SetLightsBrightness"))
	api.POST("/lights/effectspeed", svc.handleSimplePost("SetLightsEffectSpeed"))
	api.POST("/lights/colour", svc.handleSimplePost("SetLightsColour"))
	api.POST("/lights/off", svc.handleSimplePost("SetLightsOff"))
	api.POST("/lights/toggle", svc.handleSimplePost("ToggleLights"))

	// cmd_mechanical.go
	api.POST("/pump/:pump", svc.handlePostPump)
	api.GET("/pump/states", svc.handleGetList(spanet.PumpStateStrings()))
	api.POST("/blower", svc.handlePostBlower)
	api.GET("/blower/modes", svc.handleGetList(spanet.BlowerModeStrings()))
	api.POST("/blower/speed", svc.handleSimplePost("SetBlowerVariableSpeed"))

	// cmd_misc.go
	api.POST("/temperature", svc.handleSimplePost("SetTargetTemperature"))
	api.GET("/operation/modes", svc.handleGetList(spanet.OperationModeStrings()))
	api.POST("/operation/mode", svc.handleSimplePost("SetOperationMode"))
	api.POST("/sanitise", svc.handleSimplePost("ToggleSanitise"))
	api.POST("/sanitise/time", svc.handlePostSanitiseTime)
	api.POST("/filtration/runtime", svc.handleSimplePost("SetFiltrationRunTime", "Hours"))
	api.POST("/filtration/cycle", svc.handleSimplePost("SetFiltrationCycle", "Hours"))
	api.POST("/timeout", svc.handleSimplePost("SetTimeout", "Minutes"))
	api.GET("/heatpump/modes", svc.handleGetList(spanet.HeatPumpModeStrings()))
	api.POST("/heatpump/mode", svc.handleSimplePost("SetHeatPumpMode"))
	api.POST("/svelementboost", svc.handleSimplePost("SetSVElementBoost"))
	api.GET("/lock/modes", svc.handleGetList(spanet.LockModeStrings()))
	api.POST("/lock/mode", svc.handleSimplePost("SetLockMode"))

	// cmd_power.go
	api.GET("/powersave/modes", svc.handleGetList(spanet.PowerSaveModeStrings()))
	api.POST("/powersave/mode", svc.handleSimplePost("SetPowerSave", "Mode"))
	api.POST("/peak/start", svc.handlePostPeakStart)
	api.POST("/peak/end", svc.handlePostPeakEnd)

	// cmd_sleep.go
	api.GET("/sleeptimer/states", svc.handleGetList(spanet.SleepTimerStateStrings()))
	api.POST("/sleeptimer/:timer/state", svc.handlePostSetSleepTimerState)
	api.POST("/sleeptimer/:timer/start", svc.handlePostSleepTimerStart)
	api.POST("/sleeptimer/:timer/end", svc.handlePostSleepTimerEnd)
	api.POST("/sleeptimer/:timer", svc.handlePostSleepTimer)

	// cmd_timedate.go
	api.POST("/datetime/year", svc.handleSimplePost("SetYear"))
	api.POST("/datetime/month", svc.handleSimplePost("SetMonth"))
	api.POST("/datetime/day", svc.handleSimplePost("SetDay"))
	api.POST("/datetime/hour", svc.handleSimplePost("SetHour"))
	api.POST("/datetime/minute", svc.handleSimplePost("SetMinute"))
	api.POST("/datetime/year", svc.handleSimplePost("SetYear"))
	api.POST("/datetime", svc.handlePostDateTime)

	api.GET("/status", svc.handleGetStatus)
}

func init() {
	subcommands.Register(&serverCmd{}, "")
}
