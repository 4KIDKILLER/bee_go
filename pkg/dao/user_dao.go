package dao

import (
	"database/sql"
	"goserver/pkg/model"

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
	result, err = userDao.mysql.Exec("INSERT INTO `bee_user` (username,password,avatar) VALUES (?,?,?)", username, password, avatar)

	return
}

func (userDao *UserDao) QueryUserByNameAndPassword(username, password string) (beeUser *model.BeeUser, err error) {
	beeUser = &model.BeeUser{}
	err = userDao.mysql.Get(beeUser, "SELECT `username`,`avatar`,`status` FROM `bee_user` WHERE `username`=? AND `password`=? AND `status`=1", username, password)

	return
}

func (userDao *UserDao) CountUserByName(username string) (count int, err error) {
	err = userDao.mysql.Get(&count, "SELECT COUNT(`username`) FROM `bee_user` WHERE `username`=?", username)

	return
}
