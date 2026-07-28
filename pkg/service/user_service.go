package service

import (
	"errors"
	"fmt"
	"goserver/pkg/dao"
	"goserver/pkg/dto"
	"goserver/pkg/model"
	"strings"
)

var (
	ErrUsernameRequired = errors.New("用户名不能为空")
	ErrGetUserFailed    = errors.New("用户查询失败")
	ErrPasswordRequired = errors.New("密码不能为空")
	ErrRegisterFailed   = errors.New("用户注册失败")
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
		return nil, fmt.Errorf("%w: %v", ErrGetUserFailed, err)
	}
	return
}

func (userService *UserService) LoginService(username, password string) (beeUser *model.BeeUser, err error) {
	if strings.TrimSpace(username) == "" {
		return nil, ErrUsernameRequired
	}
	if strings.TrimSpace(password) == "" {
		return nil, ErrPasswordRequired
	}
	beeUser, err = userService.userDao.QueryUserByNameAndPassword(username, password)

	return
}

func (userService *UserService) UserRegisterService(registerReq *dto.RegisterReq) (bool, error) {
	if strings.TrimSpace(registerReq.Username) == "" {
		return false, ErrUsernameRequired
	}

	if strings.TrimSpace(registerReq.Password) == "" {
		return false, ErrPasswordRequired
	}
	if strings.TrimSpace(registerReq.Avatar) == "" {
		registerReq.Avatar = "https://go.dev/images/go-logo-white.svg"
	}

	count, countErr := userService.userDao.CountUserByName(registerReq.Username)
	if countErr != nil {
		return false, fmt.Errorf("%w: %v", ErrRegisterFailed, countErr)
	}
	if count > 0 {
		return false, fmt.Errorf("%w: %v", ErrRegisterFailed, "当前用户名已存在")
	}

	countUser, _ := userService.userDao.CountByUser()

	userId := 100000 + countUser

	insert, insertErr := userService.userDao.Insert(registerReq.Username, registerReq.Password, registerReq.Avatar, userId)
	if insertErr != nil {
		return false, fmt.Errorf("%w: %v", ErrRegisterFailed, insertErr)
	}

	rowsAffected, err := insert.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%w: 注册失败: %v", ErrRegisterFailed, err)
	}
	if rowsAffected != 1 {
		return false, ErrRegisterFailed
	}

	return true, nil
}
