package apps

import (
	"fmt"
	"sync"
)

type KeyLoggerApp struct{}

func (c *KeyLoggerApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Starting KeyLogger app.")
}

func (c *KeyLoggerApp) Stop() error {
	fmt.Println("Stopping KeyLogger app.")
	return nil
}

func (c *KeyLoggerApp) Validate() error {
	return nil
}
