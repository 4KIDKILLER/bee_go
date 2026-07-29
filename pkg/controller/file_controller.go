package controller

import (
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/service"
	"goserver/pkg/utils"
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
	fileController.protectedMux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {

		// 限制请求体大小（例如 10MB），防止大文件耗尽服务器资源
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		//解析 multipart 表单，32MB 以内的文件会存内存，更大则存临时文件
		parseErr := r.ParseMultipartForm(32 << 20)
		if parseErr != nil {
			fileController.writeError(w, http.StatusBadRequest, "文件过大或解析错误")
			return
		}
		file, handler, fileErr := r.FormFile("file")
		if fileErr != nil {
			fileController.writeError(w, http.StatusBadRequest, "获取文件出错")
			return
		}
		defer file.Close()

		parentId := r.FormValue("parentId")

		//【重要】安全处理文件名，防止路径穿越攻击
		// 使用 filepath.Base 去除任何路径信息，仅保留文件名本身
		safeFilename := filepath.Base(strings.ReplaceAll(handler.Filename, "\\", "/"))
		if safeFilename == "" || safeFilename == "." || safeFilename == ".." {
			fileController.writeError(w, http.StatusBadRequest, "无效的文件名")
			return
		}
		beeClaims, ok := jwt.ClaimsFromContext(r.Context())
		if !ok || beeClaims == nil {
			fileController.writeFail(w, "请登录", nil)
			return
		}
		_, resultErr := fileController.fileService.CreateFileService(file, parentId, safeFilename, beeClaims.UserId)
		if resultErr != nil {
			fileController.writeFail(w, resultErr.Error(), nil)
			return
		} else {
			fileController.writeSuccess(w, "", nil)
			return
		}
	})
}
