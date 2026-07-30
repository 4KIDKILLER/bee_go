package model

import "time"

type BeeUser struct {
	Id         int       `json:"id"`
	UserId     int       `json:"userId" db:"user_id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
	UpdateTime time.Time `json:"updateTime" db:"update_time"`
	Avatar     string    `json:"avatar"`
	Status     int       `json:"status"`
}
