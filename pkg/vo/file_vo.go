package vo

import "goserver/pkg/utils"

type FileListVo struct {
	ParentId   string        `json:"parentId"`
	FileId     string        `json:"fileId"`
	UserId     int           `json:"userId"`
	FileName   string        `json:"fileName"`
	FileSize   float64       `json:"fileSize"`
	FilePath   string        `json:"filePath"`
	FileType   int           `json:"fileType"`
	Tags       []string      `json:"tags"`
	Covers     [3]string     `json:"covers"`
	Remark     string        `json:"remark"`
	CreateTime utils.BeeTime `json:"createTime"`
	UpdateTime utils.BeeTime `json:"updateTime"`
}
