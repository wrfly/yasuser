package handler

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/short-url/handler/db"
	"github.com/wrfly/short-url/utils"
)

type Shorter struct {
	DB db.Database
}

// Short long2short
func (s *Shorter) Short(URL string) string {
	index := utils.MD5(URL)

	// mem cache first,then db
	short, err := s.DB.GetShort(index)
	if err != nil {
		logrus.Errorf("get short from db error: %s", err)
		return ""
	}
	if short != "" {
		return short
	}

	// not found
	short = s.sURL(URL)
	go func() {
		err = s.DB.SetLong(short, URL)
		if err != nil {
			logrus.Errorf("set long error: %s", err)
		}
		err = s.DB.SetShort(index, short)
		if err != nil {
			logrus.Errorf("set short error: %s", err)
		}
	}()
	return short
}

func (s *Shorter) sURL(URL string) string {
	n, err := s.DB.Len()
	if err != nil {
		logrus.Errorf("get db lenth error: %s", err)
		return "_"
	}

	shortURL := utils.CalHash(n)
	shortURL = strings.TrimLeft(shortURL, "0")

	return shortURL
}

func (s *Shorter) Long(shortURL string) string {
	// mem cache first,then db
	longURL, err := s.DB.GetLong(shortURL)
	if err != nil {
		logrus.Errorf("restore URL error: %s", err)
		return ""
	}

	return longURL
}
