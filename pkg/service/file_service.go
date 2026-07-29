package service

import (
	"errors"
	"goserver/pkg/dao"
	"goserver/pkg/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var (
	Err6250 = errors.New("6250: 文件上传失败")
	Err6251 = errors.New("6251: 文件上传失败")
	Err6252 = errors.New("6252: 文件上传失败")
	Err6253 = errors.New("6253: 文件上传失败")
)

// 错误码范围6250-6299
type FileService struct {
	fileDao *dao.FileDao
}

func NewFileService(fileDao *dao.FileDao) (fileService *FileService) {
	fileService = &FileService{fileDao}
	return
}

func (fileService *FileService) CreateFileService(file multipart.File, parentId, fileName string, userId int) (bool, error) {

	uploadDir := "./file"
	if mkdirErr := os.MkdirAll(uploadDir, 0755); mkdirErr != nil {
		return false, Err6250
	}

	//创建文件ID
	fileId, idErr := utils.GetUUID()
	if idErr != nil {
		return false, Err6251
	}

	randomFilename := fileId + filepath.Ext(fileName)
	dstPath := filepath.Join(uploadDir, randomFilename)

	//创建目标文件并保存上传内容
	dst, dstErr := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if dstErr != nil {
		return false, Err6252
	}
	defer dst.Close()

	_, copyErr := io.Copy(dst, file)
	if copyErr != nil {
		return false, Err6253
	}
	return false, nil
	// fileService.fileDao.Insert(parentId, fileId, fileName, dstPath)
}
