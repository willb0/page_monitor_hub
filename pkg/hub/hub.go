package hub

import (
	"context"
	"fmt"
	"page_monitor/pkg/pagerefresh"

	"github.com/redis/go-redis/v9"
	"gopkg.in/tomb.v1"
)


type PageRefreshHub struct {
	channels map[string](*PageRefreshRequest)
}
type PageRefreshRequest struct {
	url           string
	testtomb      *tomb.Tomb
	redis_channel string
	refresher     pagerefresh.PageRefresher
}

func (p *PageRefreshRequest) GetUrl() string{
	return p.url
}

func NewPageRefreshRequest(page_url string, redis_channel string, refresh_rate int) *PageRefreshRequest {
	testTomb := tomb.Tomb{}
	ctx, cancel := context.WithCancel(context.Background())
	refresher := pagerefresh.NewPageRefresher(page_url, redis_channel)
	go func() {
		defer testTomb.Done()
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		refresher.WatchForChangesAndNotify(ctx, rdb, refresh_rate, &testTomb)
		println("Finished that mafucka")
		cancel()
	}()
	return &PageRefreshRequest{
		url:           page_url,
		testtomb:      &testTomb,
		redis_channel: redis_channel,
		refresher:     *refresher,
	}
}

func NewPageRefresherHub() *PageRefreshHub {
	return &PageRefreshHub{
		channels: make(map[string](*PageRefreshRequest)),
	}
}

func (p *PageRefreshHub) AddPage(pg *PageRefreshRequest) {
	p.channels[pg.url] = pg
}
func (p *PageRefreshHub) RemovePage(page_url string) {
	println(p.channels[page_url])
	p.channels[page_url].testtomb.Kill(fmt.Errorf("death from above"))
	delete(p.channels, page_url)
}

func (p *PageRefreshHub) GetMonitors() map[string](*PageRefreshRequest){
	return p.channels
}