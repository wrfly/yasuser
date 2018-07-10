package shortener

import (
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
	Shorten(long string, ops *types.ShortOptions) (*types.URL, error)
	// Restore a short URL
	Restore(short string) (*types.URL, error)
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
func (stner db_Shortener) Shorten(long string, opts *types.ShortOptions) (*types.URL, error) {
	// check custom URL
	if opts == nil {
		opts = &types.ShortOptions{}
	} else if stner.customURLAlreadyExist(opts.Custom, long) {
		return nil, types.ErrAlreadyExist
	}

	// xxhash is faster than md5sum
	hashSum := utils.XXHash(long)

	// get from mem-cache
	if cacheURL, err := stner.cacher.Get(hashSum); err == nil {
		if cacheURL.Expired() {
			goto SET_URL
		}

		if !cacheURL.Expired() && opts.Custom == cacheURL.Custom {
			return cacheURL, nil
		}
		// cache expired or
		if opts.Custom == cacheURL.Custom {
			logrus.Debugf("cache get %s=%s", long, cacheURL)
			return cacheURL, nil
		} else {
			return nil, types.ErrAlreadyExist
		}
	}

	// cache not found, check database
	if shortURL, err := stner.db.GetShort(hashSum); err == nil {
		if !shortURL.Expired() && opts.Custom == shortURL.Custom {
			return shortURL, nil
		}
		logrus.Debugf("url expired or new customURL, re-create it")
	} else if err != types.ErrNotFound {
		logrus.Errorf("get shortURL from db error: %s", err)
		return nil, err
	}
	logrus.Debugf("url %s not found, create a new one", long)

SET_URL:
	short := strings.TrimLeft(utils.CalHash(stner.db.Len()), "0")
	logrus.Debugf("short %s to %s", long, short)

	u := newURL(hashSum, short, long, opts)
	stner.db.Store(u)

	stner.cacher.Store(u)
	logrus.Debugf("shorten URL: [ %s ] -> [ %s ]; opts: %v", long, short, opts)

	return u, nil
}

func newURL(hashSum, short, long string, opts *types.ShortOptions) *types.URL {
	u := new(types.URL)
	u.Hash = hashSum
	u.Short = short
	u.Ori = long

	u.Custom = opts.Custom
	if opts.TTL.Seconds() > 0 {
		exp := time.Now().Add(opts.TTL)
		u.Expire = &exp
	}
	return u
}

// Restore a short URL
func (stner db_Shortener) Restore(short string) (*types.URL, error) {
	// found in cache
	if cacheURL, err := stner.cacher.Get(short); err == nil {
		logrus.Debugf("cache get %s=%s", short, cacheURL)
		if cacheURL.Expired() {
			return nil, types.ErrURLExpired
		}
		return cacheURL, nil
	}

	longURL, err := stner.db.GetLong(short)
	if err != nil {
		if err == types.ErrNotFound {
			return nil, err
		}
		logrus.Errorf("restore URL error: %s", err)
	}

	if longURL.Expired() {
		return nil, types.ErrURLExpired
	}

	stner.cacher.Store(longURL)
	logrus.Debugf("restore url [ %s ] -> [ %s ]", short, longURL)

	return longURL, nil
}

func (stner db_Shortener) customURLAlreadyExist(short, long string) bool {
	if short == "" {
		return false
	}
	if url, err := stner.Restore(short); err == nil {
		if url.Expired() {
			return false
		}
		if url.Ori != long {
			return true
		}
	}

	return false
}
