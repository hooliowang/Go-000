package dao

import (
	"database/sql"
	"errors"
	"errortest/models"

	xerrors "github.com/pkg/errors"
)

var (
	ErrNoRecords = errors.New("dao: ErrNoRecords")
)

func GetUserById(uid int64) (*models.User, error) {
	// do some query stuff
	err := sql.ErrNoRows
	if err == sql.ErrNoRows {
		return &models.User{}, xerrors.Wrap(ErrNoRecords, "[dao.GetUserById] failed")
	}

	return &models.User{Username: "user", Password: "pass"}, nil
}
