package main

import (
	"fmt"
	"net/http"
	"page_monitor/pkg/hub"

	"github.com/gin-gonic/gin"
)





type PageRequestJson struct {
	Url          string `json:"url"`
	RedisChannel string `json:"redis_channel"`
	RefreshRate  int    `json:"refresh_rate"`
}

func main() {
	pageHub := hub.NewPageRefresherHub()
	r := gin.Default()
	r.POST("/add_page_monitor", func(context *gin.Context) {
		pgj := PageRequestJson{}
		if err := context.BindJSON(&pgj); err != nil {
			context.AbortWithError(http.StatusBadRequest, err)
			return
		}
		fmt.Println(pgj)
		pgr := hub.NewPageRefreshRequest(pgj.Url, pgj.RedisChannel, pgj.RefreshRate)
		pageHub.AddPage(pgr)
		context.JSON(http.StatusAccepted, &pgr)
	})
	r.GET("/get_all_monitors", func(context *gin.Context) {
		for key, element := range pageHub.GetMonitors() {
			fmt.Println("k:", key, "=>", "v:", element.GetUrl())
		}
		context.JSON(http.StatusAccepted, &PageRequestJson{
			Url:          "a",
			RedisChannel: "b",
		})
	})
	r.POST("/stop_page_monitor", func(context *gin.Context) {
		pgj := PageRequestJson{}
		if err := context.BindJSON(&pgj); err != nil {
			context.AbortWithError(http.StatusBadRequest, err)
			return
		}
		fmt.Println(pgj)
		pageHub.RemovePage(pgj.Url)
	})
	r.Run(":3000")
}

func NewPageRefresherHub() {
	panic("unimplemented")
}
