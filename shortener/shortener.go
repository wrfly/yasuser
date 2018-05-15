package shortener

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/yasuser/config"
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
	db db.Database
	// TODO: go-memdb for caching
	// cacher
}

func New(conf config.ShortenerConfig) Shortener {
	db, err := db.New(conf.Store)
	if err != nil {
		logrus.Fatalf("Create Shortener error: %v", err)
	}

	return db_Shortener{
		db: db,
	}
}

// Shorten a long URL
func (stner db_Shortener) Shorten(longURL string) string {
	md5sum := utils.MD5(longURL)

	// return from db if found, otherwise create a new one
	shortURL, err := stner.db.GetShort(md5sum)
	if err == nil {
		return shortURL
	}
	if err != types.ErrNotFound {
		logrus.Errorf("get shortURL from db error: %s", err)
		return ""
	}
	logrus.Debugf("url %s not found, create a new one", longURL)

	shortURL = stner.shorten(longURL)
	go func() {
		if err = stner.db.SetLong(shortURL, longURL); err != nil {
			logrus.Errorf("set long error: %s", err)
		}

		if err = stner.db.SetShort(md5sum, shortURL); err != nil {
			logrus.Errorf("set shortURL error: %s", err)
		}
	}()

	return shortURL
}

func (stner db_Shortener) shorten(longURL string) string {
	n, err := stner.db.Len()
	if err != nil {
		logrus.Errorf("get db lenth error: %s", err)
		return "_"
	}

	shortURL := utils.CalHash(int(n))
	shortURL = strings.TrimLeft(shortURL, "0")

	return shortURL
}

// Restore a short URL
func (stner db_Shortener) Restore(shortURL string) string {
	longURL, err := stner.db.GetLong(shortURL)
	if err != nil {
		logrus.Errorf("restore URL error: %s", err)
		return ""
	}

	return longURL
}
