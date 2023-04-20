package main

import (
	"page_monitor_hub/pkg/hub"
	"page_monitor_hub/router"
)



func main() {
	pageHub := hub.NewPageMonitorHub()
	r := router.SetupRouter(pageHub)
	r.Run(":3001")
}


