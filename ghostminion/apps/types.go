package apps

import (
	"encoding/json"
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
	Stop() error
	Validate() error
}
