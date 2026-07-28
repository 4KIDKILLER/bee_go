package server

import (
	"goserver/pkg/controller"
	"goserver/pkg/dao"
	"goserver/pkg/infrastructure/config"
	"goserver/pkg/infrastructure/jwt"
	"goserver/pkg/infrastructure/mysql"
	"goserver/pkg/service"
	"goserver/pkg/utils"
	"net/http"
)

func NewHttpServer(config *config.Config) *http.Server {
	serverConf := config.Server
	mysqlConf := config.Mysql

	//初始化http路由
	mux := http.NewServeMux()

	//创建数据库连接池
	mysqlDb, mysqlDbErr := mysql.NewMysqlConn(&mysqlConf)
	if mysqlDbErr != nil {
		panic("数据库链接创建失败:" + mysqlDbErr.Error())
	}

	responseJson := utils.NewResponseJson()

	//创建jwt验证器
	jwt := jwt.NewBeeJwt()
	//创建用户持久层
	userDao := dao.NewUserDao(mysqlDb)
	//创建用户服务层
	userService := service.NewUserService(userDao)
	//注册用户控制器
	userController := controller.NewUserController(jwt, mux, userService, responseJson)
	userController.BindUserController()

	return &http.Server{
		Handler: mux,
		Addr:    serverConf.Port,
	}
}
