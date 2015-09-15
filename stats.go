package stats

import "runtime"

const (
	udp        = "udp"
	bufferSize = 1024 * 10
)

// Stats contains all of the statistics to be passed and Encoded/Decoded on the Client and Server sides
type Stats struct {
	MemStats runtime.MemStats
}
