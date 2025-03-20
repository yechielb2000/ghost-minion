package apps

import (
	"fmt"
	"sync"
)

type KeyLoggerApp struct{}

func (c *KeyLoggerApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	//maybe use library https://github.com/MarinX/keylogger
}

func (c *KeyLoggerApp) Stop() error {
	fmt.Println("Stopping KeyLogger app.")
	return nil
}

func (c *KeyLoggerApp) Validate() error {
	return nil
}
