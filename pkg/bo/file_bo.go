package bo

import "goserver/pkg/overwrite"

type FileListBo struct {
	Id               int               `json:"id"`
	FileSize         float64           `json:"fileSize"`
	FileExt          string            `json:"fileExt"`
	FileOriginalName string            `json:"fileOriginalName"`
	FilePath         string            `json:"filePath"`
	FileThumbPath    string            `json:"fileThumbPath"`
	FileType         int               `json:"fileType"`
	Tags             string            `json:"tags"`
	Cover1           string            `json:"cover1"`
	Cover2           string            `json:"cover2"`
	Cover3           string            `json:"cover3"`
	Remark           string            `json:"remark"`
	CreateTime       overwrite.BeeTime `json:"createTime"`
	UpdateTime       overwrite.BeeTime `json:"updateTime"`
}
