package utils

import "encoding/json"

type ResponseJson struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type PaginationJson[T any] struct {
	List     []T `json:"list"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func jsonMarshal(response ResponseJson) []byte {
	result, err := json.Marshal(response)
	if err != nil {
		fallback := ResponseJson{
			Code:    501,
			Data:    nil,
			Message: "响应序列化失败:" + err.Error(),
		}
		result, _ = json.Marshal(fallback)
	}
	return result
}

func (this *ResponseJson) SendSuccess(message string, data any) []byte {
	if message == "" {
		message = "操作成功"
	}
	return jsonMarshal(ResponseJson{
		Code:    200,
		Data:    data,
		Message: message,
	})
}

func (this *ResponseJson) SendFail(message string, data any) []byte {
	if message == "" {
		message = "操作失败"
	}
	return jsonMarshal(ResponseJson{
		Code:    601,
		Data:    data,
		Message: message,
	})
}

func (this *ResponseJson) SendMessage(code int, data any, message string) []byte {
	return jsonMarshal(ResponseJson{
		Code:    code,
		Data:    data,
		Message: message,
	})
}

func NewResponseJson() (responseJson *ResponseJson) {
	responseJson = &ResponseJson{}
	return
}
