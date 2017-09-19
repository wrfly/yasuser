package utils

import (
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
