package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"strings"
	turlDb "turl/turl/models/db"
)

const (
	NoErrors       uint = 0
	ErrNotValidUrl uint = 1
	ErrEmptyUrl    uint = 2
	ErrUrlTooLong  uint = 3
	ErrUrlNotFound uint = 4

	ErrDbPrepair      uint = 10
	ErrDbExec         uint = 11
	ErrDbLastInsertId uint = 12
)

var MSG = map[uint]string {
	0:"Request failed",
}

var ErrMsg = map[uint]string {
	NoErrors:       "No errors",
	ErrNotValidUrl: "Field url not valid",
	ErrEmptyUrl:    "Field url is empty",
	ErrUrlTooLong:  "Field url is too long",
	ErrUrlNotFound: "Requested url not found",

	ErrDbPrepair:      MSG[0],
	ErrDbExec:         MSG[0],
	ErrDbLastInsertId: MSG[0],
}

const (
	EmptyStr     = ""
	MaxUrlLength = 2000
)

type ShortUrl struct {
	Url string `json:"url"`
}

type LongUrl struct {
	Id  int    `json:"-"`
	Url string `json:"url"`
}

var abc *ABC
var config *Config
var db turlDb.DB

func Initialize(conf *Config) (ok bool) {
	config = conf
	abc = NewABC(config.ABC(), config.ABCIdMinLen())
	db, ok = turlDb.NewDB(config.DSN(), config.ShardCount())

	return
}

func (u *ShortUrl) Validate() (errorCode uint, ok bool) {
	errorCode, ok = validateUrl(u.Url)

	if !ok {
		return
	}

	shortUrl, _ := url.ParseRequestURI(u.Url)

	if shortUrl.Host != config.ServiceUrl().Host {
		errorCode, ok = ErrNotValidUrl, false
		return
	}

	errorCode, ok = u.validatePath(shortUrl.Path)

	if !ok {
		return
	}

	return NoErrors, true
}

func (u *ShortUrl) Long() (longUrl *LongUrl, errorCode uint, ok bool) {
	if errorCode, ok = u.Validate(); !ok {
		return
	}

	ok           = false
	rowUrl, err := db.GetUrlById(u.parseUrlId())

	if nil != err {
		switch {
		case err == sql.ErrNoRows:
			errorCode = ErrUrlNotFound
			break
		default:
			errorCode = ErrDbExec
		}
		return
	}

	return &LongUrl{Url:rowUrl}, NoErrors, true
}

func (u *ShortUrl) validatePath(path string) (errorCode uint, ok bool) {
	var index int

	errorCode, ok = NoErrors, true

	if "/" == string(path[0]) {
		path = path[1:]
	}

	prefixLen := len(config.ABCIdPrefix())

	if prefixLen > 0 && config.ABCIdPrefix() != path[:prefixLen] {
		return ErrNotValidUrl, false
	}

	path = path[prefixLen:]
	l := len(path)

	for i := 0; i < l; i++ {
		index = abc.IndexOf(path[i])

		if -1 == index {
			return ErrNotValidUrl, false
		}
	}

	return
}

func (u *ShortUrl) parseUrlId() int {
	shortUrl, _ := url.ParseRequestURI(u.Url)
	path        := shortUrl.Path[1:]
	prefixLen   := len(config.ABCIdPrefix())

	if prefixLen > 0 {
		path = path[prefixLen:]
	}

	return abc.Decode(path)
}

func (u *LongUrl) Validate() (errorCode uint, ok bool) {
	return validateUrl(u.Url)
}

func (u *LongUrl) Short() (url *ShortUrl, errorCode uint, ok bool) {
	if errorCode, ok = u.Validate(); !ok {
		return
	}

	ok       = false
	id, err := db.AddUrl(u.Url)

	if nil != err {
		errorCode = ErrDbExec
		return
	}

	return &ShortUrl{Url:u.getShortUrlById(int(id))}, NoErrors, true
}

func (u *LongUrl) getShortUrlById(id int) string {
	shortUrl := make([]string, 4)

	shortUrl[0] = config.ServiceUrl().String()
	shortUrl[1] = "/"
	shortUrl[2] = config.ABCIdPrefix()
	shortUrl[3] = abc.Encode(id)

	return strings.Join(shortUrl, EmptyStr)
}

func validateUrl(u string) (errorCode uint, ok bool) {
	var length = uint(len(u))

	if 0 == length {
		errorCode, ok = ErrEmptyUrl, false
		return
	}

	if length > MaxUrlLength {
		errorCode, ok = ErrUrlTooLong, false
		return
	}

	_, err := url.ParseRequestURI(u)

	if nil != err {
		errorCode, ok = ErrNotValidUrl, false
	}

	return NoErrors, true
}
