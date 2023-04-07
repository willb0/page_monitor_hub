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
	fmt.Println(pgj)
	pageHub.RemoveMonitor(pgj.Url)
	return false
}

func AllMonitorsRoute(context *gin.Context, pageHub *hub.PageMonitorHub) {
	for key, element := range pageHub.GetMonitors() {
		fmt.Println("k:", key, "=>", "v:", element.GetUrl())
	}
	context.JSON(http.StatusAccepted, &models.PageRequestJson{
		Url:          "a",
		RedisChannel: "b",
	})
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

	return false

}
