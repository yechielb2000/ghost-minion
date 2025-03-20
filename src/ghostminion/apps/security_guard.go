package apps

import (
	"fmt"
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

func (c *SecurityGuardApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Starting SecurityGuard app.")
	time.Sleep(2 * time.Hour)
	c.mu.Lock()
	c.isSafe = false
	c.mu.Unlock()
}

func (c *SecurityGuardApp) Stop() error {
	fmt.Println("Stopping SecurityGuard app.")
	return nil
}

func (c *SecurityGuardApp) Validate() error {
	c.isSafe = true
	return nil
}
