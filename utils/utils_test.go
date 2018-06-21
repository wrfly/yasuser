package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5(t *testing.T) {
	c := MD5("hello")
	r := "5d41402abc4b2a76b9719d911017c592"
	assert.Equal(t, c, r)
}

func TestCalHash(t *testing.T) {
	c := CalHash(65)
	t.Log(c)
}

func TestXXHash(t *testing.T) {
	for index := 0; index < 10; index++ {
		h := XXHash(fmt.Sprintf("index:%d", index))
		t.Log(h)
	}
}
