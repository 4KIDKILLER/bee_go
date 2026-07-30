package controller

import (
	"encoding/json"
	"goserver/pkg/dto"
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"log"
	"net/http"
	"path/filepath"
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

		// 限制请求体大小（例如 10MB），防止大文件耗尽服务器资源
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		//解析 multipart 表单，32MB 以内的文件会存内存，更大则存临时文件
		parseErr := r.ParseMultipartForm(32 << 20)
		if parseErr != nil {
			fileController.writeError(w, http.StatusBadRequest, "文件过大或解析错误")
			return
		}
		file, fileHeader, fileErr := r.FormFile("file")
		if fileErr != nil {
			fileController.writeError(w, http.StatusBadRequest, "获取文件出错")
			return
		}
		defer file.Close()
		parentId := r.FormValue("parentId")

		//【重要】安全处理文件名，防止路径穿越攻击
		// 使用 filepath.Base 去除任何路径信息，仅保留文件名本身
		safeFilename := filepath.Base(strings.ReplaceAll(fileHeader.Filename, "\\", "/"))
		if safeFilename == "" || safeFilename == "." || safeFilename == ".." {
			fileController.writeError(w, http.StatusBadRequest, "无效的文件名")
			return
		}

		beeClaims, _ := jwt.ClaimsFromContext(r.Context())

		_, resultErr := fileController.fileService.CreateFileService(file, parentId, safeFilename, fileHeader.Size, beeClaims.UserId)
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

		_, insertErr := fileController.fileService.CreateFolderService(createFolderReq.ParentId, createFolderReq.FolderName, beeClaims.UserId)

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

		var getFileListReq dto.GetFileListReq

		decodeErr := json.NewDecoder(r.Body).Decode(&getFileListReq)
		if decodeErr != nil {
			fileController.writeFail(w, "参数解析失败", nil)
			return
		}

		beeClaims, _ := jwt.ClaimsFromContext(r.Context())

		fileList, err := fileController.fileService.GetUserFileList(getFileListReq.ParentId, beeClaims.UserId, getFileListReq.Page, getFileListReq.PageSize)

		if err != nil {
			log.Println(err)
			fileController.writeFail(w, "获取文件列表失败", nil)
		} else {
			fileController.writeSuccess(w, "获取文件列表成功", fileList)
		}
	})
}
