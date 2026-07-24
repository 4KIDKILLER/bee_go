package service

import (
	"errors"
	"fmt"
	"goserver/pkg/dao"
	"strings"
)

var (
	ErrUsernameRequired = errors.New("用户名不能为空")
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

func (userService *UserService) UserRegister(username, password, avatar string) (bool, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return false, ErrUsernameRequired
	}
	if strings.TrimSpace(password) == "" {
		return false, ErrPasswordRequired
	}
	if strings.TrimSpace(avatar) == "" {
		avatar = "https://go.dev/images/go-logo-white.svg"
	}

	result, err := userService.userDao.Insert(username, password, avatar)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrRegisterFailed, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%w: 获取写入结果失败: %v", ErrRegisterFailed, err)
	}
	if rowsAffected != 1 {
		return false, ErrRegisterFailed
	}

	return true, nil
}
