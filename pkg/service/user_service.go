package service

import (
	"errors"
	"goserver/pkg/dao"
	"goserver/pkg/dto"
	"goserver/pkg/model"
	"log"
	"strings"
)

var (
	Err6150 = errors.New("6150: 用户名不能为空")
	Err6151 = errors.New("6151: 用户查询失败")
	Err6152 = errors.New("6152: 密码不能为空")
	Err6153 = errors.New("6153: 用户注册失败")
	Err6154 = errors.New("6154: 用户名或密码错误")
	Err6155 = errors.New("6155: 用户名不能为空")
	Err6156 = errors.New("6156: 密码不能为空")
	Err6157 = errors.New("6157: 注册失败")
	Err6158 = errors.New("6158: 当前用户名已存在")
	Err6159 = errors.New("6159: 注册失败")
	Err6160 = errors.New("6160: 注册失败")
	Err6161 = errors.New("6161: 注册失败")
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) (userService *UserService) {
	userService = &UserService{userDao}
	return
}

func (userService *UserService) GetUserInfoService(userId int) (beeUser *model.BeeUser, err error) {
	beeUser, err = userService.userDao.QueryUserByUId(userId)
	if err != nil {
		return nil, Err6151
	}
	return
}

func (userService *UserService) LoginService(username, password string) (beeUser *model.BeeUser, err error) {
	if strings.TrimSpace(username) == "" {
		return nil, Err6150
	}
	if strings.TrimSpace(password) == "" {
		return nil, Err6152
	}
	beeUser, err = userService.userDao.QueryUserByNameAndPassword(username, password)

	if err != nil {
		log.Printf("%v: %v", Err6154, err)
		return nil, Err6154
	}

	return
}

func (userService *UserService) UserRegisterService(registerReq *dto.RegisterReq) (bool, error) {
	if strings.TrimSpace(registerReq.Username) == "" {
		return false, Err6155
	}

	if strings.TrimSpace(registerReq.Password) == "" {
		return false, Err6156
	}
	if strings.TrimSpace(registerReq.Avatar) == "" {
		registerReq.Avatar = "https://go.dev/images/go-logo-white.svg"
	}

	count, countErr := userService.userDao.CountUserByName(registerReq.Username)
	if countErr != nil {
		return false, Err6157
	}
	if count > 0 {
		return false, Err6158
	}

	countUser, _ := userService.userDao.CountByUser()

	userId := 100000 + countUser

	insert, insertErr := userService.userDao.Insert(registerReq.Username, registerReq.Password, registerReq.Avatar, userId)
	if insertErr != nil {
		return false, Err6159
	}

	rowsAffected, err := insert.RowsAffected()
	if err != nil {
		return false, Err6160
	}
	if rowsAffected != 1 {
		return false, Err6161
	}

	return true, nil
}
