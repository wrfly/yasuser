package routes

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/testing-kit/util/tokenbucket"

	"github.com/wrfly/yasuser/bg"
	"github.com/wrfly/yasuser/config"
	"github.com/wrfly/yasuser/filter"
	"github.com/wrfly/yasuser/routes/asset"
	s "github.com/wrfly/yasuser/shortener"
	"github.com/wrfly/yasuser/types"
)

const (
	maxPasswordLength = 60
	maxCustomLength   = 60
)

var (
	validCustomURI *regexp.Regexp
	indexTemplate  *template.Template
)

func init() {
	validURI, err := regexp.Compile("^[a-zA-Z0-9][a-zA-Z0-9_+-]+$")
	if err != nil {
		panic(err)
	}
	validCustomURI = validURI

	a, err := asset.Find("/index.html")
	if err != nil {
		panic(err)
	}
	indexTemplate = a.Template()

}

type server struct {
	domain string
	gaID   string
	limit  int64

	handler    s.Shortener
	assetFiles map[string]bool
	filter     filter.Filter

	host string
	tb   map[string]tokenbucket.Bucket
}

func newServer(conf config.SrvConfig,
	shortener s.Shortener, filter filter.Filter) server {

	u, err := url.Parse(conf.Domain)
	if err != nil {
		panic(err)
	}

	srv := server{
		domain:     conf.Domain,
		host:       u.Host,
		handler:    shortener,
		gaID:       conf.GAID,
		limit:      conf.Limit,
		assetFiles: make(map[string]bool),
		tb:         make(map[string]tokenbucket.Bucket, 0),
		filter:     filter,
	}
	for _, a := range asset.List() {
		srv.assetFiles[a.Name()] = true
	}

	return srv
}

func (s *server) handleIndex() gin.HandlerFunc {
	curlUA := regexp.MustCompile("curl*")

	return func(c *gin.Context) {
		if matched := curlUA.MatchString(c.Request.UserAgent()); matched {
			// query from curl
			c.String(200, fmt.Sprintf("curl %s -d \"%s\"",
				s.domain, "http://longlonglong.com/long/long/long?a=1&b=2"))
		} else {
			shortened, requests := s.handler.Status()
			// visit from a web browser
			indexTemplate.Execute(c.Writer, map[string]interface{}{
				"domain":    s.domain,
				"gaID":      s.gaID,
				"shortened": shortened,
				"requests":  requests,
				"bglink":    bg.Image(),
			})
		}
	}
}

func (s *server) restoreShortLink() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Debugf("request code: %s", c.Request.RequestURI)

		code := c.Param("code")
		if code == "" {
			c.String(404, "not found")
			return
		}
		if s.assetFiles[c.Request.RequestURI] {
			asset.Handler(c.Writer, c.Request)
			return
		}

		shortURL, err := s.handler.Restore(code)
		if err != nil {
			if err == types.ErrNotFound {
				c.String(http.StatusNotFound, err.Error())
			} else {
				c.String(http.StatusInternalServerError, err.Error())
			}
			return
		}

		if shortURL.Expired() {
			c.String(http.StatusNotFound, "URL expired")
			return
		}

		if err := s.invalidURL(shortURL.Original); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			c.Redirect(http.StatusPermanentRedirect, shortURL.Original)
		}
	}
}

func (s *server) shortenURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		// rate limit
		IP := c.ClientIP()
		if tb, ok := s.tb[IP]; !ok {
			s.tb[IP] = tokenbucket.New(s.limit, time.Second)
		} else {
			if !tb.TakeOne() {
				badRequest(c, fmt.Errorf("rate exceeded"))
				return
			}
		}

		buf := urlBufferPool.Get().([]byte)
		defer urlBufferPool.Put(buf)
		n, err := c.Request.Body.Read(buf)
		if err != io.EOF && err != nil {
			badRequest(c, err)
			return
		}
		if n > _MaxURLLength {
			badRequest(c, types.ErrURLTooLong)
			return
		}

		longURL := fmt.Sprintf("%s", buf[:n])
		if err := s.invalidURL(longURL); err != nil {
			badRequest(c, err)
			return
		}

		opts, err := generateOptions(c.Request.Header)
		if err != nil {
			badRequest(c, err)
			return
		}
		shortURL, err := s.handler.Shorten(longURL, opts)
		if err != nil {
			badRequest(c, err)
			return
		}

		c.String(200, fmt.Sprintf("%s/%s\n", s.domain, shortURL.Short))
	}
}

func badRequest(c *gin.Context, err error) {
	c.String(http.StatusBadRequest, fmt.Sprintln(err.Error()))
}

func generateOptions(h http.Header) (*types.ShortOptions, error) {
	var (
		duration time.Duration = -1
		err      error
	)

	customURI := h.Get("CUSTOM")
	passWord := h.Get("PASS")
	ttl := h.Get("TTL")
	if ttl != "" {
		duration, err = time.ParseDuration(ttl)
		if err != nil {
			return nil, err
		}
	}

	if len(passWord) > maxPasswordLength {
		return nil, fmt.Errorf("password length exceeded, max %d",
			maxPasswordLength)
	}

	if customURI != "" {
		if len(customURI) > maxCustomLength {
			return nil, fmt.Errorf("custom URI length exceeded, max %d",
				maxCustomLength)
		}
		if !validCustomURI.MatchString(customURI) {
			return nil, fmt.Errorf("invalid custom URI, must match %s",
				validCustomURI.String())
		}
		// TODO: return error when the short
		// URL is the same as some assset files' name
	}

	return &types.ShortOptions{
		Custom: customURI,
		TTL:    duration,
		Pass:   passWord,
	}, nil
}

func (s *server) invalidURL(URL string) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	if u.Hostname() == s.host {
		return types.ErrSameHost
	}
	switch u.Scheme {
	case "http", "https", "ftp", "tcp":
	default:
		return types.ErrScheme
	}

	return s.filter.OK(u)
}
