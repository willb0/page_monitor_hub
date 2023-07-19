package router

import (
	"log"
	"net/http"
	"page_monitor_hub/models"
	"page_monitor_hub/pkg/hub"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func SetupRouter(pageHub *hub.PageMonitorHub) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())
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
	println(models.DB == nil)
	return r
}
func AddMonitorRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
	pgj := models.PageRequestJson{}
	if err := context.ShouldBindJSON(&pgj); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return true
	}
	pageHub.AddMonitor(pgj.Url,pgj.RedisChannel,pgj.RefreshRate)
	newMonitor := models.Monitor{Url:pgj.Url,RedisChannel: pgj.RedisChannel, RefreshRate: pgj.RefreshRate}
	models.DB.Create(&newMonitor)
	//models.DB.Commit()
	context.JSON(201, &pgj)
	return false
}

func DeleteMonitorRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
	pgj := models.PageRequestJson{}
	if err := context.BindJSON(&pgj); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return true
	}
	monitor := models.Monitor{}
	res := models.DB.Where("url = ?",pgj.Url).Delete(&monitor)
	if res.Error != nil{
		context.AbortWithError(500,res.Error)
		return true
	}
	if res.RowsAffected < 1 {
		context.AbortWithError(404,res.Error)
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
	var monitors []models.Monitor
	if err := models.DB.Find(&monitors).Error; err != nil {
		log.Print(err)
		context.AbortWithStatusJSON(404,err)
	}
	context.JSON(200,monitors)
}



func DeleteAllMonitorsRoute(context *gin.Context, pageHub *hub.PageMonitorHub) bool {
	db := models.DB.Select("DELETE FROM MONITORS;")
	if db.Error!= nil{
		context.AbortWithError(500,db.Error)
		return true
	}
	models.DB.Commit()
	context.Writer.WriteHeader(200)
	return false

}

