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
			w.Write(wjson.SendMessage(7001, nil, "请登录"))
			return
		}
		claims, err := tokenParser.ParseToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(wjson.SendMessage(7001, nil, "无效的Token"))
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

type contentTypeResponseWriter struct {
	http.ResponseWriter
}

func (w contentTypeResponseWriter) WriteHeader(statusCode int) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w contentTypeResponseWriter) Write(data []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	return w.ResponseWriter.Write(data)
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有域访问
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(contentTypeResponseWriter{ResponseWriter: w}, r)
	})
}

func NewHttpServer(config *config.Config) *http.Server {
	//创建数据库连接池
	mysqlDb, mysqlDbErr := mysql.NewMysqlConn(config.Mysql)
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
	fileService := service.NewFileService(fileDao, config.Upload)

	baseController := controller.NewBaseController(config, responseJson)
	userController := controller.NewUserController(beeJwt, baseController, mux, protectedMux, userService, responseJson)
	fileController := controller.NewFileController(beeJwt, baseController, mux, protectedMux, fileService, responseJson)

	userController.BindUserController()
	fileController.BindFileController()

	mux.Handle("/", ProtectedMux(protectedMux, beeJwt, responseJson))

	return &http.Server{
		Handler: middleware(mux),
		Addr:    config.Server.Port,
	}
}
