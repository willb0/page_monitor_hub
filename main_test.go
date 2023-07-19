package main

import (
	"net/http"
	"net/http/httptest"
	"page_monitor_hub/models"
	"page_monitor_hub/pkg/hub"
	"page_monitor_hub/router"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllMonitors1(t *testing.T) {
	pageHub := hub.NewPageMonitorHub()
	models.ConnectTestDatabase()
	router := router.SetupRouter(pageHub)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET","/monitors/all",nil)
	router.ServeHTTP(rec,req)

	assert.Equal(t,200,rec.Code)
	assert.Equal(t,0,len(pageHub.GetMonitors()))
}
