package apps

import (
	"fmt"
	"github.com/MarinX/keylogger"
	"sync"
)

type KeyLoggerApp struct{}

func (c *KeyLoggerApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	keyboard := keylogger.FindKeyboardDevice()
	if len(keyboard) <= 0 {
		fmt.Println("No keyboard found...you will need to provide manual input path")
		return
	}
	fmt.Println("Found a keyboard at", keyboard)
	klgr, err := keylogger.New(keyboard)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer klgr.Close()
	for {
		events := klgr.Read()
		for e := range events {
			switch e.Type {
			case keylogger.EvKey:
				if e.KeyPress() {
					fmt.Println("[event] press key ", e.KeyString())
				} else if e.KeyRelease() {
					fmt.Println("[event] release key ", e.KeyString())
				}
				break // save to db ⬆️
			}
		}
	}
}

func (c *KeyLoggerApp) Stop() error {
	fmt.Println("Stopping KeyLogger app.")
	return nil
}

func (c *KeyLoggerApp) Validate() error {
	return nil
}
