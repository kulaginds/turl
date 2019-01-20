package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"net/url"
	"os"
	"strings"
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
	DbDriver     = "mysql"
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
var db *sql.DB

func Initialize(conf *Config) (ok bool) {
	var err error

	config = conf
	abc = NewABC(conf.ABC(), conf.ABCIdMinLen())

	db, err = sql.Open(DbDriver, config.DSN())

	if nil != err {
		fmt.Fprintln(os.Stderr, "Can't connect to DB via dsn:", config.DSN())
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	err = db.Ping()

	if nil != err {
		fmt.Println("Can't ping DB")
		fmt.Println(err.Error())
		return
	}

	return true
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
		return
	}

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

	var rowUrl string

	ok = false
	id := u.parseUrlId()
	row := db.QueryRow("SELECT url FROM urls WHERE id = ?", id)
	err := row.Scan(&rowUrl)

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

func (u *ShortUrl) parseUrlId() int {
	shortUrl, _ := url.ParseRequestURI(u.Url)

	return abc.Decode(shortUrl.Path[1:])
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

func (u *LongUrl) Short() (url *ShortUrl, errorCode uint, ok bool) {
	var result sql.Result
	var id int64

	if errorCode, ok = u.Validate(); !ok {
		return
	}

	ok = false
	stmt, err := db.Prepare("INSERT INTO urls VALUES(NULL, ?)")

	if nil != err {
		errorCode = ErrDbPrepair
		return
	}

	defer stmt.Close()

	result, err = stmt.Exec(u.Url)

	if nil != err {
		errorCode = ErrDbExec
		return
	}

	id, err = result.LastInsertId()

	if nil != err {
		errorCode = ErrDbLastInsertId
		return
	}

	return &ShortUrl{Url:u.getShortUrlById(int(id))}, NoErrors, true
}

func (u *LongUrl) getShortUrlById(id int) string {
	shortUrl := make([]string, 3)

	shortUrl[0] = config.ServiceUrl().String()
	shortUrl[1] = "/"
	shortUrl[2] = abc.Encode(id)

	return strings.Join(shortUrl, EmptyStr)
}
