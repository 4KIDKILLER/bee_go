package dto

type CreateFolderReq struct {
	ParentId   string `json:"parentId"`
	FolderName string `json:"folderName"`
	// Cover1     string `json:"cover1"`
	// Cover2     string `json:"cover2"`
	// Cover3     string `json:"cover3"`
	// Remark     string `json:"remark"`
	// Tags       string `json:"tags"`
}

type GetFileListReq struct {
	ParentId string `json:"parentId"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type DeleteSoftReq struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
}

type RenameReq struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
}

type EditRemarkReq struct {
	Id     string `json:"id"`
	Remark string `json:"remark"`
}

type SetCoverReq struct {
	Id       string `json:"id"`
	Cover    string `json:"cover"`
	Position int    `json:"position"`
}
