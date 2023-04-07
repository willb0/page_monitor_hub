# page_monitor

i'm writing this repo in go to monitor websites and send notifications of content change using with redis pub sub

you can run page monitor rn if you run a site on localhost:3000 and have docker:

```sh
docker run -p 6379:6379 --name some-redis -d redis
git clone https://github.com/willb0/page_monitor
cd page_monitor
go mod download
go build
./page_monitor -url https://localhost:3000 -refresh_rate 5
```

now you have a go server watching the webpage, and it will publish any updates with the content on the channel "page_refresh" on redis

example connecting

```py
import time,redis
conn = redis.Redis(host='localhost',port=6379,db=1)
sub = conn.pubsub()
sub.subscribe('page_refresh')
while(True):
    time.sleep(1)
    print(sub.get_message())
```
