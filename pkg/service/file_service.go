package service

import (
	"errors"
	"goserver/pkg/dao"
	"goserver/pkg/infrastructure/config"
	"goserver/pkg/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"log"

	"github.com/shopspring/decimal"
)

var (
	Err6250 = errors.New("6250: 文件上传失败")
	Err6251 = errors.New("6251: 文件上传失败")
	Err6252 = errors.New("6252: 文件上传失败")
	Err6253 = errors.New("6253: 文件上传失败")
	Err6254 = errors.New("6254: 文件上传失败")
	Err6255 = errors.New("6255: 文件上传失败")
	Err6256 = errors.New("6256: 文件上传失败")
	Err6257 = errors.New("6257: 文件上传失败")
	Err6258 = errors.New("6258: 文件夹创建失败")
	Err6259 = errors.New("6259: 文件夹创建失败")
	Err6260 = errors.New("6260: 文件夹创建失败")
)

// 错误码范围6250-6299
type FileService struct {
	fileDao    *dao.FileDao
	fileConfig config.FileConfig
}

func NewFileService(fileDao *dao.FileDao, fileConfig config.FileConfig) (fileService *FileService) {
	fileService = &FileService{
		fileDao,
		fileConfig,
	}
	return
}

func (fileService *FileService) CreateFileService(file multipart.File, parentId, fileName string, fileSize int64, userId int) (bool, error) {

	if mkdirErr := os.MkdirAll(fileService.fileConfig.Path, 0755); mkdirErr != nil {
		return false, Err6250
	}

	//创建文件ID
	fileId, _ := utils.GetUUID()

	randomFilename := fileId + filepath.Ext(fileName)
	dstPath := filepath.Join(fileService.fileConfig.Path, randomFilename)

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

	size := decimal.NewFromInt(fileSize).DivRound(decimal.NewFromInt(1024), 2)

	insert, insertErr := fileService.fileDao.Insert(parentId, fileId, fileName, dstPath, userId, 2, size)
	if insertErr != nil {
		log.Printf("%v: %v", Err6254, insertErr)
		return false, Err6254
	}

	rowsAffected, err := insert.RowsAffected()
	if err != nil {
		return false, Err6255
	}
	if rowsAffected != 1 {
		return false, Err6256
	}
	return true, nil
}

func (fileService *FileService) CreateFolderService(parentId, folderName string, userId int) (bool, error) {

	//创建文件夹ID
	folderId, _ := utils.GetUUID()

	size, _ := decimal.NewFromString(".0")
	insert, insertErr := fileService.fileDao.Insert(parentId, folderId, folderName, "", userId, 1, size)

	if insertErr != nil {
		log.Printf("%v: %v", Err6254, insertErr)
		return false, Err6257
	}

	rowsAffected, err := insert.RowsAffected()
	if err != nil {
		return false, Err6258
	}
	if rowsAffected != 1 {
		return false, Err6259
	}
	return true, nil
}
