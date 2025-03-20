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
	c.isSafe = true
	time.Sleep(10 * time.Second)
	c.isSafe = false
	fmt.Println("Starting SecurityGuard app.")
}

func (c *SecurityGuardApp) Stop() error {
	fmt.Println("Stopping SecurityGuard app.")
	return nil
}

func (c *SecurityGuardApp) Validate() error {
	return nil
}
