package shortener

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/yasuser/config"
	"github.com/wrfly/yasuser/shortener/cache"
	"github.com/wrfly/yasuser/shortener/db"
	"github.com/wrfly/yasuser/types"
	"github.com/wrfly/yasuser/utils"
)

type Shortener interface {
	// Shorten a long URL
	Shorten(longURL string, ttl time.Duration) (shortURL string)
	// Restore a short URL
	Restore(shortURL string) (longURL string)
	// Shorten the URL with a custom short URL
	ShortenWithCustomURL(customURL, longURL string, ttl time.Duration) error
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
func (stner db_Shortener) Shorten(longURL string, ttl time.Duration) string {
	// xxhash is faster than md5sum
	hashSum := utils.XXHash(longURL)

	// cache first
	if ttl < 0 {
		if shortURL, err := stner.cacher.Get(hashSum); err == nil {
			return shortURL
		}
	}

	// then the db
	shortURL := stner.getShortFromHash(hashSum, longURL, ttl)
	if shortURL == "" {
		return "" // error
	}
	stner.cacher.SetWithExpire(hashSum, shortURL, int(ttl.Seconds()))
	logrus.Debugf("shorten URL: [ %s ] -> [ %s ]; ttl: %s", longURL, shortURL, ttl)

	return shortURL
}

// getShortFromHash return the shortURL from db if found
// otherwise create a new one
func (stner db_Shortener) getShortFromHash(hashSum, longURL string, ttl time.Duration) string {
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
	stner.db.StoreWithTTL(hashSum, shortURL, longURL, ttl)
	return shortURL
}

// Restore a short URL
func (stner db_Shortener) Restore(shortURL string) string {
	// cache first
	if longURL, err := stner.cacher.Get(shortURL); err == nil {
		logrus.Debugf("cache get %s=%s", shortURL, longURL)
		return longURL
	}

	longURL, err := stner.db.GetLong(shortURL)
	if err != nil {
		if err != types.ErrNotFound {
			logrus.Errorf("restore URL error: %s", err)
		}
		return ""
	}

	// if err := stner.cacher.Set(shortURL, longURL); err != nil {
	// 	logrus.Errorf("cache shortURL error: %s", err)
	// }

	logrus.Debugf("db get %s=%s", shortURL, longURL)
	return longURL
}

// Shorten a long URL with a custom key
func (stner db_Shortener) ShortenWithCustomURL(customURL, longURL string, ttl time.Duration) error {
	// xxhash is faster than md5sum
	hashSum := utils.XXHash(longURL)

	// check if the longURL exist and compare its shortenURL with customURL
	if shortURL, err := stner.cacher.Get(hashSum); err == nil {
		if shortURL == customURL {
			return nil
		}
	}
	// cache has this key but the custom url not equal the shorten one
	// check the key (customURL) already exist
	if _, err := stner.cacher.Get(customURL); err == nil {
		return types.ErrAlreadyExist
	}
	if _, err := stner.db.GetLong(customURL); err == nil {
		return types.ErrAlreadyExist
	}

	existShortURL := stner.getShortFromHash(hashSum, longURL, ttl)
	if existShortURL == "" {
		return fmt.Errorf("error") // error
	}

	// reset the custom URL
	stner.db.StoreWithTTL(hashSum, customURL, longURL, ttl)
	stner.cacher.SetWithExpire(hashSum, customURL, int(ttl.Seconds()))
	stner.cacher.SetWithExpire(customURL, longURL, int(ttl.Seconds()))

	logrus.Debugf("shorten URL: [ %s ] -> [ %s ]", longURL, customURL)

	return nil
}
