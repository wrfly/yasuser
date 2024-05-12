package filter

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/wrfly/yasuser/config"
)

type Filter interface {
	OK(*url.URL) error
}

type list struct {
	blacklist map[string]bool
	whitelist map[string]bool
}

type urlFilter struct {
	domain  list
	keyword list
}

var simpleHTTPClient = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
		ForceAttemptHTTP2: false,
	},
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func (f *urlFilter) OK(u *url.URL) error {
	// bypass domain
	if f.domain.whitelist[u.Hostname()] {
		return nil
	}

	// bypass keyword
	for x := range f.keyword.whitelist {
		if strings.Contains(u.Path, x) {
			return nil
		}
	}

	// bad domain
	if f.domain.blacklist[u.Hostname()] {
		return ErrBadDomain
	}

	// bad keyword
	for x := range f.keyword.blacklist {
		if strings.Contains(u.Path, x) {
			return ErrBadKeyword
		}
	}

	// 3xx urls
	resp, err := simpleHTTPClient.Head(u.String())
	if err != nil {
		return err
	}
	if resp.StatusCode/100 == 3 {
		return ErrRedirect
	}

	return nil
}

func makeList(slice []string) map[string]bool {
	x := make(map[string]bool, len(slice))
	for _, s := range slice {
		x[s] = true
	}
	return x
}

// New filter
func New(conf config.Filter) Filter {
	return &urlFilter{
		domain: list{
			whitelist: makeList(conf.Domain.WhiteList),
			blacklist: makeList(conf.Domain.BlackList),
		},
		keyword: list{
			whitelist: makeList(conf.Keyword.WhiteList),
			blacklist: makeList(conf.Keyword.BlackList),
		},
	}
}
