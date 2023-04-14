package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"page_monitor_hub/models"
	"page_monitor_hub/pkg/hub"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAllMonitors1(t *testing.T) {
	pageHub := hub.NewPageMonitorHub()
	r := SetupRouter(pageHub)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET","/get_all_monitors",nil)
	r.ServeHTTP(rec,req)

	assert.Equal(t,200,rec.Code)
	assert.Equal(t,0,len(pageHub.GetMonitors()))
}

func TestAllMonitors2(t *testing.T) {
	pageHub := hub.NewPageMonitorHub()
	r := SetupRouter(pageHub)

	rec := httptest.NewRecorder()
	pgm := hub.NewPageMonitorRequest("http://purdue.edu","purdue",10)
	pageHub.AddMonitor(pgm)
	body, _ := json.Marshal(pgm)
	req, _ := http.NewRequest("GET","/get_all_monitors",bytes.NewReader(body))
	r.ServeHTTP(rec,req)
	assert.Equal(t,200,rec.Code)
	assert.Equal(t,1,len(pageHub.GetMonitors()))
}

func TestNewMonitor1(t *testing.T) {
	pageHub := hub.NewPageMonitorHub()
	r := SetupRouter(pageHub)
	pgjson := models.PageRequestJson{
		Url : "https://www.purdue.edu",
		RedisChannel: "purdue",
		RefreshRate: 10,
	}
	rec := httptest.NewRecorder()
	body, _ := json.Marshal(pgjson)
	req, _ := http.NewRequest("POST","/add_page_monitor",bytes.NewReader(body))
	r.ServeHTTP(rec,req)
	//t.Log(req.Body)
	assert.Equal(t,200,rec.Code)
	assert.Equal(t,1,len(pageHub.GetMonitors()))
}
func TestNewMonitor2(t *testing.T) {
	pageHub := hub.NewPageMonitorHub()
	r := SetupRouter(pageHub)
	pgjson := gin.H{"message": "Killing a monitor that was never started"}
	rec := httptest.NewRecorder()
	body, _ := json.Marshal(pgjson)
	req, _ := http.NewRequest("POST","/add_page_monitor",bytes.NewReader(body))
	r.ServeHTTP(rec,req)
	t.Log("the hell")
	t.Log(pageHub.GetMonitors())
	assert.Equal(t,400,rec.Code)
	assert.Equal(t,0,len(pageHub.GetMonitors()))
}

func TestDeleteMonitor1(t *testing.T) {
	pageHub := hub.NewPageMonitorHub()
	r := SetupRouter(pageHub)
	pgjson := models.PageRequestJson{
		Url : "https://www.purdue.edu",
		RedisChannel: "purdue",
		RefreshRate: 10,
	}
	rec := httptest.NewRecorder()
	body, _ := json.Marshal(pgjson)
	req, _ := http.NewRequest("POST","/add_page_monitor",bytes.NewReader(body))
	r.ServeHTTP(rec,req)
	//t.Log(req.Body)
	assert.Equal(t,200,rec.Code)
	assert.Equal(t,1,len(pageHub.GetMonitors()))

	req, _ = http.NewRequest("POST","/stop_page_monitor",bytes.NewReader(body))
	r.ServeHTTP(rec,req)

	assert.Equal(t,200,rec.Code)
	assert.Equal(t,0,len(pageHub.GetMonitors()))

}

func TestDeleteMonitor2(t *testing.T) {
	pageHub := hub.NewPageMonitorHub()
	r := SetupRouter(pageHub)
	pgjson := models.PageRequestJson{
		Url : "https://www.purdue.edu",
		RedisChannel: "purdue",
		RefreshRate: 10,
	}
	rec := httptest.NewRecorder()
	body, _ := json.Marshal(pgjson)
	req, _ := http.NewRequest("POST","/stop_page_monitor",bytes.NewReader(body))
	r.ServeHTTP(rec,req)
	assert.Equal(t,404,rec.Code)
}
func TestDeleteAllMonitors1(t *testing.T) {
	pageHub := hub.NewPageMonitorHub()
	r := SetupRouter(pageHub)
	pgjson := models.PageRequestJson{
		Url : "https://www.purdue.edu",
		RedisChannel: "purdue",
		RefreshRate: 10,
	}	
	rec := httptest.NewRecorder()
	body, _ := json.Marshal(pgjson)
	req, _ := http.NewRequest("POST","/add_page_monitor",bytes.NewReader(body))
	r.ServeHTTP(rec,req)
	assert.Equal(t,200,rec.Code)
	assert.Equal(t,1,len(pageHub.GetMonitors()))

	req, _ = http.NewRequest("GET","/stop_all_monitors",nil)
	r.ServeHTTP(rec,req)
	//t.Log(req.Body)
	assert.Equal(t,200,rec.Code)
	assert.Equal(t,0,len(pageHub.GetMonitors()))
}

func TestDeleteAllMonitors2(t *testing.T) {
	pageHub := hub.NewPageMonitorHub()
	r := SetupRouter(pageHub)
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("GET","/stop_all_monitors",nil)
	r.ServeHTTP(rec,req)
	//t.Log(req.Body)
	assert.Equal(t,404,rec.Code)
	assert.Equal(t,0,len(pageHub.GetMonitors()))
}