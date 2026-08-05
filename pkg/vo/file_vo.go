package vo

import "goserver/pkg/utils"

type FileListVo struct {
	Id           string        `json:"id"`
	ParentId     string        `json:"parentId"`
	UserId       int           `json:"userId"`
	Name         string        `json:"name"`
	OriginalName string        `json:"originalName"`
	Size         float64       `json:"size"`
	Type         int           `json:"type"`
	Tags         []string      `json:"tags"`
	Src          string        `json:"src"`
	ThumbSrc     string        `json:"thumbSrc"`
	Covers       [3]string     `json:"covers"`
	Remark       string        `json:"remark"`
	CreateTime   utils.BeeTime `json:"createTime"`
	UpdateTime   utils.BeeTime `json:"updateTime"`
}
