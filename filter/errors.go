package filter

import "errors"

// errors
var (
	ErrInBlackList = errors.New("domain in blacklist")
	ErrBadKeyword  = errors.New("url containes bad keyword")
	ErrRemoved     = errors.New("domain removed")
)
