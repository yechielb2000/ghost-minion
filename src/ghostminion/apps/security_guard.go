package apps

import (
	"fmt"
	"sync"
)

type SecurityGuardApp struct{}

func (c *SecurityGuardApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Starting SecurityGuard app.")
}

func (c *SecurityGuardApp) Stop() error {
	fmt.Println("Stopping SecurityGuard app.")
	return nil
}
