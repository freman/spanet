package server

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/subcommands"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/freman/spanet/subcmd/server/middleware/safespa"
	"github.com/freman/spanet/subcmd/server/mqttbridge"
)

// gracefulTimeout bounds how long we'll wait for in-flight requests to
// drain after the first shutdown signal before giving up.
const gracefulTimeout = 10 * time.Second

type serverCmd struct {
	spa    string
	listen string

	mqttBroker       string
	mqttUsername     string
	mqttPassword     string
	mqttNodeID       string
	mqttPollInterval time.Duration
}

func (*serverCmd) Name() string     { return "server" }
func (*serverCmd) Synopsis() string { return "A JSON bridge to your spalink" }
func (*serverCmd) Usage() string {
	return `server -spa ip:port -listen ip:port [-mqtt-broker tcp://host:1883 ...]
`
}
func (s *serverCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.spa, "spa", "", "Spa host:port")
	f.StringVar(&s.listen, "listen", ":8080", "Listen host:port")

	f.StringVar(&s.mqttBroker, "mqtt-broker", "", "MQTT broker URL, e.g. tcp://localhost:1883 - enables Home Assistant MQTT discovery when set")
	f.StringVar(&s.mqttUsername, "mqtt-username", "", "MQTT username")
	f.StringVar(&s.mqttPassword, "mqtt-password", "", "MQTT password")
	f.StringVar(&s.mqttNodeID, "mqtt-node-id", "spanet", "Unique id for this spa in Home Assistant and MQTT topics")
	f.DurationVar(&s.mqttPollInterval, "mqtt-poll-interval", 15*time.Second, "How often to poll spa status for MQTT")
}

func (s *serverCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	e := echo.New()
	e.Use(middleware.RequestLogger())

	safeSpa := safespa.New(safespa.WithAddr(s.spa))
	defer func() {
		if err := safeSpa.Close(); err != nil {
			slog.Warn("failed to close spa connection", "error", err)
		}
	}()

	defineRoutes(e, safeSpa)

	// signal.NotifyContext's stop func is documented as possibly restoring
	// the OS default (kill) behavior for a signal once called, but in
	// practice a second ctrl+c after that isn't guaranteed to do anything -
	// so we watch for a second signal ourselves and force-exit on it.
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-sigCh
		slog.Info("shutting down, press ctrl+c again to force quit")
		cancel()

		<-sigCh
		slog.Warn("force quitting, spa connection left dangling until it drops")
		os.Exit(1)
	}()

	var mqttDone chan struct{}

	if s.mqttBroker != "" {
		bridge, err := mqttbridge.NewBridge(mqttbridge.Config{
			Broker:       s.mqttBroker,
			Username:     s.mqttUsername,
			Password:     s.mqttPassword,
			NodeID:       s.mqttNodeID,
			PollInterval: s.mqttPollInterval,
		}, safeSpa)
		if err != nil {
			slog.Error("failed to start mqtt bridge", "error", err)

			return subcommands.ExitFailure
		}

		mqttDone = make(chan struct{})

		go func() {
			defer close(mqttDone)

			if err := bridge.Start(ctx); err != nil {
				slog.Error("mqtt bridge stopped", "error", err)
			}
		}()
	}

	sc := echo.StartConfig{
		Address:         s.listen,
		GracefulTimeout: gracefulTimeout,
		OnShutdownError: func(err error) {
			slog.Warn("graceful shutdown did not complete in time", "error", err)
		},
	}

	httpErr := sc.Start(ctx, e)

	if mqttDone != nil {
		<-mqttDone
	}

	if httpErr != nil {
		slog.Error("server stopped", "error", httpErr)

		return subcommands.ExitFailure
	}

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
	api.POST("/datetime", svc.handlePostDateTime)

	api.GET("/status", svc.handleGetStatus)
}

func init() {
	subcommands.Register(&serverCmd{}, "")
}
