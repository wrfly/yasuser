package types

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("custom URL already exist")
	ErrURLExpired   = errors.New("url expired")
	ErrWronPass     = errors.New("wrong password")
)
