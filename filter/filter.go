package filter

type Filter interface {
	Removed(domain string) bool
	InBlackList(domain string) bool
	InWhiteList(domain string) bool
	BadKeyword(url string) bool
}
