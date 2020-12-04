package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/byuoitav/common/status"
	"github.com/byuoitav/crestron/cpu3/matrix"
	"github.com/labstack/echo"
)

type Handlers struct {
	CreateVideoSwitcher func(string) *matrix.Matrix
}

func (h *Handlers) RegisterRoutes(group *echo.Group) {
	sixteen := group.Group("/MD16x16/:address")

	// TODO singleflight?

	// get state
	sixteen.GET("/output/:output/input", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting inputs")

		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()

		inputs, err := vs.AudioVideoInputs(ctx)
		if err != nil {
			l.Printf("unable to get inputs: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		out := c.Param("output")
		in, ok := inputs[out]
		if !ok {
			l.Printf("invalid output %q requested", out)
			return c.String(http.StatusBadRequest, "invalid output")
		}

		l.Printf("Got inputs: %+v", inputs)
		return c.JSON(http.StatusOK, status.Input{
			Input: fmt.Sprintf("%v:%v", in, out),
		})
	})

	// set state
	sixteen.GET("/output/:output/input/:input", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
		out := c.Param("output")
		in := c.Param("input")

		l.Printf("Setting AV input on %q to %q", out, in)

		ctx, cancel := context.WithTimeout(c.Request().Context(), 15*time.Second)
		defer cancel()

		err := vs.SetAudioVideoInput(ctx, out, in)
		if err != nil {
			l.Printf("unable to set AV input: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set AV input")
		return c.JSON(http.StatusOK, status.Input{
			Input: fmt.Sprintf("%v:%v", in, out),
		})
	})
}
