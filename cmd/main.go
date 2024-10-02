package main

import (
	manager "microinformer/internal/maanger"
	"microinformer/internal/settings"
	"microinformer/internal/webserver"
)

//roadmap
//TODO: dynamic html and css - rewrite by api
//TODO: ps list monitor
//TODO: set custom info to api from others machines (bot statuses, cpu, home disks, pinger by the others servers)
//TODO: show images (graphics)
//TODO: parser
// DRm credit, zkh
//TODO: item types - ping, graphic (save old values)
//TODO: system info  - left block, ping - central, graphics - bottom (current pc, others pc)
//TODO: buttons panel - social events and top process history (other page)

//TODO:  commands types - http check status, check https
//TODO: fix duplicate ID - command to b64
//TODO: admin panel for change panel html and css and quick show
//TODO: show part of info from js (date, time)
//TODO: frontend js - if item info not found in the new request - remove from page
//TODO: auth and tokens? check fail auth
//TODO: ping in the backside with custom timeout by item and read for frontend only result + check by old result
//TODO: admin panel - check items list after update by new get result
//TODO: result class - different styles for different result (not only fail)

func main() {
	settingsService := settings.NewService()
	managerService := manager.NewService(settingsService)
	managerService.Configure()
	webServer := webserver.NewService(settingsService, managerService)
	webServer.Start()
}
