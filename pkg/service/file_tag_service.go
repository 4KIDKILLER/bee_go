package service

import (
	"errors"
	"goserver/pkg/dao"
	"goserver/pkg/utils"
)

var (
	AddTagErr    = "标签创建失败"
	DeleteTagErr = "标签删除失败"
	UpdateTagErr = "标签修改失败"
	QueryTagErr  = "标签查询失败"
)
var (
	Err6350 = errors.New("6350:" + AddTagErr)
)

// 错误码范围6350-6399
type FileTagService struct {
	fileTagDao *dao.FileTagDao
}

func NewFileTagService(fileTagDao *dao.FileTagDao) (fileTagService *FileTagService) {
	fileTagService = &FileTagService{fileTagDao}
	return
}

func (fileTagService *FileTagService) CreateFileTagService(tagName, fileId string, userId int) (bool, error) {

	tagId, _ := utils.GetUUID()

	_, err := fileTagService.fileTagDao.Insert(tagName, tagId, fileId, userId)

	if err != nil {
		return false, Err6350
	}

	return true, nil
}
