package types

import (
	"bytes"
	"encoding/gob"
	"time"
)

type ShortOptions struct {
	Pass   string
	Custom string
	TTL    time.Duration
}

type URL struct {
	Original string
	Password string
	Short    string
	Hash     string
	Expire   *time.Time
	bs       []byte
}

func (u *URL) Bytes() []byte {
	if u.bs == nil {
		var data bytes.Buffer
		gob.NewEncoder(&data).Encode(u)
		u.bs = data.Bytes()
	}
	return u.bs
}

func (u *URL) Decode(byts []byte) *URL {
	var data bytes.Buffer
	data.Write(byts)
	gob.NewDecoder(&data).Decode(u)
	return u
}

func (u *URL) HashSum() []byte {
	return []byte(u.Hash)
}

func (u *URL) ShortURL() []byte {
	return []byte(u.Short)
}

func (u *URL) Expired() bool {
	if u.Expire == nil {
		return false
	}
	if u.Expire.Sub(time.Now()).Nanoseconds() > 0 {
		return false
	}
	return true
}
