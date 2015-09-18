package stats

import (
	"runtime"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

const (
	udp        = "udp"
	bufferSize = 1024 * 10
)

// GoMemory contains go specific memory metrics
type GoMemory struct {
	NumGC               uint32           `json:"numgc"`
	LastGC              uint64           `json:"lastgc"`
	LastGCPauseDuration uint64           `json:"lastgcpause"`
	Allocated           uint64           `json:"alloc"`
	HeapAlloc           uint64           `json:"heap"`
	HeapSys             uint64           `json:"sys"`
	lastNumGC           uint32           `json:"-"`
	mem                 runtime.MemStats `json:"-"`
}

// GoInfo contains go specific metrics and stats
type GoInfo struct {
	Version    string   `json:"gover"`
	Memory     GoMemory `json:"gomem"`
	GoRoutines int      `json:"goroutines"`
}

// MemInfo contains memory info including swap information
type MemInfo struct {
	Memory *mem.VirtualMemoryStat `json:"mem"`
	Swap   *mem.SwapMemoryStat    `json:"swap"`
	// GoMemstats *runtime.MemStats      `json:"gomem"`
}

// Stats contains all of the statistics to be passed and Encoded/Decoded on the Client and Server sides
type Stats struct {
	HostInfo *host.HostInfoStat `json:"hostInfo"`
	MemInfo  *MemInfo           `json:"memInfo"`
	GoInfo   *GoInfo            `json:"goInfo"`
}

// GetHostInfo populates Stats with host system information
func (s *Stats) GetHostInfo() {

	if s.GoInfo == nil {
		s.GoInfo = new(GoInfo)
	}

	info, _ := host.HostInfo()

	s.HostInfo = info
	s.GoInfo.Version = runtime.Version()
}

// GetMemoryInfo populates Stats with host and go process memory information
func (s *Stats) GetMemoryInfo() {

	if s.GoInfo == nil {
		s.GetHostInfo()
	}

	if s.MemInfo == nil {
		s.MemInfo = new(MemInfo)
	}

	runtime.ReadMemStats(&s.GoInfo.Memory.mem)
	s.GoInfo.GoRoutines = runtime.NumGoroutine()
	s.GoInfo.Memory.NumGC = s.GoInfo.Memory.mem.NumGC
	s.GoInfo.Memory.Allocated = s.GoInfo.Memory.mem.Alloc
	s.GoInfo.Memory.HeapAlloc = s.GoInfo.Memory.mem.HeapAlloc
	s.GoInfo.Memory.HeapSys = s.GoInfo.Memory.mem.HeapSys

	if s.GoInfo.Memory.LastGC != s.GoInfo.Memory.mem.LastGC {
		s.GoInfo.Memory.LastGC = s.GoInfo.Memory.mem.LastGC
		s.GoInfo.Memory.NumGC = s.GoInfo.Memory.mem.NumGC - s.GoInfo.Memory.lastNumGC
		s.GoInfo.Memory.lastNumGC = s.GoInfo.Memory.mem.NumGC
		s.GoInfo.Memory.LastGCPauseDuration = s.GoInfo.Memory.mem.PauseNs[(s.GoInfo.Memory.mem.NumGC+255)%256]
	} else {
		s.GoInfo.Memory.NumGC = 0
		s.GoInfo.Memory.LastGCPauseDuration = 0
	}

	s.MemInfo.Memory, _ = mem.VirtualMemory()
	s.MemInfo.Swap, _ = mem.SwapMemory()
}
