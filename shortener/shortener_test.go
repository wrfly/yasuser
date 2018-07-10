package shortener

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wrfly/yasuser/config"
	"github.com/wrfly/yasuser/types"
)

func TestShorter(t *testing.T) {
	tempDBPath := "/tmp/test.yasuser.bolt.db"
	os.Remove(tempDBPath)
	defer os.Remove(tempDBPath)

	c := config.ShortenerConfig{
		Store: config.StoreConfig{
			DBPath: tempDBPath,
			DBType: "bolt",
		},
	}
	s := New(c)

	t.Run("normal", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			URL := fmt.Sprintf("http://normal=%v", i)
			short, err := s.Shorten(URL, nil)
			if err != nil {
				t.Error(err)
			}
			long, _ := s.Restore(short.Short)
			assert.Equal(t, URL, long.Ori)
		}
	})

	t.Run("with ttl", func(t *testing.T) {
		shortURLs := make([]string, 10)
		for i := 0; i < 10; i++ {
			URL := fmt.Sprintf("http://with.ttl=%v", i)
			short, err := s.Shorten(URL, &types.ShortOptions{
				TTL: time.Second,
			})
			if err != nil {
				t.Error(err)
			}
			long, _ := s.Restore(short.Short)
			assert.Equal(t, URL, long.Ori)
			shortURLs[i] = short.Short
		}

		time.Sleep(time.Second)

		for _, sURL := range shortURLs {
			_, err := s.Restore(sURL)
			assert.Error(t, err, "expired")
		}
	})

	t.Run("with custom", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			URL := fmt.Sprintf("http://custom=%v", i)
			custom := fmt.Sprintf("%v%v", i, i)
			short, err := s.Shorten(URL, &types.ShortOptions{
				Custom: custom,
			})
			if err != nil {
				t.Error(err)
			}
			long, _ := s.Restore(short.Short)
			assert.Equal(t, custom, long.Custom)
		}
	})

	t.Run("custom already exist error", func(t *testing.T) {
		_, err := s.Shorten("123", &types.ShortOptions{
			Custom: "custom",
		})
		assert.NoError(t, err)

		_, err = s.Shorten("456", &types.ShortOptions{
			Custom: "custom",
		})
		assert.Error(t, err)

	})

	t.Run("custom expired and rewrite it", func(t *testing.T) {
		_, err := s.Shorten("234", &types.ShortOptions{
			Custom: "hello",
			TTL:    time.Second,
		})
		assert.NoError(t, err)
		time.Sleep(time.Second * 1)

		_, err = s.Shorten("567", &types.ShortOptions{
			Custom: "hello",
		})
		assert.NoError(t, err)
	})
}
