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
	/*
		isSafe should be false on these terms
		- it has been too much time without communicating with the C2 (3 days)
		- unknown process or user has touched the db or config file
		- someone wrote the process name in its bash history
		- any of the files that was supposed to be in its place is not anymore
		- the cpu of the target has increase too much because of our process
	*/
	fmt.Println("Stopping SecurityGuard app.")
	return nil
}

func (c *SecurityGuardApp) Validate() error {
	c.isSafe = true
	return nil
}
