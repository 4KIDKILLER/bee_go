package controller

import (
	"encoding/json"
	"goserver/pkg/dto"
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"net/http"
)

type FileTagController struct {
	*BaseController
	jwt            *jwt.BeeJwt
	mux            *http.ServeMux
	protectedMux   *http.ServeMux
	fileTagService *service.FileTagService
}

func NewFileTagController(
	jwt *jwt.BeeJwt,
	baseController *BaseController,
	mux, protectedMux *http.ServeMux,
	fileTagService *service.FileTagService,
	responseJson *utils.ResponseJson,
) (
	fileTagController *FileTagController,
) {
	fileTagController = &FileTagController{
		BaseController: baseController,
		jwt:            jwt,
		mux:            mux,
		protectedMux:   protectedMux,
		fileTagService: fileTagService,
	}
	return
}

func (fileTagController *FileTagController) BindFileTagController() {
	fileTagController.protectedMux.HandleFunc("POST /createTarget", func(w http.ResponseWriter, r *http.Request) {
		var createTargetReq dto.CreateTargetReq
		decodeErr := json.NewDecoder(r.Body).Decode(&createTargetReq)

		if decodeErr != nil {
			fileTagController.writeError(w, http.StatusBadRequest, "参数解析失败")
			return
		}

		beeClaims, _ := jwt.ClaimsFromContext(r.Context())

		_, err := fileTagController.fileTagService.CreateFileTagService(createTargetReq.TagName, createTargetReq.FileId, beeClaims.UserId)
		if err != nil {
			fileTagController.writeFail(w, "标签创建失败", err)
		} else {
			fileTagController.writeSuccess(w, "标签创建成功", nil)
		}
	})
	/*
		删除文件标签
	*/
	fileTagController.protectedMux.HandleFunc("POST /deleteTarget", func(w http.ResponseWriter, r *http.Request) {
		var deleteTargetReq dto.DeleteTargetReq

		decodeErr := json.NewDecoder(r.Body).Decode(&deleteTargetReq)

		if decodeErr != nil {
			fileTagController.writeError(w, http.StatusBadRequest, "参数解析失败")
			return
		}

		_, err := fileTagController.fileTagService.DeleteFileTagService(deleteTargetReq.Id)

		if err != nil {
			fileTagController.writeFail(w, "标签删除失败", err)
		} else {
			fileTagController.writeSuccess(w, "标签删除成功", nil)
		}
	})
}
