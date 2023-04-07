package main

import (
	"page_monitor_hub/pkg/hub"
	"page_monitor_hub/router"

	"github.com/gin-gonic/gin"
)



func main() {
	pageHub := hub.NewPageMonitorHub()
	r := gin.Default()
	r.POST("/add_page_monitor", func(context *gin.Context) {
		router.StartMonitorRoute(context, pageHub)
	})
	r.GET("/get_all_monitors", func(context *gin.Context) {
		router.AllMonitorsRoute(context, pageHub)
	})
	r.POST("/stop_page_monitor", func(context *gin.Context) {
		router.StopMonitorRoute(context, pageHub)
	})
	r.Run(":3000")
}


