package controller

import (
	"encoding/json"
	"goserver/pkg/dto"
	"goserver/pkg/infrastructure/dictionary"
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

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

		fileList, err := fileController.fileService.GetUserFileListService(fileController.config.Upload.Host, parentId, beeClaims.UserId, page, pageSize)

		if err != nil {
			fileController.writeFail(w, "获取文件列表失败", nil)
			return
		}

		fileController.writeSuccess(w, "获取文件列表成功", fileList)
	})
	/*
		软删除文件或文件夹
	*/
	fileController.protectedMux.HandleFunc("POST /deleteSoft", func(w http.ResponseWriter, r *http.Request) {
		var deleteSoftReq dto.DeleteSoftReq
		decodeErr := json.NewDecoder(r.Body).Decode(&deleteSoftReq)
		if decodeErr != nil {
			fileController.writeFail(w, "参数解析失败", nil)
			return
		}

		beeClaims, _ := jwt.ClaimsFromContext(r.Context())
		var err error
		if deleteSoftReq.Type == 1 {
			_, err = fileController.fileService.DeleteFolderSoftService(deleteSoftReq.Id, beeClaims.UserId)
		} else {
			_, err = fileController.fileService.DeleteFileSoftService(deleteSoftReq.Id, beeClaims.UserId, deleteSoftReq.Type)
		}
		fileTypeName := dictionary.FileTypeDict[deleteSoftReq.Type]
		if err != nil {
			fileController.writeFail(w, fileTypeName+"删除失败", nil)
		} else {
			fileController.writeSuccess(w, fileTypeName+"删除成功", nil)
		}
	})
	/*
		硬删除文件
	*/
	fileController.protectedMux.HandleFunc("POST /deleteHard", func(w http.ResponseWriter, r *http.Request) {
		beeClaims, _ := jwt.ClaimsFromContext(r.Context())
		_, err := fileController.fileService.DeleteFolderHardService(beeClaims.UserId)
		if err != nil {
			fileController.writeFail(w, "删除失败", nil)
		} else {
			fileController.writeFail(w, "删除成功", nil)
		}

	})
	/*
		修改文件/文件夹名称
	*/
	fileController.protectedMux.HandleFunc("POST /rename", func(w http.ResponseWriter, r *http.Request) {
		var renameReq dto.RenameReq
		decodeErr := json.NewDecoder(r.Body).Decode(&renameReq)
		if decodeErr != nil {
			fileController.writeFail(w, "参数解析失败", nil)
			return
		}

		beeClaims, _ := jwt.ClaimsFromContext(r.Context())
		_, renameErr := fileController.fileService.UpdateOriginalNameService(renameReq.Name, renameReq.Id, beeClaims.UserId)

		if renameErr != nil {
			fileController.writeFail(w, "修改失败", nil)
		} else {
			fileController.writeSuccess(w, "修改成功", nil)
		}
	})

	/*
		获取文件夹结构
	*/
	fileController.protectedMux.HandleFunc("GET /getFolderTree", func(w http.ResponseWriter, r *http.Request) {
		beeClaims, _ := jwt.ClaimsFromContext(r.Context())
		result, err := fileController.fileService.GetUserFileTreeService(beeClaims.UserId)
		if err != nil {
			fileController.writeFail(w, "获取文件夹列表失败", nil)
		} else {
			fileController.writeSuccess(w, "获取文件夹列表成功", result)
		}
	})

	/*
		清理无效文件数据
	*/
	fileController.protectedMux.HandleFunc("POST /clearInvalidRecord", func(w http.ResponseWriter, r *http.Request) {
		result, err := fileController.fileService.DeleteInvalidFileRecordService()
		if err != nil {
			fileController.writeFail(w, "无效记录清理失败", err)
		} else {
			fileController.writeSuccess(w, "无效记录清理成功，共删除"+strconv.FormatInt(result, 10)+"条数据", result)
		}
	})
}
