package apps

import (
	"fmt"
	"sync"
)

type PeriodicGetFileApp struct {
	Path     string
	MaxSize  int // maximum allowed size in bytes.
	Interval int
	CheckMD5 bool // check if file was changed since last run.
}

func (c *PeriodicGetFileApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Starting PeriodicGetFile app.")
}

func (c *PeriodicGetFileApp) Stop() error {
	fmt.Println("Stopping PeriodicGetFile app.")
	return nil
}
