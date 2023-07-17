# page_monitor_hub

### context
this is my work in progress on a REST API to manage a slew of different page monitors

I initially started this to set up text notifications when my butchers website updated

After learning some go, and redis pub/sub, I set up a page monitor module which will check for updates in the hashed HTML and publish a message over redis to all subscribers.

this architecture can be used to build notification clients for any end goal, not just SMS

### this project

I want to build a dashboard that lets me create, manage, and delete monitors while also maybe having some analytics. first part of this is:

backend: CRUD interface for page monitor objects, refactor current code

### to run
i would highly reccomend getting docker, as you need to run redis for this.
you'll also need sqlite dbs called 'page-monitor' and 'page-monitor-test', I included two empty ones with the schema in the repo
the schema is in schema.sql, run `./sqllite.sh` to set the db's up
#### redis
once you have the sqlite created (remember the program expects a db called 'page-monitor' and table called monitors), run the redis
`docker run -d -p 6379:6379 --name=redis-pgmonitor redis`

#### go
```sh

go build -o bin
./bin/page_monitor_hub

```

endpoints:
```json
{
    /* POST */
    "/add_page_monitor" : {
        "url":"https://purdue.edu",
        "redis_channel":"purdue_monitor",
        "refresh_rate":5,
    },
    /* GET */
    "/get_all_monitors" : {
    },
    /* POST */ 
    "/stop_page_monitor": {
        "url":"https://purdue.edu",
        "redis_channel":"purdue_monitor",
        "refresh_rate":5,
    },
    /* GET */
    "/stop_all_monitors" : {
    }
}
```

FUTURE WORK:
- move the redis client to async for lower memory footprint
- add a "last changed" hashmap maybe make a websocket api? 
- change to a sqlite DB for storing the monitors etc
- think about moving out the orchestrator to its own repo and making a SDK 
- could even make some structured types where you have a fucntion that passes (redis_channel,middleware func,html_comparator) for event logging etc 
