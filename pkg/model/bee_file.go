package model

import "time"

type BeeFile struct {
	Id         int       `json:"id"`
	ParentId   string    `json:"parentId" db:"parent_id"`
	FileId     string    `json:"fileId" db:"file_id"`
	UserId     int       `json:"userId" db:"user_id"`
	FileName   string    `json:"fileName" db:"file_name"`
	FileSize   float64   `json:"fileSize" db:"file_size"`
	FilePath   string    `json:"filePath" db:"file_path"`
	FileType   int       `json:"fileType" db:"file_type"`
	FileExt    string    `json:"fileExt" db:"file_ext"`
	UpdateTime time.Time `json:"updateTime" db:"update_time"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
	Type       int       `json:"type"`
	Status     int       `json:"status"`
}
