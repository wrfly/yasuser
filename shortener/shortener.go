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
	// get status (total handled and visited)
	Status() (int64, int64)
}

type shortener struct {
	db     db.Database
	cacher cache.Cacher
}

func New(conf config.StoreConfig) Shortener {
	db, err := db.New(conf)
	if err != nil {
		logrus.Fatalf("Create Shortener error: %v", err)
	}

	return shortener{
		db:     db,
		cacher: cache.NewCacher(100 * 1024 * 1024),
	}
}

// Shorten a long URL
func (s shortener) Shorten(long string, opts *types.ShortOptions) (*types.URL, error) {
	// check custom URL
	if opts == nil {
		opts = &types.ShortOptions{}
	} else if s.customURLAlreadyExist(opts.Custom, long) {
		return nil, types.ErrAlreadyExist
	}

	// a long URL and its opts made up the hashSum
	hashSum := utils.HashURL(long, opts)

	// get from mem-cache
	if cacheURL, err := s.cacher.Get(hashSum); err == nil {
		logrus.Debugf("cache get %s=%s", long, cacheURL.Short)
		return cacheURL, nil
	}

	// cache not found, check database
	if shortURL, err := s.db.GetShort(hashSum); err == nil {
		if !shortURL.Expired() {
			s.cacher.Store(shortURL)
			return shortURL, nil
		}
		logrus.Debugf("url expired or new customURL, re-create it")
	} else if err != types.ErrNotFound {
		logrus.Errorf("get shortURL from db error: %s", err)
		return nil, err
	}
	logrus.Debugf("url %s not found, create a new one", long)

	short := opts.Custom
	keyNum, err := s.db.IncKey()
	if err != nil {
		return nil, err
	}
	if short == "" {
		short = strings.TrimLeft(utils.CalHash(keyNum), "0")
	}
	logrus.Debugf("short %s to %s", long, short)

	u := newURL(hashSum, short, long, opts)
	if err := s.db.Store(u); err != nil {
		return nil, err
	}

	s.cacher.Store(u)
	logrus.Debugf("shorten URL: [ %s ] -> [ %s ]; opts: %v", long, short, opts)

	return u, nil
}

func newURL(hashSum, short, long string, opts *types.ShortOptions) *types.URL {
	u := new(types.URL)
	u.Hash = hashSum
	u.Short = short
	u.Ori = long

	if opts.TTL.Seconds() > 0 {
		exp := time.Now().Add(opts.TTL)
		u.Expire = &exp
	}
	return u
}

// Restore a short URL
func (s shortener) Restore(short string) (*types.URL, error) {
	// found in cache
	if cacheURL, err := s.cacher.Get(short); err == nil {
		logrus.Debugf("cache get %s=%s", short, cacheURL)
		if cacheURL.Expired() {
			return nil, types.ErrURLExpired
		}
		s.db.IncVisited()
		return cacheURL, nil
	}

	URL, err := s.db.GetLong(short)
	if err != nil {
		if err == types.ErrNotFound {
			return nil, err
		}
		logrus.Errorf("restore URL error: %s", err)
	}

	if URL.Expired() {
		return nil, types.ErrURLExpired
	}

	s.cacher.Store(URL)
	logrus.Debugf("restore url [ %s ] -> [ %s ]", short, URL.Ori)
	s.db.IncVisited()

	return URL, nil
}

func (s shortener) customURLAlreadyExist(short, long string) bool {
	if short == "" {
		return false
	}
	if url, err := s.Restore(short); err == nil {
		if url.Expired() {
			return false
		}
		if url.Ori != long {
			return true
		}
	}

	return false
}

func (s shortener) Status() (int64, int64) {
	k, _ := s.db.Keys()
	v, _ := s.db.Visited()
	return k, v
}
