package apps

import (
	"fmt"
	"sync"
)

type ClipboardApp struct{}

func (c *ClipboardApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Starting Clipboard app.")
}

func (c *ClipboardApp) Stop() error {
	fmt.Println("Stopping Clipboard app.")
	return nil
}

func (c *ClipboardApp) Validate() error {
	return nil
}
