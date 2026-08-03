package model

import (
	"goserver/pkg/utils"
)

type BeeFile struct {
	Id               int           `json:"id"`
	ParentId         string        `json:"parentId" db:"parent_id"`
	FileId           string        `json:"fileId" db:"file_id"`
	UserId           int           `json:"userId" db:"user_id"`
	FileSize         float64       `json:"fileSize" db:"file_size"`
	FileExt          string        `json:"fileExt" db:"file_ext"`
	FileOriginalName string        `json:"fileOriginalName" db:"file_original_name"`
	FilePath         string        `json:"filePath" db:"file_path"`
	FileThumbPath    string        `json:"fileThumbPath" db:"file_thumb_path"`
	FileType         int           `json:"fileType" db:"file_type"`
	Tags             string        `json:"tags"`
	Cover1           string        `json:"cover1" db:"cover_1"`
	Cover2           string        `json:"cover2" db:"cover_2"`
	Cover3           string        `json:"cover3" db:"cover_3"`
	Remark           string        `json:"remark"`
	CreateTime       utils.BeeTime `json:"createTime" db:"create_time"`
	UpdateTime       utils.BeeTime `json:"updateTime" db:"update_time"`
	Status           int           `json:"status"`
}
