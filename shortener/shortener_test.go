package shortener

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wrfly/short-url/shortener/db"
)

func TestShorter(t *testing.T) {
	db, err := db.NewDB("/tmp/my1x.db")
	assert.NoError(t, err)
	s := Shorter{
		DB: db,
	}

	for i := 0; i < 64; i++ {
		u := fmt.Sprintf("%v", i)
		short := s.Short(u)
		long := s.Long(short)
		assert.Equal(t, u, long)
	}
}
