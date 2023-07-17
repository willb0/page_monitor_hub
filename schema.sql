create table if not exists monitors (
    monitor_id integer primary key,
    url text not null,
    redis_channel text not null,
    refresh_rate integer not null
    );