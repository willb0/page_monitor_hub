# page_monitor_hub

### context
this is my work in progress on a REST API to manage a slew of different page monitors

I initially started this to set up text notifications when my butchers website updated

After learning some go, and redis pub/sub, I set up a page monitor module which will check for updates in the hashed HTML and publish a message over redis to all subscribers.

this architecture can be used to build notification clients for any end goal, not just SMS

### this project

I want to build a dashboard that lets me create, manage, and delete monitors while also maybe having some analytics. first part of this is:

backend: CRUD interface for page monitor objects, refactor current code

next will be frontend/deployment tweaks


### to run
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
    }
}
```
