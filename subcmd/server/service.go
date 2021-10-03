package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/freman/spanet/subcmd/server/middleware/safespa"
	"github.com/labstack/echo/v4"
)

type service struct {
	spa *safespa.SafeSpa
}

func (s *service) handleGetStatus(c echo.Context) error {
	status, err := s.spa.GetStatus()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, status)
}

func (s *service) handleGetList(list []string) echo.HandlerFunc {
	blob, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}
	return func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, blob)
	}
}

func (s *service) handlePostLights(c echo.Context) error {
	var lights struct {
		Mode        spanet.LightsMode
		Brightness  int
		EffectSpeed int
		Colour      int
	}

	if err := c.Bind(&lights); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := s.spa.SetLightsMode(lights.Mode); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if _, err := s.spa.SetLightsEffectSpeed(lights.EffectSpeed); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if _, err := s.spa.SetLightsColour(lights.Colour); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if _, err := s.spa.SetLightsBrightness(lights.Brightness); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func (s *service) handlePostLightsMode(c echo.Context) error {
	var input struct {
		Mode spanet.LightsMode
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := s.spa.SetLightsMode(input.Mode); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostLightsBrightness(c echo.Context) error {
	var input struct {
		Brightness int
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := s.spa.SetLightsBrightness(input.Brightness); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostLightsEffectSpeed(c echo.Context) error {
	var input struct {
		EffectSpeed int
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := s.spa.SetLightsEffectSpeed(input.EffectSpeed); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostLightsColour(c echo.Context) error {
	var input struct {
		Colour int
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := s.spa.SetLightsColour(input.Colour); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostLightsOff(c echo.Context) error {
	if err := s.spa.SetLightsOff(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostToggleLights(c echo.Context) error {
	if err := s.spa.ToggleLights(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostPump(c echo.Context) error {
	var input struct {
		State spanet.PumpState
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	pump, err := strconv.ParseInt(c.Param("pump"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := s.spa.ControlPump(int(pump), input.State); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostBlower(c echo.Context) error {
	var input struct {
		Mode  spanet.BlowerMode
		Speed int
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := s.spa.ControlBlower(input.Mode); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if input.Speed > 0 {
		if _, err := s.spa.SetBlowerVariableSpeed(input.Speed); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostBlowerSpeed(c echo.Context) error {
	var input struct {
		Speed int
	}
	if _, err := s.spa.SetBlowerVariableSpeed(input.Speed); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}
