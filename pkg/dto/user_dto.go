package dto

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterReq struct {
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
	Password string `json:"password"`
}
