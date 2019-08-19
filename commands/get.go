package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/pooled"
)

var pool = pooled.NewMap(45*time.Second, 300*time.Millisecond, getConnection)

//GetInput .
func GetInput(address, output string) (string, error) {
	var input string

	work := func(conn pooled.Conn) error {
		conn.Log().Infof("Getting the current input for output %s", output)

		cmd := []byte(fmt.Sprintf("dumpdmrouteinfo\r\n"))
		_, err := conn.Write(cmd)
		if err != nil {
			return err
		}

		b, err := conn.ReadUntil(LINE_FEED, 5*time.Second)
		if err != nil {
			return err
		}

		var n []byte
		conn.Log().Debugf("Response from command: %s", b)

		for {
			n, err = conn.ReadUntil(LINE_FEED, 5*time.Second)
			if err != nil {
				return err
			}

			s := fmt.Sprintf("%s", n)

			if strings.Contains(s, "Routing Information for Output Card at Slot "+output) {
				for {
					n, err = conn.ReadUntil(LINE_FEED, 5*time.Second)
					if err != nil {
						return err
					}

					s := fmt.Sprintf("%s", n)

					if strings.Contains(s, "VideoSwitch") {
						input = strings.Split(s, "->In")[1]
						input = strings.TrimSuffix(input, "\r\n")
						log.L.Debugf("Response from command2: %s", n)
						break
					}
				}
				break
			}
		}

		return nil
	}

	err := pool.Do(address, work)
	if err != nil {
		return "", fmt.Errorf("failed to get input for output %s: %s", output, err)
	}
	return input, nil
}
