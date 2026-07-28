package controller

import (
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileController struct {
	jwt          *jwt.BeeJwt
	mux          *http.ServeMux
	protectedMux *http.ServeMux
	fileService  *service.FileService
	responseJson *utils.ResponseJson
}

func NewFileController(
	jwt *jwt.BeeJwt,
	mux, protectedMux *http.ServeMux,
	fileService *service.FileService,
	responseJson *utils.ResponseJson,
) (
	fileController *FileController,
) {
	fileController = &FileController{
		jwt,
		mux,
		protectedMux,
		fileService,
		responseJson,
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
			http.Error(w, "File too large or parse error", http.StatusBadRequest)
			return
		}
		file, handler, fileErr := r.FormFile("file")
		if fileErr != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		uploadDir := "./file"
		if mkdirErr := os.MkdirAll(uploadDir, 0755); mkdirErr != nil {
			http.Error(w, "Error creating upload dir", http.StatusInternalServerError)
			return
		}

		//【重要】安全处理文件名，防止路径穿越攻击
		// 使用 filepath.Base 去除任何路径信息，仅保留文件名本身
		safeFilename := filepath.Base(strings.ReplaceAll(handler.Filename, "\\", "/"))
		if safeFilename == "" || safeFilename == "." || safeFilename == ".." {
			http.Error(w, "Invalid filename", http.StatusBadRequest)
			return
		}

		fileId, uuidErr := uuid.NewRandom()
		if uuidErr != nil {
			http.Error(w, "Error generating filename", http.StatusInternalServerError)
			return
		}
		randomFilename := strings.ReplaceAll(fileId.String(), "-", "") + filepath.Ext(safeFilename)
		dstPath := filepath.Join(uploadDir, randomFilename)

		//创建目标文件并保存上传内容
		dst, dstErr := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if dstErr != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, copyErr := io.Copy(dst, file)
		if copyErr != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
	})
}
