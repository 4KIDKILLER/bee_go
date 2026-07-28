package utils

import "encoding/json"

type ResponseJson struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
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

func (this *ResponseJson) SendSuccess(data any) []byte {
	return jsonMarshal(ResponseJson{
		Code:    200,
		Data:    data,
		Message: "操作成功",
	})
}

func (this *ResponseJson) SendFail(data any) []byte {
	return jsonMarshal(ResponseJson{
		Code:    601,
		Data:    data,
		Message: "操作失败",
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
