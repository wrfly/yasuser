package filter

import "errors"

// errors
var (
	ErrBadDomain  = errors.New("domain in blacklist")
	ErrBadKeyword = errors.New("url contains bad keyword")
	ErrRedirect   = errors.New("url is a redirect")
)
