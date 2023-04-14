package router

import (
	"fmt"
	"net/http"
	"page_monitor_hub/models"
	"page_monitor_hub/pkg/hub"

	"github.com/gin-gonic/gin"
)

func SetupRouter(pageHub *hub.PageMonitorHub) *gin.Engine {
	r := gin.Default()
	r.POST("/add_page_monitor", func(context *gin.Context) {
		StartMonitorRoute(context, pageHub)
	})
	r.GET("/get_all_monitors", func(context *gin.Context) {
		AllMonitorsRoute(context, pageHub)
	})
	r.POST("/stop_page_monitor", func(context *gin.Context) {
		StopMonitorRoute(context, pageHub)
	})
	r.GET("/stop_all_monitors", func(context *gin.Context) {
		StopAllMonitorsRoute(context, pageHub)
	})
	return r
}

func StopMonitorRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
	pgj := models.PageRequestJson{}
	if err := context.BindJSON(&pgj); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return true
	}
	_,ok := pageHub.GetMonitors()[pgj.Url]
	if !ok {
		context.AbortWithStatusJSON(404,gin.H{"message": "Killing a monitor that was never started"})
		return true
	}
	pageHub.RemoveMonitor(pgj.Url)
	context.JSON(200,pgj)
	return false
}

func AllMonitorsRoute(context *gin.Context, pageHub *hub.PageMonitorHub) {
	keys := make([]string, len(pageHub.GetMonitors()))
	i := 0;
	for key := range pageHub.GetMonitors() {
		keys[i] = key
		i++
	}
	context.JSON(200,keys)
}

func StartMonitorRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
	pgj := models.PageRequestJson{}
	if err := context.BindJSON(&pgj); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return true
	}
	fmt.Println(pgj)
	pgr := hub.NewPageMonitorRequest(pgj.Url, pgj.RedisChannel, pgj.RefreshRate)
	pageHub.AddMonitor(pgr)
	context.JSON(200, &pgr)
	return false
}

func StopAllMonitorsRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
	if len(pageHub.GetMonitors()) == 0{
		context.AbortWithStatusJSON(404,gin.H{"message": "No monitors are currently running"})
		return true
	}
	for key, _ := range pageHub.GetMonitors() {
		pageHub.RemoveMonitor(key)
	}
	context.Writer.WriteHeader(200)
	return false

}

