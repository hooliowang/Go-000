package service

import (
	"errortest/dao"

	xerrors "github.com/pkg/errors"
)

func GetUser(uid int64) (bool, error) {
	user, err := dao.GetUserById(uid)
	if err != nil {
		return false, xerrors.WithMessage(err, "[service.GetUser] failed")
	}

	user.Gold++

	return true, nil
}
