package monitors

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CPUMonitor struct {
	pid       int
	threshold float64
	interval  time.Duration
	isSafe    bool
	stop      chan struct{}
	isCPUSafe chan bool
}

func NewCPUMonitor(pid int, threshold float64, interval time.Duration) *CPUMonitor {
	return &CPUMonitor{
		pid:       pid,
		threshold: threshold,
		interval:  interval,
		isSafe:    true,
		stop:      make(chan struct{}),
		isCPUSafe: make(chan bool, 1),
	}
}

func (m *CPUMonitor) Start() {
	go func() {
		lastState := m.isSafe
		for {
			select {
			case <-m.stop:
				return
			default:
				usage, err := getCPUUsage(m.pid, m.interval)
				if err == nil {
					current := usage <= m.threshold
					m.isSafe = current
					if current != lastState {
						m.isCPUSafe <- current
						lastState = current
					}
				}
			}
		}
	}()
}

func (m *CPUMonitor) Stop() {
	close(m.stop)
}

func (m *CPUMonitor) IsCPUSafe() <-chan bool {
	return m.isCPUSafe
}

func readProcessCPU(pid int) (uint64, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	utime, _ := strconv.ParseUint(fields[13], 10, 64)
	stime, _ := strconv.ParseUint(fields[14], 10, 64)
	return utime + stime, nil
}

func readSystemCPU() (uint64, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(strings.Split(string(data), "\n")[0]) // "cpu ..."
	var total uint64
	for _, v := range fields[1:] {
		val, _ := strconv.ParseUint(v, 10, 64)
		total += val
	}
	return total, nil
}

func getCPUUsage(pid int, interval time.Duration) (float64, error) {
	p1, err := readProcessCPU(pid)
	if err != nil {
		return 0, err
	}
	s1, err := readSystemCPU()
	if err != nil {
		return 0, err
	}
	time.Sleep(interval)
	p2, err := readProcessCPU(pid)
	if err != nil {
		return 0, err
	}
	s2, err := readSystemCPU()
	if err != nil {
		return 0, err
	}
	procDelta := float64(p2 - p1)
	sysDelta := float64(s2 - s1)
	if sysDelta == 0 {
		return 0, nil
	}
	return (procDelta / sysDelta) * 100.0, nil
}
