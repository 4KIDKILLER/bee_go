package server

import (
	"goserver/pkg/controller"
	"goserver/pkg/dao"
	"goserver/pkg/infrastructure/config"
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/infrastructure/mysql"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"log"
	"net/http"
)

func ProtectedMux(mux http.Handler, jwt *jwt.BeeJwt, wjson *utils.ResponseJson) http.Handler {
	//初始化http路由
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, string(wjson.SendMessage(601, nil, "请登录")), http.StatusUnauthorized)
			return
		}
		_, err := jwt.ParseToken(authHeader)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, string(wjson.SendMessage(601, nil, "Invalid Token")), http.StatusUnauthorized)
			return
		}
		mux.ServeHTTP(w, r)
	})
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
	//创建jwt验证器
	beeJwt := jwt.NewBeeJwt()
	//创建用户持久层
	userDao := dao.NewUserDao(mysqlDb)
	//创建用户服务层
	userService := service.NewUserService(userDao)
	//注册用户控制器
	userController := controller.NewUserController(beeJwt, mux, protectedMux, userService, responseJson)
	userController.BindUserController()

	mux.Handle("/", ProtectedMux(protectedMux, beeJwt, responseJson))

	return &http.Server{
		Handler: mux,
		Addr:    serverConf.Port,
	}
}
