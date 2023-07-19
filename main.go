package main

import (
	"page_monitor_hub/models"
	"page_monitor_hub/pkg/hub"
	"page_monitor_hub/router"
)

func main() {
	pageHub := hub.NewPageMonitorHub()
	models.ConnectDatabase()
	r := router.SetupRouter(pageHub)
	println("listening on :3001")
	r.Run(":3001")
}


