package service

import (
	"errors"
	"fmt"
	"goserver/pkg/dao"
	"goserver/pkg/dto"
	"goserver/pkg/infrastructure/config"
	"goserver/pkg/model"
	"goserver/pkg/utils"
	"goserver/pkg/vo"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"log"
)

var (
	UploadErr           = "文件上传失败"
	CreateFolderErr     = "文件夹创建失败"
	GetFileListErr      = "获取文件列表失败"
	CreateThumbErr      = "预览图创建失败"
	RemoveFileErr       = "文件删除失败"
	FileRenameErr       = "文件名称修改失败"
	GetFolderErr        = "获取文件夹列表失败"
	GetTgsErr           = "获取文件标签失败"
	UpdateRemarkErr     = "修改文件备注失败"
	UpdateCoverErr      = "设置文件夹封面失败"
	PositionValiFailErr = "非法的position参数"
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
	Err6269 = errors.New("6269:" + RemoveFileErr)
	Err6270 = errors.New("6270:" + RemoveFileErr)
	Err6271 = errors.New("6271:" + RemoveFileErr)
	Err6272 = errors.New("6272:" + RemoveFileErr)
	Err6273 = errors.New("6273:" + RemoveFileErr)
	Err6274 = errors.New("6274:" + RemoveFileErr)
	Err6275 = errors.New("6275:" + RemoveFileErr)
	Err6276 = errors.New("6276:" + RemoveFileErr)
	Err6277 = errors.New("6277:" + RemoveFileErr)
	Err6278 = errors.New("6278:" + GetTgsErr)
	Err6279 = errors.New("6279:" + UpdateRemarkErr)
	Err6280 = errors.New("6280:" + UpdateRemarkErr)
	Err6281 = errors.New("6281:" + PositionValiFailErr)
	Err6282 = errors.New("6282:" + UpdateCoverErr)
	Err6283 = errors.New("6283:" + UpdateCoverErr)
)

const (
	//缩略图创建协程最大限制
	thumbnailWorkerCount = 12
	//任务队列最大限制
	thumbnailQueueSize = 300
)

// 错误码范围6250-6299
type FileService struct {
	fileDao        *dao.FileDao
	fileTagDao     *dao.FileTagDao
	fileConfig     config.FileConfig
	thumbnailTasks chan thumbnailTask
}

type thumbnailTask struct {
	srcPath   string
	dstPath   string
	fileExt   string
	thumbPath string
	fileId    string
	userId    int
}

type FileTreeNode struct {
	ParentId   string          `json:"parentId"`
	FileId     string          `json:"fileId"`
	UserId     int             `json:"userId"`
	FolderName string          `json:"folderName"`
	Children   []*FileTreeNode `json:"children"`
}

func NewFileService(fileDao *dao.FileDao, fileTagDao *dao.FileTagDao, fileConfig config.FileConfig) (fileService *FileService) {
	fileService = &FileService{
		fileDao:        fileDao,
		fileConfig:     fileConfig,
		fileTagDao:     fileTagDao,
		thumbnailTasks: make(chan thumbnailTask, thumbnailQueueSize),
	}
	// golang中chan是并发安全的，不会出现资源竞争，所以不需要考虑加锁
	/*
		开启10个协程消费管道数据
	*/
	for range thumbnailWorkerCount {
		go fileService.thumbnailWorker()
	}
	return
}

func (fileService *FileService) thumbnailWorker() {
	for task := range fileService.thumbnailTasks {
		compErr := utils.ImageCompression(task.srcPath, task.dstPath, task.fileExt)
		if compErr != nil {
			log.Printf("%v: %v", Err6264, compErr)
			continue
		}

		_, thumbErr := fileService.fileDao.UpdateRowThumbPathByFileId(task.thumbPath, task.fileId, task.userId)
		if thumbErr != nil {
			log.Printf("%v: %v", Err6265, thumbErr)
		}
	}
}

