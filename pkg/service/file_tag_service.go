package service

import (
	"errors"
	"goserver/pkg/dao"
	"goserver/pkg/utils"
	"goserver/pkg/vo"
)

var (
	AddTagErr    = "标签创建失败"
	DeleteTagErr = "标签删除失败"
	UpdateTagErr = "标签修改失败"
	QueryTagErr  = "标签查询失败"
)
var (
	Err6350 = errors.New("6350:" + AddTagErr)
	Err6351 = errors.New("6351:" + DeleteTagErr)
	Err6352 = errors.New("6352:" + DeleteTagErr)
)

// 错误码范围6350-6399
type FileTagService struct {
	fileTagDao *dao.FileTagDao
}

func NewFileTagService(fileTagDao *dao.FileTagDao) (fileTagService *FileTagService) {
	fileTagService = &FileTagService{fileTagDao}
	return
}

func (fileTagService *FileTagService) CreateFileTagService(tagName, fileId string, userId int) (bool, *vo.FileTagVo, error) {

	tagId, _ := utils.GetUUID()

	_, err := fileTagService.fileTagDao.Insert(tagName, tagId, fileId, userId)

	data := new(vo.FileTagVo)

	if err != nil {
		return false, data, Err6350
	}

	data.Id = tagId
	data.TagName = tagName

	return true, data, nil
}

func (fileTagService *FileTagService) DeleteFileTagService(tagId string) (bool, error) {
	result, err := fileTagService.fileTagDao.DeleteRowByTagId(tagId)
	if err != nil {
		return false, Err6351
	}
	row, _ := result.RowsAffected()
	if row == 1 {
		return true, nil
	} else {
		return false, Err6352
	}
}
