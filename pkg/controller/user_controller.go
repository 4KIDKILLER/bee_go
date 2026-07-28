package controller

import (
	"encoding/json"
	"goserver/pkg/dto"
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"log"
	"net/http"
)

type UserController struct {
	jwt          *jwt.BeeJwt
	mux          *http.ServeMux
	protectedMux *http.ServeMux
	userService  *service.UserService
	responseJson *utils.ResponseJson
}

func NewUserController(
	jwt *jwt.BeeJwt,
	mux, protectedMux *http.ServeMux,
	userService *service.UserService,
	responseJson *utils.ResponseJson,
) (
	userController *UserController,
) {
	userController = &UserController{
		jwt,
		mux,
		protectedMux,
		userService,
		responseJson,
	}
	return
}

func (this *UserController) BindUserController() {
	/*
		获取用户信息
	*/
	this.protectedMux.HandleFunc("GET /getUserInfo", func(w http.ResponseWriter, r *http.Request) {

		data := []string{
			"123123",
			"123123",
			"231321",
		}

		result := this.responseJson.SendSuccess(data)

		w.Write(result)
	})

	/*
		用户登录
	*/
	this.mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var loginReq dto.LoginReq

		err := json.NewDecoder(r.Body).Decode(&loginReq)

		if err != nil {
			this.responseJson.SendMessage(602, nil, "参数解析失败:"+err.Error())
		}

		beeUser, loginErr := this.userService.LoginService(loginReq.Username, loginReq.Password)
		var result []byte
		if loginErr != nil {
			log.Println(loginErr.Error())
			result = this.responseJson.SendMessage(601, nil, loginErr.Error())
		} else {
			token, tokenErr := this.jwt.GenerateToken(beeUser.Username, beeUser.UserId)
			if tokenErr != nil {
				result = this.responseJson.SendMessage(601, nil, tokenErr.Error())
			} else {
				loginInfo := map[string]string{
					"avatar":   beeUser.Avatar,
					"username": beeUser.Username,
					"token":    "Bearer " + token,
				}
				result = this.responseJson.SendSuccess(loginInfo)
			}
		}
		w.Write(result)
	})

	//用户注册
	this.mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		var registerReq dto.RegisterReq

		err := json.NewDecoder(r.Body).Decode(&registerReq)

		if err != nil {
			this.responseJson.SendMessage(602, nil, "参数解析失败:"+err.Error())
		}

		_, registerErr := this.userService.UserRegisterService(&registerReq)
		var result []byte
		if registerErr != nil {
			log.Println(registerErr.Error())
			result = this.responseJson.SendMessage(601, nil, registerErr.Error())
		} else {
			result = this.responseJson.SendSuccess(nil)
		}
		w.Write(result)
	})
}
