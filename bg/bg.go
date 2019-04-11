package bg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const bingDailyURL = "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN"

type bingDaily struct {
	Images []struct {
		URL     string `json:"url"`
		URLBase string `json:"urlbase"`
	} `json:"images"`
}

var (
	link string
	last time.Time
)

func Image() string {
	if time.Now().Sub(last) < time.Minute*10 {
		return link
	}

	resp, err := http.Get(bingDailyURL)
	if err != nil {
		logrus.Errorf("get image link error: %s", err)
		return ""
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("read body error: %s", err)
		return ""
	}
	r := new(bingDaily)
	if err := json.Unmarshal(bs, r); err != nil {
		logrus.Errorf("unmarshal json error: %s", err)
		return ""
	}

	if len(r.Images) >= 1 {
		link = "https://www.bing.com" + r.Images[0].URL
		last = time.Now()
		return link
	}
	logrus.Error("no images")
	return ""
}
