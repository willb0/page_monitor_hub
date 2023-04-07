package models

type PageRequestJson struct {
	Url          string `json:"url"`
	RedisChannel string `json:"redis_channel"`
	RefreshRate  int    `json:"refresh_rate"`
}