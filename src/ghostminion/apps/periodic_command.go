package apps

import (
	"fmt"
	"sync"
)

type PeriodicCommandApp struct {
	Command string
	Timeout uint // in seconds default is 120
}

func (c *PeriodicCommandApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Starting PeriodicCommand app.")
}

func (c *PeriodicCommandApp) Stop() error {
	fmt.Println("Stopping PeriodicCommand app.")
	return nil
}
