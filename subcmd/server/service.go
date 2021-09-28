package server

import (
	"net/http"

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

func (s *service) handlePostLights(c echo.Context) error {
	var lights struct {
		Mode        spanet.LightsMode
		Brightness  int
		EffectSpeed int
		Colour      int
	}

	if err := c.Bind(&lights); err != nil {
		return err
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

func (s *service) handleGetLightsModes(c echo.Context) error {
	return c.JSON(http.StatusOK, spanet.LightsModeNames())
}

func (s *service) handlePostLightsMode(c echo.Context) error {
	var input struct {
		Mode spanet.LightsMode
	}

	if err := c.Bind(&input); err != nil {
		return err
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
		return err
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
		return err
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
		return err
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
