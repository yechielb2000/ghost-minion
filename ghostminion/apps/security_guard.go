package apps

import (
	"errors"
	"ghostminion/monitors"
	"os"
	"sync"
	"time"
)

type SecurityGuardApp struct {
	*BaseApp
	isSafe         bool
	mu             sync.Mutex
	FilesExistence []string `json:"FilesExistence"`
}

func (ctx *SecurityGuardApp) setIsSafe(isSafe bool) {
	ctx.mu.Lock()
	ctx.isSafe = isSafe
	ctx.mu.Unlock()
}

func (ctx *SecurityGuardApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	lgr.Info("Starting SecurityGuard App")
	ctx.checkFileExistence(ctx.FilesExistence)

	CPUMonitor := monitors.NewCPUMonitor(os.Getpid(), 50.0, 1*time.Second)
	CPUMonitor.Start()
	for isCPUSafe := range CPUMonitor.IsCPUSafe() {
		if !isCPUSafe {
			lgr.Warn("CPU is too high")
			ctx.setIsSafe(isCPUSafe)
		}
	}

	CPUMonitor.Stop()
	/*
		isSafe should be false on one of these terms
		- it has been too much time without communicating with the C2 (3 days)
		- unknown process or user has touched the db or config file
		- someone wrote the process name in its bash history
		- any of the files that was supposed to be in its place is not anymore
	*/
}

func (ctx *SecurityGuardApp) Stop() {
}

func (ctx *SecurityGuardApp) validateParams() error {
	return nil
}

func (ctx *SecurityGuardApp) checkFileExistence(files []string) {
	for _, file := range files {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			lgr.Warn("File does not exist: ", file)
			ctx.setIsSafe(false)
			break
		}
	}
}
