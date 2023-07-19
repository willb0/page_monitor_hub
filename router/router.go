package router

import (
	"net/http"
	"page_monitor_hub/models"
	"page_monitor_hub/pkg/hub"

	"github.com/gin-gonic/gin"
)

func SetupRouter(pageHub *hub.PageMonitorHub) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/monitors/create", func(context *gin.Context) {
		AddMonitorRoute(context, pageHub)
	})
	r.GET("/monitors/all", func(context *gin.Context) {
		AllMonitorsRoute(context, pageHub)
	})
	r.POST("/monitors/delete", func(context *gin.Context) {
		DeleteMonitorRoute(context, pageHub)
	})
	r.GET("/monitors/all/stop", func(context *gin.Context) {
		DeleteAllMonitorsRoute(context, pageHub)
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	  })
	return r
}
func AddMonitorRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
	pgj := models.PageRequestJson{}
	if err := context.ShouldBindJSON(&pgj); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return true
	}
	pageHub.AddMonitor(pgj.Url,pgj.RedisChannel,pgj.RefreshRate)
	context.JSON(201, &pgj)
	return false
}

func DeleteMonitorRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
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
	context.JSON(204,pgj)
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



func DeleteAllMonitorsRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
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

