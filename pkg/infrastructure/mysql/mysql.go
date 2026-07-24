package mysql

import (
	"fmt"
	"goserver/pkg/infrastructure/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMysqlConn(confg *config.Config) (sqlxDB *sqlx.DB, sqlxErr error) {
	mysqlConf := confg.Mysql

	dbDrive := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		mysqlConf.Username,
		mysqlConf.Password,
		mysqlConf.Host,
		mysqlConf.Port,
		mysqlConf.DB,
	)

	sqlxDB, sqlxErr = sqlx.Open("mysql", dbDrive)
	if sqlxErr != nil {
		return
	}
	//设置数据库连接池中允许存在的最大打开连接数
	sqlxDB.SetMaxOpenConns(mysqlConf.MaxOpenConns)
	//设置连接池中最多可以保留的空闲连接数量
	sqlxDB.SetMaxIdleConns(mysqlConf.MaxIdleConns)

	return
}
