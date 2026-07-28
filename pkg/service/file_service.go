package service

import "goserver/pkg/dao"

type FileService struct {
	userDao *dao.FileDao
}

func NewFileService(fileDao *dao.FileDao) (fileService *FileService) {
	fileService = &FileService{fileDao}
	return
}