func (fileService *FileService) UploadFileService(file multipart.File, parentId, fileOriginalName, remark string, fileSize int64, userId int) (bool, error) {

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

	_, copyErr := io.Copy(dst, file)
	if copyErr != nil {
		_ = dst.Close()
		return false, Err6253
	}
	// 入队前关闭文件，确保缩略图Worker读取到完整内容。
	if closeErr := dst.Close(); closeErr != nil {
		return false, Err6253
	}

	insert, insertErr := fileService.fileDao.Insert(parentId, fileId, fileOriginalName, fileExt, uploadDir.Original, "", "", "", "", remark, fileSize, userId, 2)
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

	thumbDstPath := filepath.Join(fileService.fileConfig.Path, uploadDir.Thumb, randomFilename)

	// 将任务交给固定的缩略图Worker；队列中的300个等待槽都被占用时在这里等待。
	fileService.thumbnailTasks <- thumbnailTask{
		srcPath:   dstPath,
		dstPath:   thumbDstPath,
		fileExt:   fileExt,
		thumbPath: uploadDir.Thumb,
		fileId:    fileId,
		userId:    userId,
	}
	return true, nil
}

func (fileService *FileService) CreateFolderService(reqData *dto.CreateFolderReq, userId int) (bool, error) {

	//创建文件夹ID
	folderId, _ := utils.GetUUID()

	insert, insertErr := fileService.fileDao.Insert(reqData.ParentId, folderId, reqData.FolderName, "", "", "", "", "", "", "", 0, userId, 1)

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

func (fileService *FileService) GetUserFileListService(uploadHost, parentId string, userId, page, pageSize int) (*utils.PaginationJson[vo.FileListVo], error) {
	fileCount, countErr := fileService.fileDao.CountRowByParentId(userId, parentId)
	if countErr != nil {
		log.Printf("%v: %v", Err6261, countErr)
		return nil, Err6261
	}
	fileList, err := fileService.fileDao.SelectRowsLimitByUserId(parentId, userId, (page-1)*pageSize, pageSize)
	if err != nil {
		log.Printf("%v: %v", Err6262, err)
		return nil, Err6262
	}

	if len(fileList) == 0 {
		return &utils.PaginationJson[vo.FileListVo]{
			List:     make([]vo.FileListVo, 0),
			Total:    0,
			Page:     page,
			PageSize: pageSize,
		}, nil
	}

	fileIds := make([]string, 0, len(fileList))
	for _, item := range fileList {
		fileIds = append(fileIds, item.FileId)
	}

	fileTags, err := fileService.fileTagDao.SelectTagsByFileIds(fileIds, userId)
	if err != nil {
		log.Printf("%v: %v", Err6278, err)
		return nil, Err6278
	}

	tagMap := make(map[string][]vo.FileTagVo, len(fileIds))
	for _, tag := range fileTags {
		tagMap[tag.FileId] = append(tagMap[tag.FileId], vo.FileTagVo{
			Id:      tag.TagId,
			TagName: tag.TagName,
		})
	}

	dataList := make([]vo.FileListVo, 0, len(fileList))

	for _, item := range fileList {
		covers := [3]string{item.Cover1, item.Cover2, item.Cover3}
		tags := tagMap[item.FileId]
		if tags == nil {
			tags = make([]vo.FileTagVo, 0)
		}
		name := item.FileId + item.FileExt
		src := ""
		thumbSrc := ""
		if item.FileType == 2 {
			src = fmt.Sprintf("%s/%s/%s", uploadHost, item.FilePath, name)
			thumbSrc = fmt.Sprintf("%s/%s/%s", uploadHost, item.FileThumbPath, name)
		}
		dataList = append(dataList, vo.FileListVo{
			ParentId:     item.ParentId,
			Id:           item.FileId,
			UserId:       item.UserId,
			Name:         name,
			OriginalName: item.FileOriginalName,
			Size:         item.FileSize,
			Type:         item.FileType,
			Tags:         tags,
			Src:          src,
			ThumbSrc:     thumbSrc,
			Covers:       covers,
			Remark:       item.Remark,
			CreateTime:   item.CreateTime,
			UpdateTime:   item.UpdateTime,
		})
	}

	resultData := &utils.PaginationJson[vo.FileListVo]{
		List:     dataList,
		Total:    fileCount,
		Page:     page,
		PageSize: pageSize,
	}

	return resultData, nil
}

func (fileService *FileService) DeleteFileSoftService(fileId string, userId, fileType int) (bool, error) {
	_, err := fileService.fileDao.UpdateRowStatusByFileIdAndFileType(fileId, userId, fileType, 2)
	if err != nil {
		log.Printf("%v: %v", Err6266, err)
		return false, Err6266
	}

	return true, nil
}

func (fileService *FileService) UpdateOriginalNameService(name, fileId string, userId int) (bool, error) {
	_, err := fileService.fileDao.UploadRowOriginalNameByFileId(name, fileId, userId)

	if err != nil {
		log.Printf("%v: %v", Err6267, err)
		return false, Err6267
	}

	return true, nil
}

func (fileService *FileService) GetUserFileTreeService(userId int) ([]*FileTreeNode, error) {
	folderList, folderErr := fileService.fileDao.SelectRowsByUserId(userId)
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

func (fileService *FileService) DeleteFolderSoftService(fileId string, userId int) (bool, error) {
	fileIds, fileIdsErr := fileService.fileDao.SelectRowsRecursionByFileId(fileId, userId)
	if fileIdsErr != nil {
		log.Printf("%v: %v", Err6269, fileIdsErr)
		return false, Err6269
	}
	_, statusErr := fileService.fileDao.UpdateRowsStatusByFileIdInIds(fileIds, 2, userId)
	if statusErr != nil {
		log.Printf("%v: %v", Err6270, statusErr)
		return false, Err6270
	}

	return true, nil
}

func (fileService *FileService) DeleteFolderHardService(userId int) (bool, error) {
	resultList, err := fileService.fileDao.SelectRowsByStatusOrderByFileType(2, userId)
	if err != nil {
		log.Printf("%v: %v", Err6271, err)
		return false, Err6271
	}
	var folderIds []string
	var fileList []*model.BeeFile

	for _, item := range resultList {
		if item.FileType == 1 {
			folderIds = append(folderIds, item.FileId)
		} else {
			fileList = append(fileList, item)
		}
	}
	//删除文件夹并更新文件夹状态
	if len(folderIds) > 0 {
		_, err := fileService.fileDao.UpdateRowsStatusByFileIdInIds(folderIds, 3, userId)
		if err != nil {
			log.Printf("%v: %v", Err6272, err)
		}
	}

	if len(fileList) > 0 {
		//删除文件并更新文件状态
		for _, item := range fileList {
			fileName := item.FileId + item.FileExt
			originalPath := filepath.Join(fileService.fileConfig.Path, item.FilePath, fileName)
			thumbPath := filepath.Join(fileService.fileConfig.Path, item.FileThumbPath, fileName)
			originalErr := os.Remove(originalPath)
			if originalErr != nil && !os.IsNotExist(originalErr) {
				log.Printf("%v-%v: %v", Err6273, item.FileOriginalName, err)
			}
			thumbErr := os.Remove(thumbPath)
			if thumbErr != nil && !os.IsNotExist(thumbErr) {
				log.Printf("%v-%v: %v", Err6274, item.FileOriginalName, err)
			}

			if originalErr == nil || os.IsNotExist(originalErr) && thumbErr == nil || os.IsNotExist(thumbErr) {
				_, err := fileService.fileDao.UpdateRowStatusByFileIdAndFileType(item.FileId, userId, item.FileType, 3)
				if err != nil {
					log.Printf("%v-%v: %v", Err6275, item.FileOriginalName, err)
				}
			}
		}
	}

	return true, nil
}

func (fileService *FileService) DeleteInvalidFileRecordService() (int64, error) {
	result, err := fileService.fileDao.DeleteRowsByStatus(3)
	if err != nil {
		return 0, Err6276
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, Err6277
	}

	return rows, nil
}

func (fileService *FileService) UpdateRemarkService(fileId, remark string, userId int) (int64, error) {
	result, err := fileService.fileDao.UpdateRemarkByFileId(fileId, remark, userId)
	if err != nil {
		return 0, Err6279
	}
	rows, err := result.RowsAffected()

	if err != nil {
		return 0, Err6280
	}

	return rows, nil
}

func (fileService *FileService) UpdateFolderCoverService(cover, fileId string, position, userId int) (int64, error) {

	if position < 1 || position > 3 {
		return 0, Err6281
	}

	result, err := fileService.fileDao.UpdateCoverByFileId(cover, fileId, position, userId)

	if err != nil {
		return 0, Err6282
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return 0, Err6283
	}

	return rows, nil
}
