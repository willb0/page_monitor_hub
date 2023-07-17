package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"page_monitor_hub/models"
	"page_monitor_hub/pkg/hub"

	"github.com/gin-gonic/gin"
)

func SetupRouter(pageHub *hub.PageMonitorHub,db *sql.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	rows, err := db.Query("SELECT * from monitors")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next(){
		var monitor_id int
		var url string
		var redis_channel string
		var refresh_rate int

		err = rows.Scan(&monitor_id,&url,&redis_channel,&refresh_rate)
		if err != nil {
			log.Fatal(err)
		}
		pgr := hub.NewPageMonitorRequest(url,redis_channel,refresh_rate)
		pageHub.AddMonitor(pgr)
	}

	r.POST("/add_page_monitor", func(context *gin.Context) {
		StartMonitorRoute(context, pageHub,db)
	})
	r.GET("/get_all_monitors", func(context *gin.Context) {
		AllMonitorsRoute(context, pageHub)
	})
	r.POST("/stop_page_monitor", func(context *gin.Context) {
		StopMonitorRoute(context, pageHub,db)
	})
	r.GET("/stop_all_monitors", func(context *gin.Context) {
		StopAllMonitorsRoute(context, pageHub)
	})
	return r
}

func StopMonitorRoute(context *gin.Context, pageHub *hub.PageMonitorHub, db *sql.DB) bool {
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
	sql := fmt.Sprintf("delete from monitors where url = '%s' and redis_channel='%s' and refresh_rate=%d",pgj.Url,pgj.RedisChannel,pgj.RefreshRate)
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
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

func StartMonitorRoute(context *gin.Context, pageHub *hub.PageMonitorHub, db *sql.DB) bool {
	pgj := models.PageRequestJson{}
	if err := context.ShouldBindJSON(&pgj); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return true
	}
	fmt.Println(pgj)
	pgr := hub.NewPageMonitorRequest(pgj.Url, pgj.RedisChannel, pgj.RefreshRate)
	pageHub.AddMonitor(pgr)
	sql := fmt.Sprintf("insert into monitors (url,redis_channel,refresh_rate) values ('%s','%s',%d)",pgj.Url,pgj.RedisChannel,pgj.RefreshRate)
	_, err := db.Exec(sql)
	if err != nil {
		print("mf sql string")
		log.Fatal(err)
	}
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

