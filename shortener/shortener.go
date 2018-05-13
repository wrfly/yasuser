package shortener

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/short-url/config"
	"github.com/wrfly/short-url/shortener/db"
	"github.com/wrfly/short-url/utils"
)

type Shortener interface {
	// Shorten a long URL
	Shorten(longURL string) (shortURL string)
	// Restore a short URL
	Restore(shortURL string) (longURL string)
}

type db_shortener struct {
	db db.Database
	// TODO: go-memdb for caching
	// cacher
}

func New(conf config.ShortenerConfig) Shortener {
	db, err := db.New(conf.Store)
	if err != nil {
		logrus.Fatalf("Create Shortener error: %v", err)
	}

	return db_shortener{
		db: db,
	}
}

// Shorten a long URL
func (s db_shortener) Shorten(longURL string) string {
	index := utils.MD5(longURL)

	// return from db if found, otherwise create a new one
	shortURL, err := s.db.GetShort(index)
	if err != nil {
		logrus.Errorf("get shortURL from db error: %s", err)
		return ""
	}
	if shortURL != "" {
		return shortURL
	}

	// didn't find it, create a new one and store
	shortURL = s.shorten(longURL)
	go func() {
		if err = s.db.SetLong(shortURL, longURL); err != nil {
			logrus.Errorf("set long error: %s", err)
		}

		if err = s.db.SetShort(index, shortURL); err != nil {
			logrus.Errorf("set shortURL error: %s", err)
		}
	}()

	return shortURL
}

func (s db_shortener) shorten(longURL string) string {
	n, err := s.db.Len()
	if err != nil {
		logrus.Errorf("get db lenth error: %s", err)
		return "_"
	}

	shortURL := utils.CalHash(n)
	shortURL = strings.TrimLeft(shortURL, "0")

	return shortURL
}

// Restore a short URL
func (s db_shortener) Restore(shortURL string) string {
	longURL, err := s.db.GetLong(shortURL)
	if err != nil {
		logrus.Errorf("restore URL error: %s", err)
		return ""
	}

	return longURL
}
