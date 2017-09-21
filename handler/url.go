package handler

import (
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
		logrus.Error(err)
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
			logrus.Error(err)
		}
		err = s.DB.SetShort(index, short)
		if err != nil {
			logrus.Error(err)
		}
	}()
	return short
}

func (s *Shorter) sURL(URL string) string {
	n, err := s.DB.Len()
	if err != nil {
		logrus.Error(err)
		return "_"
	}

	shortURL := utils.CalHash(n)

	return shortURL
}

func (s *Shorter) Long(shortURL string) string {
	// mem cache first,then db
	longURL, err := s.DB.GetLong(shortURL)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	return longURL
}
