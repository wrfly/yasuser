package types

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("custom URL already exist")
)
