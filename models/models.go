package models

type PageRequestJson struct {
	Url          string `json:"url" binding:"required"`
	RedisChannel string `json:"redis_channel" binding:"required"`
	RefreshRate  int    `json:"refresh_rate" binding:"required"`

}