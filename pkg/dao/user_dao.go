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

func (userDao *UserDao) Insert(username, password, avatar string, userId int) (result sql.Result, err error) {
	result, err = userDao.mysql.Exec("INSERT INTO `bee_user` (`username`,`password`,`avatar`,`user_id`) VALUES (?,?,?,?)", username, password, avatar, userId)

	return
}

func (userDao *UserDao) QueryUserByNameAndPassword(username, password string) (beeUser *model.BeeUser, err error) {
	beeUser = &model.BeeUser{}
	err = userDao.mysql.Get(beeUser, "SELECT `username`,`avatar`,`status`,`user_id` FROM `bee_user` WHERE `username`=? AND `password`=? AND `status`=1", username, password)

	return
}

func (userDao *UserDao) CountByUser() (count int, err error) {
	err = userDao.mysql.Get(&count, "SELECT COUNT(`username`) FROM `bee_user`")

	return
}

func (userDao *UserDao) QueryUserByUId(userId int) (beeUser *model.BeeUser, err error) {
	beeUser = &model.BeeUser{}
	err = userDao.mysql.Get(beeUser, "SELECT `user_id`,`username`,`create_time`,`update_time`,`avatar`,`status` FROM `bee_user` WHERE `user_id`=?", userId)

	return
}

func (userDao *UserDao) CountUserByName(username string) (count int, err error) {
	err = userDao.mysql.Get(&count, "SELECT COUNT(`username`) FROM `bee_user` WHERE `username`=?", username)

	return
}
