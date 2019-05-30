package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/pooled"
)

//SwitchInput takes the IP address, the output and the input from the user and
//switches the input to the one requested
func SwitchInput(address, output, input string) (string, error) {

	work := func(conn pooled.Conn) error {

		log.L.Infof("Setting output %s to %s", output, input)

		// _, err := conn.Write([]byte("\n"))
		// if err != nil {
		// 	return err
		// }

		cmd := []byte(fmt.Sprintf("setAVroute %s %s\r\n", input, output))

		// send command
		_, err := conn.Write(cmd)
		if err != nil {
			return err
		}

		var b []byte

		b, err = conn.ReadUntil(BIG_GATOR, 10*time.Second)
		if err != nil {
			return err
		}

		log.L.Infof("Response from command: %s", b)
		if strings.Contains(string(b), "failed") {
			return fmt.Errorf("input or output is out of range")
		}

		log.L.Infof("Set input to %s successful", input)
		return nil
	}

	err := pool.Do(address, work)
	if err != nil {
		return "", fmt.Errorf("Failed to switch input: %s", err)
	}

	return input, nil
}
