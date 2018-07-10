package types

import (
	"bytes"
	"encoding/gob"
	"time"
)

type URL struct {
	Ori    string
	Passwd string
	Short  string
	Hash   string
	Expire time.Time
}

func (u *URL) Bytes() []byte {
	var data bytes.Buffer
	gob.NewEncoder(&data).Encode(u)
	return data.Bytes()
}

func (u *URL) Decode(byts []byte) {
	var data bytes.Buffer
	data.Write(byts)
	gob.NewDecoder(&data).Decode(u)
}

func (u *URL) HashSum() []byte {
	return []byte(u.Hash)
}

func (u *URL) ShortURL() []byte {
	return []byte(u.Short)
}
