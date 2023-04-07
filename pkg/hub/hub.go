package hub

import (
	"context"
	"fmt"
	"log"
	"page_monitor_hub/pkg/pagerefresh"

	"github.com/redis/go-redis/v9"
	"gopkg.in/tomb.v1"
)


type PageMonitorHub struct {
	channels map[string](*PageMonitorRequest)
}
type PageMonitorRequest struct {
	url           string
	testtomb      *tomb.Tomb
	redis_channel string
	refresher     pagerefresh.PageMonitor
}

func (p *PageMonitorRequest) GetUrl() string{
	return p.url
}

func NewPageMonitorRequest(page_url string, redis_channel string, refresh_rate int) *PageMonitorRequest {
	testTomb := tomb.Tomb{}
	ctx, cancel := context.WithCancel(context.Background())
	refresher := pagerefresh.NewPageMonitor(page_url, redis_channel)
	go func() {
		defer testTomb.Done()
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		if err := rdb.Ping(ctx); err != nil {
			println("You need to connect to redis lil bro. get a redis running with no pw on localhost:6379")
			log.Fatal(err)
		}
		refresher.WatchForChangesAndNotify(ctx, rdb, refresh_rate, &testTomb)
		println("Finished that mafucka")
		cancel()
	}()
	return &PageMonitorRequest{
		url:           page_url,
		testtomb:      &testTomb,
		redis_channel: redis_channel,
		refresher:     *refresher,
	}
}

func NewPageMonitorHub() *PageMonitorHub {
	return &PageMonitorHub{
		channels: make(map[string](*PageMonitorRequest)),
	}
}

func (p *PageMonitorHub) AddMonitor(pg *PageMonitorRequest) {
	p.channels[pg.url] = pg
}
func (p *PageMonitorHub) RemoveMonitor(page_url string) {
	println(p.channels[page_url])
	p.channels[page_url].testtomb.Kill(fmt.Errorf("death from above"))
	delete(p.channels, page_url)
}

func (p *PageMonitorHub) GetMonitors() map[string](*PageMonitorRequest){
	return p.channels
}