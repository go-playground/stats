package stats

import (
	"encoding/gob"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"time"
)

// ClientConfig is used to initialize a new ClientStats object
type ClientConfig struct {
	Domain string
	Port   int
}

// ClientStats is the object used to collect and send data to the server for processing
type ClientStats struct {
	localAddr  string
	serverAddr string
	stop       chan struct{}
}

// NewClient create a new client object for use
func NewClient(clientConfig *ClientConfig, serverConfig *ServerConfig) (*ClientStats, error) {
	return &ClientStats{
		localAddr:  clientConfig.Domain + ":" + strconv.Itoa(clientConfig.Port),
		serverAddr: serverConfig.Domain + ":" + strconv.Itoa(serverConfig.Port),
		stop:       make(chan struct{}),
	}, nil
}

// Run starts sending the profiling stats to the server
// NOTE: the server must be running prior to starting
func (c *ClientStats) Run() {

	var localAddr *net.UDPAddr
	var serverAddr *net.UDPAddr
	var client *net.UDPConn
	var err error

	serverAddr, err = net.ResolveUDPAddr(udp, c.serverAddr)
	if err != nil {
		panic(err)
	}

	localAddr, err = net.ResolveUDPAddr(udp, c.localAddr)
	if err != nil {
		panic(err)
	}

	client, err = net.DialUDP(udp, localAddr, serverAddr)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	stats := new(Stats)
	stats.MemStats = new(runtime.MemStats)
	encoder := gob.NewEncoder(client)
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			runtime.ReadMemStats(stats.MemStats)
			err = encoder.Encode(stats)
			if err != nil {
				fmt.Println(stats, err)
			}

		case <-c.stop:
			fmt.Println("done")
			return
		}
	}
}

// Stop halts the client from sending any more data to the server,
// but may be run again at any time.
func (c *ClientStats) Stop() {
	c.stop <- struct{}{}
}
