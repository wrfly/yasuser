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
			// normal web browser
			// TODO: a pretty index web page
			welcome := fmt.Sprintf("Welcome [%s], you'll see a pretty index page later...",
				UA)
			c.String(200, welcome)
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
	case "":
		return true
	case "http":
	case "https":
	case "ftp":
	case "tcp":
	default:
		return true
	}

	return false
}
