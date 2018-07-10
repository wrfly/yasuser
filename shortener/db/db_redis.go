package db

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/wrfly/yasuser/types"
)

type redisDB struct {
	cli *redis.Client
}

func newRedisDB(redisAddr string) (*redisDB, error) {
	opts, err := redis.ParseURL(redisAddr)
	if err != nil {
		return nil, err
	}
	cli := redis.NewClient(opts)

	if err := cli.Ping().Err(); err != nil {
		return nil, err
	}

	return &redisDB{
		cli: cli,
	}, nil
}

func (r *redisDB) Close() error {
	return r.cli.Close()
}

func (r *redisDB) Len() int64 {
	length, err := r.cli.Incr("KEY_NUMS").Result()
	if err != nil {
		panic(err)
	}
	length /= 2
	length += skipKeyNums

	return length
}

func (r *redisDB) Store(hashSum, shortURL, longURL string) error {
	if err := r.set(hashSum, shortURL); err != nil {
		return err
	}
	// TODO: what if set failed, need to rollback?
	return r.set(shortURL, longURL)
}

func (r *redisDB) GetShort(hashSum string) (short string, err error) {
	return r.get(hashSum)
}

func (r *redisDB) GetLong(shortURL string) (long string, err error) {
	return r.get(shortURL)
}

func (r *redisDB) set(key, value string) error {
	return r.cli.Set(key, value, -1).Err()
}

func (r *redisDB) get(key string) (value string, err error) {
	stringCmd := r.cli.Get(key)
	if err := stringCmd.Err(); err != nil {
		if err == redis.Nil {
			return "", types.ErrNotFound
		}
		return "", err
	}
	return stringCmd.String(), nil
}

func (r *redisDB) StoreWithTTL(hashSum, shortURL, longURL string, ttl time.Duration) error {
	if err := r.cli.Set(hashSum, shortURL, ttl).Err(); err != nil {
		return err
	}
	// TODO: what if set failed, need to rollback?
	return r.cli.Set(shortURL, longURL, ttl).Err()
}
