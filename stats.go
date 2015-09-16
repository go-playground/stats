package stats

import (
	"os"
	"runtime"
)

const (
	udp        = "udp"
	bufferSize = 1024 * 10
)

// HostInfo contains all of the host specific information such as
// os, platform etc....
type HostInfo struct {
	Hostname          string `json:"hostname"`
	OS                string `json:"os"`
	LogicalProcessors int    `json:"procs"`
}

// Stats contains all of the statistics to be passed and Encoded/Decoded on the Client and Server sides
type Stats struct {
	HostInfo *HostInfo        `json:"host"`
	MemStats runtime.MemStats `json:"mem"`
}

// GetHostInfo return host system information
func GetHostInfo() *HostInfo {

	hostname, err := os.Hostname()
	if err != nil {
		hostname = err.Error()
	}

	return &HostInfo{
		Hostname:          hostname,
		OS:                runtime.GOOS,
		LogicalProcessors: runtime.NumCPU(),
	}
}
