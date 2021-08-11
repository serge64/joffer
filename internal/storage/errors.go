package storage

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrUserNotValid   = errors.New("user not valid")
	ErrEmailNotUnique = errors.New("email not unique")
)
