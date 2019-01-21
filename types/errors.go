package types

import "errors"

var (
	ErrNotFound     = errors.New("not found\n")
	ErrAlreadyExist = errors.New("custom URL already exist\n")
	ErrURLExpired   = errors.New("url expired\n")
	ErrWronPass     = errors.New("wrong password\n")
	ErrURLTooLong   = errors.New("URL too long\n")
	ErrSameHost     = errors.New("can not shorten URL has the same host\n")
	ErrScheme       = errors.New("unsupported scheme\n")
	ErrRateExceded  = errors.New("rate exceeded\n")
)
