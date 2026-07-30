package dto

type CreateFolderReq struct {
	ParentId   string `json:"parentId"`
	FolderName string `json:"folderName"`
}
