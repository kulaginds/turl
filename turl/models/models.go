package models

import (
	"net/url"
)

const (
	NoErrors uint = 0
	ErrNotValidUrl uint = 1
	ErrEmptyUrl uint = 2
	ErrUrlTooLong uint = 3
)

const (
	EmptyStr = ""
	MaxUrlLength = 2000
)

type ShortUrl struct {
	Id  int
	Url string
}

type LongUrl struct {
	Url string
}

var abc *ABC
var config *Config

func Initialize(conf *Config) {
	config = conf
	abc = NewABC(conf.ABC(), conf.ABCIdMinLen())
}

func (u *ShortUrl) Validate() (errorCode uint, ok bool) {
	var length = uint(len(u.Url))

	if 0 == length {
		errorCode, ok = ErrEmptyUrl, false
		return
	}

	if length > MaxUrlLength {
		errorCode, ok = ErrUrlTooLong, false
		return
	}

	shortUrl, err := url.ParseRequestURI(u.Url)

	if nil != err {
		errorCode, ok = ErrNotValidUrl, false
	}

	if shortUrl.Host != config.ServiceUrl().Host {
		errorCode, ok = ErrNotValidUrl, false
	}

	errorCode, ok = u.validatePath(shortUrl.Path)

	if !ok {
		return
	}

	return NoErrors, true
}

func (u *ShortUrl) validatePath(path string) (errorCode uint, ok bool) {
	var index int

	errorCode, ok = NoErrors, true

	if "/" == string(path[0]) {
		path = path[1:]
	}

	l := len(path)

	for i := 0; i < l; i++ {
		index = abc.IndexOf(path[i])

		if -1 == index {
			errorCode, ok = ErrNotValidUrl, false
			return
		}
	}

	return
}

func (u *LongUrl) Validate() (errorCode uint, ok bool) {
	var length = uint(len(u.Url))

	if 0 == length {
		errorCode, ok = ErrEmptyUrl, false
		return
	}

	if length > MaxUrlLength {
		errorCode, ok = ErrUrlTooLong, false
		return
	}

	_, err := url.ParseRequestURI(u.Url)

	if nil != err {
		errorCode, ok = ErrNotValidUrl, false
	}

	return NoErrors, true
}
