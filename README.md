# page_monitor_hub



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
