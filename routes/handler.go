package routes

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"regexp"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"github.com/wrfly/yasuser/shortener"
)

type server struct {
	domain        string
	scheme        string
	stener        shortener.Shortener
	indexTemplate *template.Template
	fileMap       map[string]bool
}

func (s *server) init() {
	// init
	fileMap := map[string]bool{}
	for _, fileName := range AssetNames() {
		fileMap[fileName] = true
	}
	s.fileMap = fileMap

	bs, err := Asset("index.html")
	if err != nil {
		panic(err)
	}
	s.indexTemplate, err = template.New("index").Parse(string(bs))
	if err != nil {
		panic(err)
	}
}

func (s *server) handleIndex() gin.HandlerFunc {
	curlUA := regexp.MustCompile("curl*")

	return func(c *gin.Context) {
		UA := c.Request.UserAgent()
		if s.scheme == "" {
			if c.Request.URL.Scheme == "" {
				s.scheme = "http"
			} else {
				s.scheme = "https"
			}
		}

		if matched := curlUA.MatchString(UA); matched {
			// query from curl
			c.String(200, fmt.Sprintf("curl %s://%s -d \"%s\"", s.scheme,
				s.domain, "http://longlonglong.com/long/long/long?a=1&b=2"))
		} else {
			// visit from a web browser
			s.indexTemplate.Execute(c.Writer, map[string]string{
				"domain": s.scheme + "://" + s.domain,
			})
		}
	}
}

func (s *server) handleURI() gin.HandlerFunc {

	return func(c *gin.Context) {
		URI := c.Param("URI")

		switch {
		case URI == "":
			c.String(404, fmt.Sprintln("not found"))

		case s.fileMap[URI]:
			// handle static files
			http.FileServer(&assetfs.AssetFS{
				Asset:     Asset,
				AssetDir:  AssetDir,
				AssetInfo: AssetInfo,
				Prefix:    "/",
			}).ServeHTTP(c.Writer, c.Request)

		default:
			// handle shortURL
			if realURL := s.stener.Restore(URI); realURL == "" {
				c.String(404, fmt.Sprintln("not found"))
			} else {
				c.Redirect(302, realURL)
			}
		}
	}
}

func (s *server) handleLongURL() gin.HandlerFunc {
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

		short := s.stener.Shorten(longURL)
		if short == "" {
			c.String(500, "something bad happend")
			return
		}
		shortURL := fmt.Sprintf("%s://%s/%s", s.scheme, s.domain, short)
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
