package controller

import (
	"encoding/json"
	"goserver/pkg/dto"
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"net/http"
)

// 错误码范围6100-6150
type UserController struct {
	*BaseController
	jwt          *jwt.BeeJwt
	mux          *http.ServeMux
	protectedMux *http.ServeMux
	userService  *service.UserService
}

func NewUserController(
	jwt *jwt.BeeJwt,
	baseController *BaseController,
	mux, protectedMux *http.ServeMux,
	userService *service.UserService,
	responseJson *utils.ResponseJson,
) (
	userController *UserController,
) {
	userController = &UserController{
		BaseController: baseController,
		jwt:            jwt,
		mux:            mux,
		protectedMux:   protectedMux,
		userService:    userService,
	}
	return
}

func (userController *UserController) BindUserController() {
	/*
		获取用户信息
	*/
	userController.protectedMux.HandleFunc("GET /getUserInfo", func(w http.ResponseWriter, r *http.Request) {

		data := []string{
			"123123",
			"123123",
			"231321",
		}

		userController.writeSuccess(w, "", data)
	})

	/*
		用户登录
	*/
	userController.mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var loginReq dto.LoginReq

		err := json.NewDecoder(r.Body).Decode(&loginReq)

		if err != nil {
			userController.writeError(w, http.StatusBadRequest, "参数解析失败")
			return
		}

		beeUser, loginErr := userController.userService.LoginService(loginReq.Username, loginReq.Password)
		if loginErr != nil {
			userController.writeFail(w, loginErr.Error(), nil)
			return
		}

		token, tokenErr := userController.jwt.GenerateToken(beeUser.Username, beeUser.UserId)
		if tokenErr != nil {
			userController.writeFail(w, tokenErr.Error(), nil)
			return
		}

		loginInfo := map[string]string{
			"avatar":   beeUser.Avatar,
			"username": beeUser.Username,
			"token":    "Bearer " + token,
		}
		userController.writeSuccess(w, "登录成功", loginInfo)
	})

	//用户注册
	userController.mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		var registerReq dto.RegisterReq

		err := json.NewDecoder(r.Body).Decode(&registerReq)

		if err != nil {
			userController.writeError(w, http.StatusBadRequest, "参数解析失败")
			return
		}

		_, registerErr := userController.userService.UserRegisterService(&registerReq)
		if registerErr != nil {
			userController.writeFail(w, registerErr.Error(), nil)
			return
		}

		userController.writeSuccess(w, "注册成功", nil)
	})
}
