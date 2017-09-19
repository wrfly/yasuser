package handler

import (
	"github.com/Sirupsen/logrus"
	"github.com/wrfly/short-url/handler/db"
	"github.com/wrfly/short-url/utils"
)

type Shorter struct {
	DB db.Database
}

func (s *Shorter) Short(URL string) string {
	index := utils.MD5(URL)

	// mem cache first,then db
	short, err := s.DB.Get(index)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	// not found
	if short == "" {
		short = s.sURL(URL)
	}
	err = s.DB.Set(index, short)
	if err != nil {
		logrus.Error(err)
		return ""
	}
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
