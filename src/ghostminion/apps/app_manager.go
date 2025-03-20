package apps

import (
	"fmt"
	"log"
	"sync"
)

type App interface {
	Start(wg *sync.WaitGroup)
	Stop() error
}

type AppManager struct {
	apps map[string]App
	mu   sync.Mutex
}

func NewAppManager() *AppManager {
	return &AppManager{
		apps: make(map[string]App),
	}
}

func (am *AppManager) AddApp(name string, app App) {
	am.mu.Lock()
	defer am.mu.Unlock()
	if am.apps[name] != nil {
		log.Printf("App %s was already exists. Overwriting app", name)
	}
	am.apps[name] = app
}

func (am *AppManager) RemoveApp(name string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	delete(am.apps, name)
}

func (am *AppManager) StartAll(wg *sync.WaitGroup) {
	am.mu.Lock()
	defer am.mu.Unlock()

	for name, app := range am.apps {
		wg.Add(1)
		go app.Start(wg)
		fmt.Printf("Started app: %s\n", name)
	}
}

func (am *AppManager) StopAll() {
	am.mu.Lock()
	defer am.mu.Unlock()

	for name, app := range am.apps {
		err := app.Stop()
		if err != nil {
			log.Printf("Error stopping app: %s\n", name)
		} else {
			fmt.Printf("Stopped app: %s\n", name)
		}
	}
}
