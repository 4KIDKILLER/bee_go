package model

type BeeUser struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
	Avatar     string `json:"avatar"`
	Status     int    `json:"status"`
}
