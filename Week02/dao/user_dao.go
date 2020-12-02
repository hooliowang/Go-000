package dao

import (
	"database/sql"
	"errortest/models"

	"github.com/pkg/errors"
)

func GetUserById(uid int64) (*models.User, error) {
	// do some query stuff
	err := sql.ErrNoRows

	return &models.User{}, errors.Wrap(err, "[dao.GetUserById] failed")
}
