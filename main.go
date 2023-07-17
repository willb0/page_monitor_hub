package main

import (
	"database/sql"
	"log"
	"page_monitor_hub/pkg/hub"
	"page_monitor_hub/router"

	_ "github.com/mattn/go-sqlite3"
)



func main() {
	pageHub := hub.NewPageMonitorHub()
	db, err:= sql.Open("sqlite3","page-monitor")
	if err != nil {
		log.Fatal(err)
	}
	r := router.SetupRouter(pageHub,db)
	r.Run(":3001")
}


