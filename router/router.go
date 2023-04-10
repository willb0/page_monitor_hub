package router

import (
	"fmt"
	"net/http"
	"page_monitor_hub/models"
	"page_monitor_hub/pkg/hub"

	"github.com/gin-gonic/gin"
)

func StopMonitorRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
	pgj := models.PageRequestJson{}
	if err := context.BindJSON(&pgj); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return true
	}
	_,ok := pageHub.GetMonitors()[pgj.Url]
	if !ok {
		context.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"message": "Aborting a monitor that was never started"})
		return true
	}
	fmt.Println(pgj)
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
	context.JSON(http.StatusAccepted, &pgr)
	return false
}

func StopAllMonitorsRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
	for key, _ := range pageHub.GetMonitors() {
		pageHub.RemoveMonitor(key)
	}
	context.Writer.WriteHeader(200)
	return false

}

