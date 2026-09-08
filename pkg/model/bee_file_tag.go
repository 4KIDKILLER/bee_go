package model

import "goserver/pkg/overwrite"

type BeeFileTag struct {
	Id         int               `json:"id"`
	TagId      string            `json:"tagId" db:"tag_id"`
	FileId     string            `json:"fileId" db:"file_id"`
	UserId     int               `json:"userId" db:"user_id"`
	TagName    string            `json:"tagName" db:"tag_name"`
	CreateTime overwrite.BeeTime `json:"createTime" db:"create_time"`
	UpdateTime overwrite.BeeTime `json:"updateTime" db:"update_time"`
}
