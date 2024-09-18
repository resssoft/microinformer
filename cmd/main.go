package main

import (
	manager "microinformer/internal/maanger"
	"microinformer/internal/settings"
	"microinformer/internal/webserver"
)

//TODO: dynamic html and css - rewrite by api
//TODO: ps list monitor
//TODO: set custom info to api from others machines (bot statuses, cpu, home disks, pinger by the others servers)
//TODO: show images (graphics)
//TODO: parser
// DRm credit, zkh
//TODO: item types - ping, graphic (save old values)
//TODO: system info  - left block, ping - central, graphics - bottom (current pc, others pc)
//TODO: buttons panel - social events and top process history (other page)

func main() {
	settingsService := settings.NewService()
	managerService := manager.NewService(settingsService)
	managerService.Configure()
	webServer := webserver.NewService(settingsService, managerService)
	webServer.Start()
}
