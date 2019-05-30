package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/crestron-switcher-ms/commands"
	"github.com/labstack/echo"
)

//SwitchInput .
func SwitchInput(ectx echo.Context) error {
	address := ectx.Param("address")
	output := ectx.Param("output")
	input := ectx.Param("input")

	log.L.Infof("Switching input for output %s to %s", output, input)

	input, err := commands.SwitchInput(address, output, input)

	if err != nil {
		log.L.Errorf("Tis the error you receiveth: %s", err)
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	log.L.Infof("Successfully changed input for output %s to %s", output, input)
	return ectx.JSON(http.StatusOK, status.Input{
		Input: fmt.Sprintf("%v:%v", input, output),
	})
}

//GetInput .
func GetInput(ectx echo.Context) error {
	address := ectx.Param("address")
	output := ectx.Param("output")

	log.L.Infof("Getting input for output %s", output)

	input, err := commands.GetInput(address, output)

	if err != nil {
		return ectx.String(http.StatusInternalServerError, err.Error())
	}

	log.L.Infof("Input for output %v is %v", output, input)
	return ectx.JSON(http.StatusOK, status.Input{
		Input: fmt.Sprintf("%v:%v", input, output),
	})
}
