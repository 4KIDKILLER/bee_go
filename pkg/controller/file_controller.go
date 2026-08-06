package controller

import (
	"encoding/json"
	"fmt"
	"goserver/pkg/dto"
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"goserver/pkg/vo"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// 错误码范围6200-6250
type FileController struct {
	*BaseController
	jwt          *jwt.BeeJwt
	mux          *http.ServeMux
	protectedMux *http.ServeMux
	fileService  *service.FileService
}

func NewFileController(
	jwt *jwt.BeeJwt,
	baseController *BaseController,
	mux, protectedMux *http.ServeMux,
	fileService *service.FileService,
	responseJson *utils.ResponseJson,
) (
	fileController *FileController,
) {
	fileController = &FileController{
		BaseController: baseController,
		jwt:            jwt,
		mux:            mux,
		protectedMux:   protectedMux,
		fileService:    fileService,
	}
	return
}

func (fileController *FileController) BindFileController() {
	/*
		文件上传
	*/
	fileController.protectedMux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {

		const maxFileSize int64 = 10 << 20

		/*
			限制请求体大小（例如 10MB）。如不设置，若有人直接调用接口上传一个20GB文件，后端不会在10MiB时拒绝，
			而是先接收完整的20GB。可能会导致文件耗尽服务器资源。通过设置http.MaxBytesReader，若超出请求体最大
			限制服务端会直接重置请求，并且不会解析文件。
		*/
		r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

		//解析 multipart 表单，32MB 以内的文件会存内存，更大则存临时文件
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			fileController.writeError(w, http.StatusRequestEntityTooLarge, "上传文件过大")
			return
		}
		defer r.MultipartForm.RemoveAll()

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			fileController.writeError(w, http.StatusBadRequest, "获取文件出错")
			return
		}
		defer file.Close()

		if fileHeader.Size > maxFileSize {
			fileController.writeError(w, http.StatusRequestEntityTooLarge, "图片不能超过10MB")
			return
		}

		if len(fileHeader.Filename) > 100 {
			fileController.writeFail(w, "文件名称需小于50字符", nil)
			return
		}

		parentId := r.FormValue("parentId")
		tags := r.FormValue("tags")
		remark := r.FormValue("remark")

		//【重要】安全处理文件名，防止路径穿越攻击
		// 使用 filepath.Base 去除任何路径信息，仅保留文件名本身
		safeFilename := strings.ToLower(filepath.Base(strings.ReplaceAll(fileHeader.Filename, "\\", "/")))
		if safeFilename == "" || safeFilename == "." || safeFilename == ".." {
			fileController.writeError(w, http.StatusBadRequest, "无效的文件名")
			return
		}

		beeClaims, _ := jwt.ClaimsFromContext(r.Context())

		_, resultErr := fileController.fileService.UploadFileService(file, parentId, safeFilename, tags, remark, fileHeader.Size, beeClaims.UserId)
		if resultErr != nil {
			fileController.writeFail(w, resultErr.Error(), nil)
		} else {
			fileController.writeSuccess(w, "文件上传成功", nil)
		}
	})
	/*
		创建文件夹
	*/
	fileController.protectedMux.HandleFunc("POST /createFolder", func(w http.ResponseWriter, r *http.Request) {

		var createFolderReq dto.CreateFolderReq

		decodeErr := json.NewDecoder(r.Body).Decode(&createFolderReq)
		if decodeErr != nil {
			fileController.writeFail(w, "参数解析失败", nil)
			return
		}
		beeClaims, _ := jwt.ClaimsFromContext(r.Context())
		_, insertErr := fileController.fileService.CreateFolderService(&createFolderReq, beeClaims.UserId)

		if insertErr != nil {
			fileController.writeFail(w, "文件夹创建失败", nil)
		} else {
			fileController.writeSuccess(w, "文件夹创建成功", nil)
		}
	})
	/*
		获取用户文件列表
	*/
	fileController.protectedMux.HandleFunc("GET /getFileList", func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()

		parentId := query.Get("parentId")
		pageStr := query.Get("page")
		if pageStr == "" {
			fileController.writeFail(w, "缺少page参数", nil)
			return
		}
		pageSizeStr := query.Get("pageSize")
		if pageStr == "" {
			fileController.writeFail(w, "缺少pageSize参数", nil)
			return
		}
		beeClaims, _ := jwt.ClaimsFromContext(r.Context())

		page, pageErr := strconv.Atoi(pageStr)
		pageSize, pageSizeErr := strconv.Atoi(pageSizeStr)
		if pageErr != nil || pageSizeErr != nil {
			fileController.writeFail(w, "分页参数异常", nil)
			return
		}

		fileCount, fileList, err := fileController.fileService.GetUserFileList(parentId, beeClaims.UserId, page, pageSize)

		if err != nil {
			fileController.writeFail(w, "获取文件列表失败", nil)
			return
		}

		dataList := make([]vo.FileListVo, 0, len(fileList))

		for _, item := range fileList {
			covers := [3]string{item.Cover1, item.Cover2, item.Cover3}
			tags := make([]string, 0, 3)
			if item.Tags != "" {
				tags = strings.Split(item.Tags, ",")
			}
			name := item.FileId + item.FileExt
			src := ""
			thumbSrc := ""
			if item.FileType == 2 {
				src = fmt.Sprintf("%s/%s/%s", fileController.config.Upload.Host, item.FilePath, name)
				thumbSrc = fmt.Sprintf("%s/%s/%s", fileController.config.Upload.Host, item.FileThumbPath, name)
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

		resultData := utils.PaginationJson[vo.FileListVo]{
			List:     dataList,
			Total:    fileCount,
			Page:     page,
			PageSize: pageSize,
		}

		fileController.writeSuccess(w, "获取文件列表成功", resultData)
	})
}
