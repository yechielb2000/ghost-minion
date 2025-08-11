package apps

import (
	"errors"
	"fmt"
	"ghostminion/config"
	"os"
	"sync"
	"time"
)

type SecurityGuardApp struct {
	isSafe bool
	mu     sync.Mutex
}

func (c *SecurityGuardApp) IsSafeToRun() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isSafe
}

var stopSecurityGuardApp = false

func (c *SecurityGuardApp) Start(wg *sync.WaitGroup) {
	configInstance, _ := config.GetConfig()
	defer wg.Done()
	fmt.Println("Starting SecurityGuard app.")
	time.Sleep(2 * time.Hour)
	c.mu.Lock()
	c.isSafe = false
	c.mu.Unlock()
	allFilesInPlace(configInstance)
	/*
		isSafe should be false on one of these terms
		- it has been too much time without communicating with the C2 (3 days)
		- unknown process or user has touched the db or config file
		- someone wrote the process name in its bash history
		- any of the files that was supposed to be in its place is not anymore
		- the cpu of the target has increase too much because of our process
	*/
}

func (c *SecurityGuardApp) Stop() error {
	stopSecurityGuardApp = true
	return nil
}

func (c *SecurityGuardApp) Validate() error {
	c.isSafe = true
	return nil
}

func allFilesInPlace(configInstance *config.Config) bool {
	if _, err := os.Stat(configInstance.Installation.ConfigFile); errors.Is(err, os.ErrNotExist) {
		return false
	}
	if _, err := os.Stat(configInstance.Installation.DBPath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func DidSomeoneSearchMe(configInstance *config.Config) bool {
	return false
}
