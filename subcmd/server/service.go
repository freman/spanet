package server

import (
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/freman/spanet/subcmd/server/middleware/safespa"
)

type service struct {
	spa *safespa.SafeSpa
}

func (s *service) handleGetStatus(c echo.Context) error {
	status, err := s.spa.GetStatus()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if _, err := s.spa.SetLightsMode(lights.Mode); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	if _, err := s.spa.SetLightsEffectSpeed(lights.EffectSpeed); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	if _, err := s.spa.SetLightsColour(lights.Colour); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	if _, err := s.spa.SetLightsBrightness(lights.Brightness); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	return nil
}

func (s *service) handlePostPump(c echo.Context) error {
	var input struct {
		State spanet.PumpState
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	pump, err := strconv.ParseInt(c.Param("pump"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if err := s.spa.ControlPump(int(pump), input.State); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostBlower(c echo.Context) error {
	var input struct {
		Mode  spanet.BlowerMode
		Speed int
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if err := s.spa.ControlBlower(input.Mode); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	if input.Speed > 0 {
		if _, err := s.spa.SetBlowerVariableSpeed(input.Speed); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
		}
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) parseTimeInput(c echo.Context) (time.Time, error) {
	var input struct {
		Time string
	}

	if err := c.Bind(&input); err != nil {
		return time.Time{}, echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	t, err := time.Parse("15:04", input.Time)
	if err != nil {
		return time.Time{}, echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	return t, nil
}

func (s *service) handlePostSanitiseTime(c echo.Context) error {
	t, err := s.parseTimeInput(c)
	if err != nil {
		return err
	}

	if _, err := s.spa.SetAutoSanitiseTime(t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostPeakStart(c echo.Context) error {
	t, err := s.parseTimeInput(c)
	if err != nil {
		return err
	}

	if _, err := s.spa.SetPeakStart(t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostPeakEnd(c echo.Context) error {
	t, err := s.parseTimeInput(c)
	if err != nil {
		return err
	}

	if _, err := s.spa.SetPeakEnd(t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostSetSleepTimerState(c echo.Context) error {
	var input struct {
		State spanet.SleepTimerState
	}

	timer, err := strconv.ParseInt(c.Param("timer"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if _, err := s.spa.SetSleepTimerState(int(timer), input.State); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostSleepTimerStart(c echo.Context) error {
	var input struct {
		Time timeParam
	}

	timer, err := strconv.ParseInt(c.Param("timer"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if _, err := s.spa.SetSleepTimerStart(int(timer), input.Time.Time); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostSleepTimerEnd(c echo.Context) error {
	var input struct {
		Time string
	}

	timer, err := strconv.ParseInt(c.Param("timer"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	t, err := time.Parse("15:04", input.Time)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if _, err := s.spa.SetSleepTimerEnd(int(timer), t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostSleepTimer(c echo.Context) error {
	var input struct {
		State spanet.SleepTimerState
		Start timeParam
		End   timeParam
	}

	timer, err := strconv.ParseInt(c.Param("timer"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if _, err := s.spa.SetSleepTimerState(int(timer), input.State); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	if _, err := s.spa.SetSleepTimerEnd(int(timer), input.Start.Time); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	if _, err := s.spa.SetSleepTimerEnd(int(timer), input.End.Time); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

func (s *service) handlePostDateTime(c echo.Context) error {
	var input struct {
		DateTime timeDateParam
	}

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error(), err)
	}

	if err := s.spa.SetDateTime(input.DateTime.Time); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
	}

	return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
}

var reSmarts = regexp.MustCompile(`[A-Z][^A-Z]*`)

func (s *service) handleSimplePost(fn string, name ...string) echo.HandlerFunc {
	rfn, found := reflect.TypeOf(s.spa.Spanet).MethodByName(fn)
	if !found {
		panic("no such method")
	}

	erridx := rfn.Type.NumOut() - 1

	if rfn.Type.NumIn() == 2 {
		sname := ""
		if len(name) != 0 {
			sname = name[0]
		} else {
			submatchall := reSmarts.FindAllString(fn, -1)
			sname = submatchall[len(submatchall)-1]
		}

		arg := rfn.Type.In(1)

		bindable := reflect.StructOf([]reflect.StructField{{
			Name: sname,
			Type: arg,
		}})

		return func(c echo.Context) error {
			v := reflect.New(bindable)

			if err := c.Bind(v.Interface()); err != nil {
				if err, isa := err.(*echo.HTTPError); isa {
					return err
				}

				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			out := reflect.ValueOf(s.spa.Spanet).MethodByName(fn).Call([]reflect.Value{v.Elem().FieldByName(sname)})
			if !out[erridx].IsNil() {
				err, isa := out[erridx].Interface().(error)
				if isa {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
				}

				return echo.NewHTTPError(http.StatusInternalServerError, out[erridx].Interface())
			}

			return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
		}
	}

	return func(c echo.Context) error {
		out := reflect.ValueOf(s.spa.Spanet).MethodByName(fn).Call(nil)
		if !out[erridx].IsNil() {
			err, isa := out[erridx].Interface().(error)
			if isa {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error(), err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, out[erridx].Interface())
		}

		return c.JSONBlob(http.StatusOK, []byte(`"ok"`))
	}
}
