package db

type Database interface {
	Close() error
	SetShort(index, shortURL string) error
	GetShort(index string) (string, error)
	SetLong(index, longURL string) error
	GetLong(index string) (string, error)
	Len() (int, error)
}
