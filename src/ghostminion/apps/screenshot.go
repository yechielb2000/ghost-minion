package apps

import (
	"fmt"
	"sync"
)

type ScreenshotApp struct {
	Interval uint // in seconds
}

func (c *ScreenshotApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Starting Screenshot app.")
}

func (c *ScreenshotApp) Stop() error {
	fmt.Println("Stopping Screenshot app.")
	return nil
}

func (c *ScreenshotApp) Validate() error {
	return nil
}
