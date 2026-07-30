package vo

import "goserver/pkg/utils"

type FileListVo struct {
	Id           string        `json:"id"`
	ParentId     string        `json:"parentId"`
	UserId       int           `json:"userId"`
	Name         string        `json:"name"`
	OriginalName string        `json:"originalName"`
	Size         float64       `json:"size"`
	Path         string        `json:"path"`
	Type         int           `json:"type"`
	Tags         []string      `json:"tags"`
	Covers       [3]string     `json:"covers"`
	Remark       string        `json:"remark"`
	CreateTime   utils.BeeTime `json:"createTime"`
	UpdateTime   utils.BeeTime `json:"updateTime"`
}
