package stats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
)

const (
	bufferSize = 2048
)

// ServerConfig is used to initialize a new ServerStats object
type ServerConfig struct {
	Domain string
	Port   int
}

// ServerStats is the object used to receive, store and send data for usage
type ServerStats struct {
	addr string
}

// NewServer create a new server object for use
func NewServer(config *ServerConfig) (*ServerStats, error) {
	return &ServerStats{
		addr: config.Domain + ":" + strconv.Itoa(config.Port),
	}, nil
}

// Run starts receiving the profiling stats for storage and usage
func (s *ServerStats) Run() {

	var serverAddr *net.UDPAddr
	var server *net.UDPConn
	var err error

	serverAddr, err = net.ResolveUDPAddr(udp, s.addr)
	if err != nil {
		panic(err)
	}

	server, err = net.ListenUDP(udp, serverAddr)
	if err != nil {
		panic(err)
	}
	defer server.Close()

	server.SetReadBuffer(bufferSize)

	var addr *net.UDPAddr
	var buf bytes.Buffer
	var bytesRead int
	buff := make([]byte, bufferSize)
	stats := new(Stats)

	decoder := gob.NewDecoder(&buf)

	for {

		bytesRead, addr, err = server.ReadFromUDP(buff)
		if err != nil {
			fmt.Println("Error:", err)
			buf.Reset()
			continue
		}

		buf.Write(buff)
		err = decoder.Decode(stats)
		buf.Reset()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Printf("Recieved: %v from %s Read %d bytes\n", stats.MemStats, addr, bytesRead)
	}
}
