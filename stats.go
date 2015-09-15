package stats

import "runtime"

const (
	udp = "udp"
)

// Stats contains all of the statistics to be passed and Encoded/Decoded on the Client and Server sides
type Stats struct {
	MemStats *runtime.MemStats
}
