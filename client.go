package stats

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type httpStats struct {
	lock     sync.RWMutex
	requests []*HTTPRequest
}

// Add adds an entry to the httpStats array
func (h *httpStats) Add(r *HTTPRequest) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.requests = append(h.requests, r)
}

func (h *httpStats) extract() []*HTTPRequest {
	h.lock.Lock()
	defer h.lock.Unlock()

	old := h.requests
	h.requests = []*HTTPRequest{}
	return old
}

// ClientConfig is used to initialize a new ClientStats object
type ClientConfig struct {
	Domain           string
	Port             int
	PollInterval     int
	Debug            bool
	LogHostInfo      bool
	LogCPUInfo       bool
	LogTotalCPUTimes bool
	LogPerCPUTimes   bool
	LogMemory        bool
	LogGoMemory      bool
	CustomBufferSize int
}

// ClientStats is the object used to collect and send data to the server for processing
type ClientStats struct {
	localAddr        string
	serverAddr       string
	stop             chan struct{}
	interval         int
	debug            bool
	httpStats        *httpStats
	logHostInfo      bool
	logCPUInfo       bool
	logTotalCPUTimes bool
	logPerCPUTimes   bool
	logMemory        bool
	logGoMemory      bool
	bufferSize       int
}

// NewClient create a new client object for use
func NewClient(clientConfig *ClientConfig, serverConfig *ServerConfig) (*ClientStats, error) {

	bSize := clientConfig.CustomBufferSize
	if bSize == 0 {
		bSize = defaultBufferSize
	}

	return &ClientStats{
		localAddr:        clientConfig.Domain + ":" + strconv.Itoa(clientConfig.Port),
		serverAddr:       serverConfig.Domain + ":" + strconv.Itoa(serverConfig.Port),
		interval:         clientConfig.PollInterval,
		stop:             make(chan struct{}),
		debug:            clientConfig.Debug,
		httpStats:        new(httpStats),
		logHostInfo:      clientConfig.LogHostInfo,
		logCPUInfo:       clientConfig.LogCPUInfo,
		logTotalCPUTimes: clientConfig.LogTotalCPUTimes,
		logPerCPUTimes:   clientConfig.LogPerCPUTimes,
		logMemory:        clientConfig.LogMemory,
		logGoMemory:      clientConfig.LogGoMemory,
		bufferSize:       bSize,
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

	client.SetWriteBuffer(c.bufferSize)

	stats := new(Stats)

	if c.logHostInfo {
		stats.GetHostInfo()
	}

	if c.logCPUInfo {
		stats.GetCPUInfo()
	}

	var bytesWritten int
	var bytes []byte
	ticker := time.NewTicker(time.Millisecond * time.Duration(c.interval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			if c.logTotalCPUTimes {
				stats.GetTotalCPUTimes()
			}

			if c.logPerCPUTimes {
				stats.GetCPUTimes()
			}

			stats.GetMemoryInfo(c.logMemory, c.logGoMemory)
			stats.HTTPRequests = c.httpStats.extract()

			bytes, err = json.Marshal(stats)
			bytesWritten, err = client.Write(bytes)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if c.debug {
				fmt.Println("Wrote:", bytesWritten, "bytes")
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
