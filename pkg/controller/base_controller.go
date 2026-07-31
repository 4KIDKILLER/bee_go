package controller

import (
	"goserver/pkg/infrastructure/config"
	"goserver/pkg/utils"
	"net/http"
)

type BaseController struct {
	config       *config.Config
	responseJson *utils.ResponseJson
}

func NewBaseController(config *config.Config, responseJson *utils.ResponseJson) *BaseController {
	return &BaseController{
		config,
		responseJson,
	}
}

func (baseController *BaseController) writeSuccess(w http.ResponseWriter, message string, data any) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(baseController.responseJson.SendSuccess(message, data))
}

func (baseController *BaseController) writeFail(w http.ResponseWriter, message string, data any) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(baseController.responseJson.SendFail(message, data))
}

func (baseController *BaseController) writeError(w http.ResponseWriter, statusCode int, message string) {
	if statusCode < 100 || statusCode > 599 {
		statusCode = http.StatusInternalServerError
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(baseController.responseJson.SendMessage(statusCode, nil, message))
}
