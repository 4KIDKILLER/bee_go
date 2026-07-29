package model

type BeeFile struct {
	Id         int    `json:"id"`
	ParentId   int    `json:"parentId" db:"parent_id"`
	FileId     int    `json:"fileId" db:"file_id"`
	UserId     int    `json:"userId" db:"user_id"`
	FileName   int    `json:"fileName" db:"file_name"`
	FileSize   int    `json:"fileSize" db:"file_size"`
	FilePath   int    `json:"filePath" db:"file_path"`
	FileType   int    `json:"fileType" db:"file_type"`
	UpdateTime string `json:"updateTime" db:"update_time"`
	CreateTime string `json:"createTime" db:"create_time"`
	Type       int    `json:"type"`
	Status     int    `json:"status"`
}
