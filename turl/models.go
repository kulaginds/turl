package turl

import (
	"math"
	"net/url"
	"strings"
)

const (
	noErrors uint = 0
	errNotValidUrl uint = 1
	errEmptyUrl uint = 2
	errUrlTooLong uint = 3
)

const (
	emptyStr = ""
	maxUrlLength = 2000

	minABCIdLength = 4
	ABC = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ABCBase = 62
)

type Url struct {
	Id  int
	Url string
}

func (u *Url) Validate() (errorCode uint, ok bool) {
	var length = uint(len(u.Url))

	if 0 == length {
		errorCode, ok = errEmptyUrl, false
		return
	}

	if length > maxUrlLength {
		errorCode, ok = errUrlTooLong, false
		return
	}

	_, err := url.ParseRequestURI(u.Url)

	if nil != err {
		errorCode, ok = errNotValidUrl, false
	}

	return noErrors, true
}

// Return Id in own ABC
// Read from left to right
func (u *Url) GetABCId() string {
	var quotient, remainder = u.Id, 0
	digitsCount := u.getABCIdDigitsCount()
	chars := make([]string, digitsCount)

	for i := 0; i < digitsCount; i++ {
		remainder = quotient % ABCBase
		quotient = quotient / ABCBase

		chars[i] = string(ABC[remainder])
	}

	return strings.Join(chars, emptyStr)
}

func (u *Url) getABCIdDigitsCount() int {
	// log(Url.Id) with ABC base
	l := math.Log2(float64(u.Id)) / math.Log2(float64(ABCBase))
	digitsCount := int(math.Ceil(l))

	if digitsCount < minABCIdLength {
		digitsCount = minABCIdLength
	}

	return digitsCount
}
