package dto

type CreateFolderReq struct {
	ParentId   string `json:"parentId"`
	FolderName string `json:"folderName"`
}

type GetFileListReq struct {
	ParentId string `json:"parentId"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}
