package apps

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type AppManager struct {
	apps map[string]App
	mu   sync.Mutex
	wg   sync.WaitGroup
}

var (
	appManagerInstance *AppManager
	once               sync.Once
)

func GetAppManagerInstance() *AppManager {
	once.Do(func() {
		appManagerInstance = &AppManager{
			apps: make(map[string]App),
		}
	})
	return appManagerInstance
}

func (am *AppManager) GetApp(name string) (App, error) {
	am.mu.Lock()
	defer am.mu.Unlock()
	app := am.apps[name]
	if app == nil {
		return nil, errors.New("am: app not found")
	}
	return app, nil
}

func (am *AppManager) ListApps() map[string]App {
	am.mu.Lock()
	defer am.mu.Unlock()
	return am.apps
}

func (am *AppManager) StartApp(name string, app App) {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.apps[name] != nil {
		log.Printf("App \"%s\" already exists. Overwriting app", name)
	}

	am.wg.Add(1)
	go app.Start(&am.wg)
	am.apps[name] = app
}

func (am *AppManager) StopApp(name string) error {
	am.mu.Lock()
	defer am.mu.Unlock()
	app := am.apps[name]
	if err := app.Stop(); err != nil {
		return err
	}
	am.wg.Done()
	delete(am.apps, name)
	return nil
}

func (am *AppManager) StartAll() {
	for name, app := range am.apps {
		am.StartApp(name, app)
	}
}

func (am *AppManager) StopAll() {
	am.mu.Lock()
	defer am.mu.Unlock()

	for name, _ := range am.apps {
		err := am.StopApp(name)
		if err != nil {
			log.Printf("Error stopping app: %s\n", name)
		} else {
			fmt.Printf("Stopped app: %s\n", name)
		}
	}
}
