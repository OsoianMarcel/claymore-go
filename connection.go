package claymore

import (
	"net"
	"fmt"
	"io/ioutil"
)

// Connection
type Connection struct{
	ServerAddress string
}

// Create new instance of connection
func NewConnection(serverAddress string) Connection {
	return Connection{serverAddress}
}

// Send request to claymore server
func (c *Connection) Request(packet []byte) ([]byte, error) {
	emptyBytes := []byte{}

	tcpAddr, err := net.ResolveTCPAddr("tcp", c.ServerAddress)
	if err != nil {
		return emptyBytes, fmt.Errorf("resolve TCP address failed: %s", err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return emptyBytes, fmt.Errorf("dial failed: %s", err.Error())
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return emptyBytes, fmt.Errorf("write to server failed: %s", err.Error())
	}

	reply, err := ioutil.ReadAll(conn)
	if err != nil {
		return emptyBytes, fmt.Errorf("read from server failed: %s", err.Error())
	}

	return reply, nil
}