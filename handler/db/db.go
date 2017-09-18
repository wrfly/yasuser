package db

type Database interface {
	Close() error
	Set(index, shortURL string) error
	Get(index string) (string, error)
	Len() (uint64, error)
}
