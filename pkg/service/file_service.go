package service

import (
	"errors"
	"goserver/pkg/dao"
	"goserver/pkg/dto"
	"goserver/pkg/infrastructure/config"
	"goserver/pkg/model"
	"goserver/pkg/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"log"

	"github.com/shopspring/decimal"
)

var (
	UploadErr       = "文件上传失败"
	CreateFolderErr = "文件夹创建失败"
	GetFileListErr  = "获取文件列表失败"
)

var (
	Err6250 = errors.New("6250:" + UploadErr)
	Err6251 = errors.New("6251:" + UploadErr)
	Err6252 = errors.New("6252:" + UploadErr)
	Err6253 = errors.New("6253:" + UploadErr)
	Err6254 = errors.New("6254:" + UploadErr)
	Err6255 = errors.New("6255:" + UploadErr)
	Err6256 = errors.New("6256:" + UploadErr)
	Err6257 = errors.New("6257:" + UploadErr)
	Err6258 = errors.New("6258:" + CreateFolderErr)
	Err6259 = errors.New("6259:" + CreateFolderErr)
	Err6260 = errors.New("6260:" + CreateFolderErr)
	Err6261 = errors.New("6261:" + GetFileListErr)
	Err6262 = errors.New("6262:" + GetFileListErr)
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

func (fileService *FileService) CreateFileService(file multipart.File, parentId, fileName, tags, remark string, fileSize int64, userId int) (bool, error) {

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

	insert, insertErr := fileService.fileDao.Insert(parentId, fileId, fileName, dstPath, tags, "", "", "", remark, userId, 2, size)
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

func (fileService *FileService) CreateFolderService(reqData *dto.CreateFolderReq, userId int) (bool, error) {

	//创建文件夹ID
	folderId, _ := utils.GetUUID()

	size, _ := decimal.NewFromString(".0")
	insert, insertErr := fileService.fileDao.Insert(reqData.ParentId, folderId, reqData.FolderName, reqData.DistPath, reqData.Tags, reqData.Cover1, reqData.Cover2, reqData.Cover3, reqData.Remark, userId, 1, size)

	if insertErr != nil {
		log.Printf("%v: %v", Err6254, insertErr)
		return false, Err6258
	}

	rowsAffected, err := insert.RowsAffected()
	if err != nil {
		return false, Err6259
	}
	if rowsAffected != 1 {
		return false, Err6260
	}
	return true, nil
}

func (fileService *FileService) GetUserFileList(parentId string, userId, page, pageSize int) (int, []*model.BeeFile, error) {
	fileCount, countErr := fileService.fileDao.CountFileByUserIdAndParentId(userId, parentId)
	if countErr != nil {
		log.Printf("%v: %v", Err6261, countErr)
		return 0, nil, Err6261
	}
	fileList, fileErr := fileService.fileDao.QueryFilesByUserId(parentId, userId, (page-1)*pageSize, pageSize)
	if fileErr != nil {
		log.Printf("%v: %v", Err6261, fileErr)
		return 0, nil, Err6262
	}

	return fileCount, fileList, nil
}
