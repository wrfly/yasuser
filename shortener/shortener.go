package shortener

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/yasuser/config"
	"github.com/wrfly/yasuser/shortener/cache"
	"github.com/wrfly/yasuser/shortener/db"
	"github.com/wrfly/yasuser/types"
	"github.com/wrfly/yasuser/utils"
)

type Shortener interface {
	// Shorten a long URL
	Shorten(longURL string) (shortURL string)
	// Restore a short URL
	Restore(shortURL string) (longURL string)
}

type db_Shortener struct {
	db     db.Database
	cacher cache.Cacher
}

func New(conf config.ShortenerConfig) Shortener {
	db, err := db.New(conf.Store)
	if err != nil {
		logrus.Fatalf("Create Shortener error: %v", err)
	}

	return db_Shortener{
		db:     db,
		cacher: cache.NewCacher(100 * 1024 * 1024),
	}
}

// Shorten a long URL
func (stner db_Shortener) Shorten(longURL string) string {
	// xxhash is faster than md5sum
	hashSum := utils.XXHash(longURL)

	// cache first
	if shortURL, err := stner.cacher.Get(hashSum); err == nil {
		return shortURL
	}

	// then the db
	shortURL := stner.getShortFromHash(hashSum, longURL)
	if shortURL != "" {
		stner.cacher.Set(hashSum, shortURL)
	}
	logrus.Debugf("shorten URL: [ %s ] -> [ %s ]", longURL, shortURL)

	return shortURL
}

func (stner db_Shortener) getShortFromHash(hashSum, longURL string) string {
	// return from db if found, otherwise create a new one
	shortURL, err := stner.db.GetShort(hashSum)
	if err == nil {
		return shortURL
	}
	if err != types.ErrNotFound {
		logrus.Errorf("get shortURL from db error: %s", err)
		return ""
	}
	logrus.Debugf("url %s not found, create a new one", longURL)

	shortURL = strings.TrimLeft(utils.CalHash(stner.db.Len()), "0")

	stner.store(shortURL, hashSum, longURL)

	return shortURL
}

// Restore a short URL
func (stner db_Shortener) Restore(shortURL string) string {
	// cache first
	if longURL, err := stner.cacher.Get(shortURL); err == nil {
		return longURL
	}

	longURL, err := stner.db.GetLong(shortURL)
	if err != nil {
		if err != types.ErrNotFound {
			logrus.Errorf("restore URL error: %s", err)
		}
		return ""
	}

	if err := stner.cacher.Set(shortURL, longURL); err != nil {
		logrus.Errorf("cache shortURL error: %s", err)
	}

	return longURL
}

// Restore a short URL
func (stner db_Shortener) store(shortURL, hashSum, longURL string) {
	if err := stner.db.SetLong(shortURL, longURL); err != nil {
		logrus.Errorf("set long error: %s", err)
	}

	if err := stner.db.SetShort(hashSum, shortURL); err != nil {
		logrus.Errorf("set shortURL error: %s", err)
	}
}
