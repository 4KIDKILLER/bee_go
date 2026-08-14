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
)

var (
	UploadErr       = "文件上传失败"
	CreateFolderErr = "文件夹创建失败"
	GetFileListErr  = "获取文件列表失败"
	CreateThumbErr  = "预览图创建失败"
	RemoveFileErr   = "删除失败"
	FileRenameErr   = "文件名称修改失败"
	GetFolderErr    = "获取文件夹列表失败"
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
	Err6263 = errors.New("6263:" + CreateFolderErr)
	Err6264 = errors.New("6264:" + CreateThumbErr)
	Err6265 = errors.New("6265:" + CreateThumbErr)
	Err6266 = errors.New("6266:" + RemoveFileErr)
	Err6267 = errors.New("6267:" + FileRenameErr)
	Err6268 = errors.New("6268:" + GetFolderErr)
)

// 错误码范围6250-6299
type FileService struct {
	fileDao    *dao.FileDao
	fileConfig config.FileConfig
}

type FileTreeNode struct {
	ParentId   string          `json:"parentId"`
	FileId     string          `json:"fileId"`
	UserId     int             `json:"userId"`
	FolderName string          `json:"folderName"`
	Children   []*FileTreeNode `json:"children"`
}

func NewFileService(fileDao *dao.FileDao, fileConfig config.FileConfig) (fileService *FileService) {
	fileService = &FileService{
		fileDao,
		fileConfig,
	}
	return
}

func (fileService *FileService) UploadFileService(file multipart.File, parentId, fileOriginalName, tags, remark string, fileSize int64, userId int) (bool, error) {

	uploadDir := utils.GetUploadDir(fileService.fileConfig)

	oriAbsPath := filepath.Join(fileService.fileConfig.Path, uploadDir.Original)
	if mkdirErr := os.MkdirAll(oriAbsPath, 0755); mkdirErr != nil {
		return false, Err6250
	}

	//创建文件ID
	fileId, _ := utils.GetUUID()

	fileExt := filepath.Ext(fileOriginalName)
	randomFilename := fileId + fileExt
	dstPath := filepath.Join(oriAbsPath, randomFilename)

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

	insert, insertErr := fileService.fileDao.Insert(parentId, fileId, fileOriginalName, fileExt, uploadDir.Original, "", tags, "", "", "", remark, fileSize, userId, 2)
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

	var thumbDstPath = filepath.Join(fileService.fileConfig.Path, uploadDir.Thumb, randomFilename)

	//创建预览图生成协程
	go func() {
		compErr := utils.ImageCompression(dstPath, thumbDstPath, fileExt)
		if compErr != nil {
			log.Printf("%v: %v", Err6264, compErr)
		}
		_, thumbErr := fileService.fileDao.UpdateThumbPathByFileId(uploadDir.Thumb, fileId, userId)
		if compErr != nil {
			log.Printf("%v: %v", Err6265, thumbErr)
		}
	}()
	return true, nil
}

func (fileService *FileService) CreateFolderService(reqData *dto.CreateFolderReq, userId int) (bool, error) {

	//创建文件夹ID
	folderId, _ := utils.GetUUID()

	insert, insertErr := fileService.fileDao.Insert(reqData.ParentId, folderId, reqData.FolderName, "", "", "", "", "", "", "", "", 0, userId, 1)

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

func (fileService *FileService) GetUserFileListService(parentId string, userId, page, pageSize int) (int, []*model.BeeFile, error) {
	fileCount, countErr := fileService.fileDao.CountFileByParentId(userId, parentId)
	if countErr != nil {
		log.Printf("%v: %v", Err6261, countErr)
		return 0, nil, Err6261
	}
	fileList, fileErr := fileService.fileDao.QueryUserFiles(parentId, userId, (page-1)*pageSize, pageSize)
	if fileErr != nil {
		log.Printf("%v: %v", Err6261, fileErr)
		return 0, nil, Err6262
	}

	return fileCount, fileList, nil
}

func (fileService *FileService) RemoveFileService(fileId string, userId, fileType int) (bool, error) {
	_, err := fileService.fileDao.UpdateStatusByFileIdAndFileType(fileId, userId, fileType, 2)
	if err != nil {
		log.Printf("%v: %v", Err6266, err)
		return false, Err6266
	}

	return true, nil
}

func (fileService *FileService) UpdateOriginalNameService(name, fileId string, userId int) (bool, error) {
	_, err := fileService.fileDao.UploadOriginalNameByFileId(name, fileId, userId)

	if err != nil {
		log.Printf("%v: %v", Err6267, err)
		return false, Err6267
	}

	return true, nil
}

func (fileService *FileService) GetUserFileTreeService(userId int) ([]*FileTreeNode, error) {
	folderList, folderErr := fileService.fileDao.QueryUserFolders(userId)
	if folderErr != nil {
		log.Printf("%v: %v", Err6268, folderErr)
		return nil, Err6268
	}

	tree := make([]*FileTreeNode, 0, len(folderList))
	if len(folderList) == 0 {
		return tree, nil
	}

	nodeMap := make(map[string]*FileTreeNode, len(folderList))
	for _, folder := range folderList {
		nodeMap[folder.FileId] = &FileTreeNode{
			ParentId:   folder.ParentId,
			FileId:     folder.FileId,
			UserId:     folder.UserId,
			FolderName: folder.FileOriginalName,
			Children:   make([]*FileTreeNode, 0),
		}
	}

	for _, folder := range folderList {
		//按顺序从映射表中获取文件夹节点
		node := nodeMap[folder.FileId]
		//获取当前文件夹父级节点
		parent, hasParent := nodeMap[folder.ParentId]

		//如果当节点不存在父级节点，则为顶级节点，直接添加到tree列表
		if folder.ParentId == "" || folder.ParentId == folder.FileId || !hasParent {
			tree = append(tree, node)
			continue
		}
		/*
			否则设置到父级文件夹的children字段中,此处有一个非常总要的知识点由于tree中存放的是
			节点的指针，nodeMap[folder.ParentId]和nodeMap[folder.FileId]获取的是文件指
			针，所以能很方便的执行append操作，因为每次获取的都是一个地址只向的那个对象。指针直
			接抹平了层级访问带来的对象获取问题。
		*/
		parent.Children = append(parent.Children, node)
	}

	return tree, nil
}
