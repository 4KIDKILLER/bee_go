package controller

import (
	"encoding/json"
	"goserver/pkg/dto"
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"net/http"
)

type UserController struct {
	*BaseController
	mux          *http.ServeMux
	protectedMux *http.ServeMux
	userService  *service.UserService
}

func NewUserController(
	baseController *BaseController,
	mux, protectedMux *http.ServeMux,
	userService *service.UserService,
	responseJson *utils.ResponseJson,
) (
	userController *UserController,
) {
	userController = &UserController{
		BaseController: baseController,
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

		loginInfo, err := userController.userService.LoginService(loginReq.Username, loginReq.Password)
		if err != nil {
			userController.writeFail(w, err.Error(), nil)
			return
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

	//私密模式验证
	userController.protectedMux.HandleFunc("POST /private", func(w http.ResponseWriter, r *http.Request) {
		var privateReq dto.PrivateReq

		err := json.NewDecoder(r.Body).Decode(&privateReq)

		if err != nil {
			userController.writeError(w, http.StatusBadRequest, "参数解析失败")
			return
		}

		beeClaims, _ := jwt.ClaimsFromContext(r.Context())

		loginInfo, err := userController.userService.UserPrivateLoginService(beeClaims.UserId, privateReq.Password)

		if err != nil {
			userController.writeFail(w, err.Error(), nil)
			return
		}

		userController.writeSuccess(w, "验证成功", loginInfo)
	})

	//用户token参数调试
	userController.protectedMux.HandleFunc("GET /tokenInfo", func(w http.ResponseWriter, r *http.Request) {
		beeClaims, _ := jwt.ClaimsFromContext(r.Context())

		userController.writeSuccess(w, "解析成功", beeClaims)
	})
}
