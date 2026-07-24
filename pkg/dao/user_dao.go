package dao

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UserDao struct {
	mysql *sqlx.DB
}

func NewUserDao(mysql *sqlx.DB) (userDao *UserDao) {
	userDao = &UserDao{mysql}
	return
}

func (userDao *UserDao) Insert(username, password, avatar string) (result sql.Result, err error) {
	result, err = userDao.mysql.DB.Exec("INSERT INTO `bee_user` (username,password,avatar) VALUES (?,?,?)", username, password, avatar)

	return
}
