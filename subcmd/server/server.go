package server

import (
	"context"
	"flag"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/freman/spanet/subcmd/server/middleware/safespa"
	"github.com/google/subcommands"
	"github.com/labstack/echo/v4"
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

	safeSpa := safespa.New(s.spa)
	svc := service{
		spa: safeSpa,
	}

	api := e.Group("/spa", safeSpa.Mutex)

	// cmd_lights.go
	api.POST("/lights", svc.handlePostLights)
	api.GET("/lights/modes", svc.handleGetList(spanet.LightsModeNames()))
	api.POST("/lights/mode", svc.handleSimplePost("SetLightsMode"))
	api.POST("/lights/brightness", svc.handleSimplePost("SetLightsBrightness"))
	api.POST("/lights/effectspeed", svc.handleSimplePost("SetLightsEffectSpeed"))
	api.POST("/lights/colour", svc.handleSimplePost("SetLightsColour"))
	api.POST("/lights/off", svc.handleSimplePost("SetLightsOff"))
	api.POST("/lights/toggle", svc.handleSimplePost("ToggleLights"))

	// cmd_mechanical.go
	api.POST("/pump/:pump:", svc.handlePostPump)
	api.GET("/pump/states", svc.handleGetList(spanet.PumpStateNames()))
	api.POST("/blower", svc.handlePostBlower)
	api.GET("/blower/modes", svc.handleGetList(spanet.BlowerModeNames()))
	api.POST("/blower/speed", svc.handleSimplePost("SetBlowerVariableSpeed"))

	// cmd_misc.go
	api.POST("/temperature", svc.handleSimplePost("SetTargetTemperature"))
	api.GET("/operation/modes", svc.handleGetList(spanet.OperationModeNames()))
	api.POST("/operation/mode", svc.handleSimplePost("SetOperationMode"))
	api.POST("/sanitise", svc.handleSimplePost("ToggleSanitise"))
	api.POST("/sanitise/time", svc.handlePostSanitiseTime)
	api.POST("/filtration/runtime", svc.handleSimplePost("SetFiltrationRunTime", "Hours"))
	api.POST("/filtration/cycles", svc.handleSimplePost("SetFiltrationCycle", "Hours"))
	api.POST("/timeout", svc.handleSimplePost("SetTimeout", "Minutes"))
	api.GET("/heatpump/modes", svc.handleGetList(spanet.HeatPumpModeNames()))
	api.POST("/heatpump/mode", svc.handleSimplePost("SetHeatPumpMode"))
	api.POST("/svelementboost", svc.handleSimplePost("SetSVElementBoost"))
	api.GET("/lock/modes", svc.handleGetList(spanet.LockModeNames()))
	api.POST("/lock/mode", svc.handleSimplePost("SetLockMode"))

	api.GET("/status", svc.handleGetStatus)

	e.Start(s.listen)

	return subcommands.ExitSuccess
}

func init() {
	subcommands.Register(&serverCmd{}, "")
}
