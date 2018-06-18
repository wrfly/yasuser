package routes

import (
	"fmt"
	"io"
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/wrfly/yasuser/shortener"
)

func handleIndex(prefix string) gin.HandlerFunc {
	example := fmt.Sprintf("curl %s -d \"%s\"",
		prefix, "http://longlonglong.com/long/long/long?a=1&b=2")

	curlUA := regexp.MustCompile("curl*")

	return func(c *gin.Context) {
		UA := c.Request.UserAgent()
		if matched := curlUA.MatchString(UA); matched {
			// query from curl
			c.String(200, example)
		} else {
			c.Redirect(301, "/index.html")
		}
	}
}

func handleShortURL(s shortener.Shortener) gin.HandlerFunc {
	return func(c *gin.Context) {
		realURL := s.Restore(c.Param("s"))
		if realURL == "" {
			c.String(404, fmt.Sprintln("not found"))
		} else {
			c.Redirect(302, realURL)
		}
	}
}

func handleLongURL(prefix string, s shortener.Shortener) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf := urlBufferPool.Get().([]byte)
		defer urlBufferPool.Put(buf)
		n, err := c.Request.Body.Read(buf)
		if err != io.EOF && err != nil {
			c.String(500, fmt.Sprintf("error: %s\n", err))
			return
		}
		if n > MAX_URL_LENGTH {
			c.String(400, fmt.Sprintln("Bad request, URL too long"))
			return
		}

		longURL := fmt.Sprintf("%s", buf[:n])
		if invalidURL(longURL) {
			c.String(400, fmt.Sprintln("invalid URL"))
			return
		}

		short := s.Shorten(longURL)
		if short == "" {
			c.String(500, "something bad happend")
			return
		}
		shortURL := fmt.Sprintf("%s/%s", prefix, short)
		c.String(200, fmt.Sprintln(shortURL))
	}
}

func invalidURL(URL string) bool {
	u, err := url.Parse(URL)
	if err != nil {
		return true
	}

	switch u.Scheme {
	case "http", "https", "ftp", "tcp":
		return false
	default:
		return true
	}
}
