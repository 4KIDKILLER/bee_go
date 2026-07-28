package model

type BeeUser struct {
	Id         int    `json:"id"`
	UserId     int    `json:"userId" db:"user_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CreateTime string `json:"createTime" db:"create_time"`
	UpdateTime string `json:"updateTime" db:"update_time"`
	Avatar     string `json:"avatar"`
	Status     int    `json:"status"`
}
