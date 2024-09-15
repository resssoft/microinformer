package main

import (
	manager "microinformer/internal/maanger"
	"microinformer/internal/settings"
	"microinformer/internal/webserver"
)

//TODO: Save to disk simplify (by start and changes by api)
//TODO: step to step load by setTImeout
//TODO: dynamic change commands list
//TODO: dynamic setTImeout time from settings and changed by ip
//TODO: more often update for app settings
//TODO: save to disk items list (json file)
//TODO: one time display info show
//TODO: ps list monitor
//TODO: set custom info to api from others machines (bot statuses, cpu, home disks, pinger by the others servers)
//TODO: show images (grafics)
//TODO: parser
// DRm credit, zkh

func main() {
	settingsService := settings.NewService()
	managerService := manager.NewService(settingsService)
	managerService.Configure()
	webServer := webserver.NewService(settingsService, managerService)
	webServer.Start()
}
