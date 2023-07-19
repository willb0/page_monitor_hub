create table if not exists monitors (
    id integer primary key,
    url text not null,
    redis_channel text not null,
    refresh_rate integer not null
    );