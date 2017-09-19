package handler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wrfly/short-url/handler/db"
)

func TestShort(t *testing.T) {
	db, err := db.NewDB("/tmp/myxx.db")
	assert.NoError(t, err)
	s := Shorter{
		DB: db,
	}

	p := fmt.Println

	for i := 0; i < 64; i++ {
		u := fmt.Sprintf("%v", i)
		p(s.Short(u))
	}

}
