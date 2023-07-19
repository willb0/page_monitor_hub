package hub

import (
	"context"
	"fmt"
	"page_monitor_hub/pkg/pagerefresh"

	"github.com/redis/go-redis/v9"
	"gopkg.in/tomb.v1"
)


type PageMonitorHub struct {
	channels map[string](tomb.Tomb)
}


func NewPageMonitorHub() *PageMonitorHub {
	return &PageMonitorHub{
		channels: make(map[string](tomb.Tomb)),
	}
}

func (p *PageMonitorHub) AddMonitor(url, redis_channel string, refresh_rate int) {
	testTomb := tomb.Tomb{}
	ctx, cancel := context.WithCancel(context.Background())
	refresher := pagerefresh.NewPageMonitor(url, redis_channel)
	go func() {
		defer testTomb.Done()
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			panic("wheeeee no redis you loser")
		}
		refresher.WatchForChangesAndNotify(ctx, rdb, refresh_rate, &testTomb)
		println("Finished that mafucka")
		cancel()
	}()
	p.channels[url] = testTomb
}
func (p *PageMonitorHub) RemoveMonitor(page_url string) {
	tomb := p.channels[page_url]
	(&tomb).Kill(fmt.Errorf("death from above"))
	delete(p.channels, page_url)
}

func (p *PageMonitorHub) GetMonitors() map[string](tomb.Tomb){
	return p.channels
}