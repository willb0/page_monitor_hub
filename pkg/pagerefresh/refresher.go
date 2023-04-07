package pagerefresh

import (
	"context"
	"crypto/sha256"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"gopkg.in/tomb.v1"
)

type PageMonitor struct {
	page_url        string
	redis_channel   string
	current_content []byte
}

func NewPageMonitor(page_url string, redis_channel string) *PageMonitor {
	return &PageMonitor{
		page_url:        page_url,
		redis_channel:   redis_channel,
		current_content: []byte(page_url),
	}
}

func (p *PageMonitor) CheckForChanges() bool {
	html := p.GetHTML()
	hasher := sha256.New()
	hasher.Write([]byte(html))
	hashed_html := hasher.Sum(nil)
	res := reflect.DeepEqual(hashed_html, p.current_content)
	if !res {
		p.current_content = hashed_html
		return true
	}
	return false
}

func (p *PageMonitor) GetHTML() []byte {
	resp, err := http.Get(p.page_url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return html
}

func (p *PageMonitor) WatchForChangesAndNotify(ctx context.Context, r *redis.Client, refresh_rate int, t *tomb.Tomb) {
	for {
		select {
		case <-t.Dying():
			return
		default:
			time.Sleep(time.Second * time.Duration(refresh_rate))
			println("checking for changes on: ", p.page_url)
			if p.CheckForChanges() {
				println("it changed!")
				r.Publish(ctx, p.redis_channel, p.GetHTML())
			}
		}
	}
}
