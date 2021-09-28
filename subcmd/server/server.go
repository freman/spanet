package server

import (
	"context"
	"flag"

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

	api.GET("/status", svc.handleGetStatus)
	api.POST("/lights", svc.handlePostLights)
	api.GET("/lights/modes", svc.handleGetLightsModes)
	api.POST("/lights/mode", svc.handlePostLightsMode)
	api.POST("/lights/brightness", svc.handlePostLightsBrightness)
	api.POST("/lights/effectspeed", svc.handlePostLightsEffectSpeed)
	api.POST("/lights/colour", svc.handlePostLightsColour)
	api.POST("/lights/off", svc.handlePostLightsOff)
	api.POST("/lights/toggle", svc.handlePostToggleLights)

	e.Start(s.listen)

	return subcommands.ExitSuccess
}

func init() {
	subcommands.Register(&serverCmd{}, "")
}
