package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/pooled"
)

//SwitchInput takes the IP address, the output and the input from the user and
//switches the input to the one requested
func SwitchInput(address, output, input string) (string, error) {

	work := func(conn pooled.Conn) error {
		conn.Log().Infof("Setting output %s to %s", output, input)

		cmd := []byte(fmt.Sprintf("setAVroute %s %s\r\n", input, output))
		_, err := conn.Write(cmd)
		if err != nil {
			return err
		}

		var b []byte

		b, err = conn.ReadUntil(BIG_GATOR, 10*time.Second)
		if err != nil {
			return err
		}

		conn.Log().Debugf("Response from command: %s", b)
		if strings.Contains(string(b), "failed") {
			return fmt.Errorf("input or output is out of range")
		}

		conn.Log().Infof("Set input to %s successful", input)
		return nil
	}

	err := pool.Do(address, work)
	if err != nil {
		return "", fmt.Errorf("Failed to switch input: %s", err)
	}

	return input, nil
}
