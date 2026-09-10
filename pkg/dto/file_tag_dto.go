package dto

type CreateTargetReq struct {
	TagName string `json:"tagName"`
	FileId  string `json:"fileId"`
}

type DeleteTargetReq struct {
	Id string `json:"id"`
}
