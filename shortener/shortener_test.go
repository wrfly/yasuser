package shortener

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wrfly/yasuser/config"
)

func TestShorter(t *testing.T) {
	c := config.ShortenerConfig{
		Store: config.StoreConfig{
			DBPath: "/tmp/test.yasuser.bolt.db",
			DBType: "bolt",
		},
	}
	s := New(c)

	for i := 0; i < 10; i++ {
		URL := fmt.Sprintf("http://%v", i)
		short := s.Shorten(URL)
		long := s.Restore(short)
		assert.Equal(t, URL, long)
	}
}
