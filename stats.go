package stats

import (
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

const (
	udp        = "udp"
	bufferSize = 1024 * 10
)

// CPUPercentages contains the CPU percentage information
type CPUPercentages struct {
	CPU       string
	User      float64
	System    float64
	Idle      float64
	Nice      float64
	IOWait    float64
	IRQ       float64
	SoftIRQ   float64
	Steal     float64
	Guest     float64
	GuestNice float64
	Stolen    float64
	Total     float64
}

// GoMemory contains go specific memory metrics
type GoMemory struct {
	NumGC               uint32 `json:"numgc"`
	LastGC              uint64 `json:"lastgc"`
	LastGCPauseDuration uint64 `json:"lastgcpause"`
	Alloc               uint64 `json:"alloc"`
	HeapAlloc           uint64 `json:"heap"`
	HeapSys             uint64 `json:"sys"`
	lastNumGC           uint32
	mem                 *runtime.MemStats
}

// GoInfo contains go specific metrics and stats
type GoInfo struct {
	Version    string   `json:"gover"`
	Memory     GoMemory `json:"gomem"`
	GoRoutines int      `json:"goroutines"`
}

// MemInfo contains memory info including swap information
type MemInfo struct {
	Memory *mem.VirtualMemoryStat `json:"mem,omitempty"`
	Swap   *mem.SwapMemoryStat    `json:"swap,omitempty"`
}

// CPUInfo contains CPU information
type CPUInfo struct {
	CPU            []cpu.CPUInfoStat  `json:"cpu,omitempty"`
	PerCPUTimes    []cpu.CPUTimesStat `json:"percputimes,omitempty"`
	TotalTimes     []cpu.CPUTimesStat `json:"totaltimes,omitempty"`
	PrevCPUTimes   []cpu.CPUTimesStat `json:"prevpercputimes,omitempty"`
	PrevTotalTimes []cpu.CPUTimesStat `json:"prevtotaltimes,omitempty"`
}

// Stats contains all of the statistics to be passed and Encoded/Decoded on the Client and Server sides
type Stats struct {
	HostInfo     *host.HostInfoStat `json:"hostInfo,omitempty"`
	CPUInfo      *CPUInfo           `json:"cpu,omitempty"`
	MemInfo      *MemInfo           `json:"memInfo,omitempty"`
	GoInfo       *GoInfo            `json:"goInfo,omitempty"`
	HTTPRequests []*HTTPRequest     `json:"http"`
}

func (s *Stats) initGoInfo() {
	s.GoInfo = new(GoInfo)
	s.GoInfo.Memory.mem = new(runtime.MemStats)
}

// GetHostInfo populates Stats with host system information
func (s *Stats) GetHostInfo() {

	if s.GoInfo == nil {
		s.initGoInfo()
	}

	info, _ := host.HostInfo()

	s.HostInfo = info
	s.GoInfo.Version = runtime.Version()
}

// GetCPUInfo populates Stats with hosts CPU information
func (s *Stats) GetCPUInfo() {

	if s.CPUInfo == nil {
		s.CPUInfo = new(CPUInfo)
	}

	s.CPUInfo.CPU, _ = cpu.CPUInfo()
}

// GetCPUTimes populates Stats with hosts CPU timing information
func (s *Stats) GetCPUTimes() {

	if s.CPUInfo == nil {
		s.CPUInfo = new(CPUInfo)
	}

	s.CPUInfo.PrevCPUTimes = s.CPUInfo.PerCPUTimes
	s.CPUInfo.PerCPUTimes, _ = cpu.CPUTimes(true)

	if len(s.CPUInfo.PrevCPUTimes) == 0 {
		s.CPUInfo.PrevCPUTimes = s.CPUInfo.PerCPUTimes
	}
}

// CalculateCPUTimes calculates the total CPU times percentages per core
func (s *Stats) CalculateCPUTimes() []CPUPercentages {

	percentages := make([]CPUPercentages, len(s.CPUInfo.PerCPUTimes))

	if len(s.CPUInfo.PrevCPUTimes) == 0 || len(s.CPUInfo.PerCPUTimes) == 0 {
		return percentages
	}

	var diff float64
	var total float64
	var prevTotal float64
	var prevStat cpu.CPUTimesStat
	var cpuStat *CPUPercentages

	for i, t := range s.CPUInfo.PerCPUTimes {
		cpuStat = &percentages[i]
		prevStat = s.CPUInfo.PrevCPUTimes[i]

		total = t.User + t.System + t.Idle + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Guest + t.GuestNice + t.Stolen
		prevTotal = prevStat.User + prevStat.System + prevStat.Idle + prevStat.Nice + prevStat.Iowait + prevStat.Irq + prevStat.Softirq + prevStat.Steal + prevStat.Guest + prevStat.GuestNice + prevStat.Stolen

		diff = total - prevTotal

		cpuStat.CPU = t.CPU
		cpuStat.User = (t.User - prevStat.User) / diff * 100
		cpuStat.System = (t.System - prevStat.System) / diff * 100
		cpuStat.Idle = (t.Idle - prevStat.Idle) / diff * 100
		cpuStat.Nice = (t.Nice - prevStat.Nice) / diff * 100
		cpuStat.IOWait = (t.Iowait - prevStat.Iowait) / diff * 100
		cpuStat.IRQ = (t.Irq - prevStat.Irq) / diff * 100
		cpuStat.SoftIRQ = (t.Softirq - prevStat.Softirq) / diff * 100
		cpuStat.Steal = (t.Steal - prevStat.Steal) / diff * 100
		cpuStat.Guest = (t.Guest - prevStat.Guest) / diff * 100
		cpuStat.GuestNice = (t.GuestNice - prevStat.GuestNice) / diff * 100
		cpuStat.Stolen = (t.Stolen - prevStat.Stolen) / diff * 100
		cpuStat.Total = 100 * (diff - (t.Idle - prevStat.Idle)) / diff
	}

	return percentages
}

// GetAllCPUInfo populates Stats with hosts CPU information and Timings
func (s *Stats) GetAllCPUInfo() {
	s.GetCPUInfo()
	s.GetCPUTimes()
}

// GetTotalCPUTimes populates Stats with hosts CPU timing information
func (s *Stats) GetTotalCPUTimes() {

	if s.CPUInfo == nil {
		s.CPUInfo = new(CPUInfo)
	}

	s.CPUInfo.PrevTotalTimes = s.CPUInfo.TotalTimes
	s.CPUInfo.TotalTimes, _ = cpu.CPUTimes(false)

	if len(s.CPUInfo.PrevTotalTimes) == 0 {
		s.CPUInfo.PrevTotalTimes = s.CPUInfo.TotalTimes
	}
}

// CalculateTotalCPUTimes calculates the total CPU times percentages
func (s *Stats) CalculateTotalCPUTimes() []CPUPercentages {

	percentages := make([]CPUPercentages, len(s.CPUInfo.TotalTimes))

	if len(s.CPUInfo.PrevTotalTimes) == 0 || len(s.CPUInfo.TotalTimes) == 0 {
		return percentages
	}

	var diff float64
	var total float64
	var prevTotal float64
	var prevStat cpu.CPUTimesStat
	var cpuStat *CPUPercentages

	for i, t := range s.CPUInfo.TotalTimes {
		cpuStat = &percentages[i]
		prevStat = s.CPUInfo.PrevTotalTimes[i]

		total = t.User + t.System + t.Idle + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Guest + t.GuestNice + t.Stolen
		prevTotal = prevStat.User + prevStat.System + prevStat.Idle + prevStat.Nice + prevStat.Iowait + prevStat.Irq + prevStat.Softirq + prevStat.Steal + prevStat.Guest + prevStat.GuestNice + prevStat.Stolen

		diff = total - prevTotal

		cpuStat.CPU = t.CPU
		cpuStat.User = (t.User - prevStat.User) / diff * 100
		cpuStat.System = (t.System - prevStat.System) / diff * 100
		cpuStat.Idle = (t.Idle - prevStat.Idle) / diff * 100
		cpuStat.Nice = (t.Nice - prevStat.Nice) / diff * 100
		cpuStat.IOWait = (t.Iowait - prevStat.Iowait) / diff * 100
		cpuStat.IRQ = (t.Irq - prevStat.Irq) / diff * 100
		cpuStat.SoftIRQ = (t.Softirq - prevStat.Softirq) / diff * 100
		cpuStat.Steal = (t.Steal - prevStat.Steal) / diff * 100
		cpuStat.Guest = (t.Guest - prevStat.Guest) / diff * 100
		cpuStat.GuestNice = (t.GuestNice - prevStat.GuestNice) / diff * 100
		cpuStat.Stolen = (t.Stolen - prevStat.Stolen) / diff * 100
		cpuStat.Total = 100 * (diff - (t.Idle - prevStat.Idle)) / diff
	}

	return percentages
}

// GetMemoryInfo populates Stats with host and go process memory information
func (s *Stats) GetMemoryInfo(logMemory, logGoMemory bool) {

	if logGoMemory {
		if s.GoInfo == nil {
			s.initGoInfo()
		}

		runtime.ReadMemStats(s.GoInfo.Memory.mem)
		s.GoInfo.GoRoutines = runtime.NumGoroutine()
		s.GoInfo.Memory.Alloc = s.GoInfo.Memory.mem.Alloc
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
	}

	if logMemory {

		if s.MemInfo == nil {
			s.MemInfo = new(MemInfo)
		}

		s.MemInfo.Memory, _ = mem.VirtualMemory()
		s.MemInfo.Swap, _ = mem.SwapMemory()
	}
}
