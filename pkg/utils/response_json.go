package utils

import "encoding/json"

type ResponseJson struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func jsonMarshal(this *ResponseJson) []byte {
	result, err := json.Marshal(this)
	if err != nil {
		return []byte(`{
			data: "",
			code: 501,
			message: "响应序列化失败:"` + `err.Error()
		}`)
	}
	return result
}

func (this *ResponseJson) SendSuccess(data any) []byte {
	this.Code = 200
	this.Data = data
	this.Message = "操作成功"

	return jsonMarshal(this)
}

func (this *ResponseJson) SendFail(data any) []byte {
	this.Code = 601
	this.Data = data
	this.Message = "操作失败"

	return jsonMarshal(this)
}

func (this *ResponseJson) SendMessage(code int, data any, message string) []byte {
	this.Code = code
	this.Data = data
	this.Message = message

	return jsonMarshal(this)
}

func NewResponseJson() (responseJson *ResponseJson) {
	responseJson = &ResponseJson{}
	return
}
