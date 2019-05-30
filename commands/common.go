package commands

import (
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/pooled"
)

const (
	CARRIAGE_RETURN           = 0x0D
	LINE_FEED                 = 0x0A
	SPACE                     = 0x20
	DELAY_BETWEEN_CONNECTIONS = time.Second * 10
	BIG_GATOR                 = 0x3e
)

//var tlsConfig *tls.Config
//var caller telnet.Caller

//func init() {
//	tlsConfig = &tls.Config{}
//}

func getConnection(key interface{}) (pooled.Conn, error) {
	address, ok := key.(string)
	if !ok {
		return nil, fmt.Errorf("key must be a string")
	}

	conn, err := net.DialTimeout("tcp", address+":23", 30*time.Second)
	if err != nil {
		return nil, err
	}

	pconn := pooled.Wrap(conn)

	// send command
	time.Sleep(50 * time.Millisecond) // time for the switcher to chill out
	// read first new line
	var b []byte
	b, err = pconn.ReadUntil(LINE_FEED, 3*time.Second)
	if err != nil {
		return nil, err
	}

	log.L.Infof("Response from command: %s", b)

	return pconn, nil
}
