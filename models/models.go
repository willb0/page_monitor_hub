package models

type Monitor struct {
	ID     uint   `json:"id" gorm:"primary_key";autoIncrement:true`
	Url  string `json:"url"`
	RefreshRate int `json:"refresh_rate"`
	RedisChannel string `json:"redis_channel"`
}

type PageRequestJson struct {
	Url          string `json:"url" binding:"required"`
	RedisChannel string `json:"redis_channel" binding:"required"`
	RefreshRate  int    `json:"refresh_rate" binding:"required"`

}