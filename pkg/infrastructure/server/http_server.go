package server

import (
	"goserver/pkg/controller"
	"goserver/pkg/dao"
	"goserver/pkg/infrastructure/config"
	beeJwt "goserver/pkg/infrastructure/jwt"
	"goserver/pkg/infrastructure/mysql"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"net/http"
	"strings"
)

func ProtectedMux(mux http.Handler, tokenParser *beeJwt.BeeJwt, wjson *utils.ResponseJson) http.Handler {
	//初始化http路由
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := tokenFromRequest(r)
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(wjson.SendMessage(601, nil, "请登录"))
			return
		}
		claims, err := tokenParser.ParseToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(wjson.SendMessage(601, nil, "Invalid Token"))
			return
		}
		mux.ServeHTTP(w, r.WithContext(beeJwt.WithClaims(r.Context(), claims)))
	})
}

func tokenFromRequest(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader != "" {

		authParts := strings.Fields(authHeader)
		if len(authParts) == 2 && strings.EqualFold(authParts[0], "Bearer") {
			return authParts[1]
		}
		return authHeader
	} else {
		return ""
	}
}

func NewHttpServer(config *config.Config) *http.Server {
	serverConf := config.Server
	mysqlConf := config.Mysql

	//创建数据库连接池
	mysqlDb, mysqlDbErr := mysql.NewMysqlConn(&mysqlConf)
	if mysqlDbErr != nil {
		panic("数据库链接创建失败:" + mysqlDbErr.Error())
	}

	responseJson := utils.NewResponseJson()
	mux := http.NewServeMux()
	protectedMux := http.NewServeMux()
	beeJwt := beeJwt.NewBeeJwt()

	userDao := dao.NewUserDao(mysqlDb)
	fileDao := dao.NewFileDao(mysqlDb)

	userService := service.NewUserService(userDao)
	fileService := service.NewFileService(fileDao)

	userController := controller.NewUserController(beeJwt, mux, protectedMux, userService, responseJson)
	fileController := controller.NewFileController(beeJwt, mux, protectedMux, fileService, responseJson)

	userController.BindUserController()
	fileController.BindFileController()

	mux.Handle("/", ProtectedMux(protectedMux, beeJwt, responseJson))

	return &http.Server{
		Handler: mux,
		Addr:    serverConf.Port,
	}
}
