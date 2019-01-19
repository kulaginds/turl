package models

import (
	"math"
	"strings"
)

// Tool for encoding and decoding from decimal to own ABC
type ABC struct {
	// min ABC id length to return in Encode
	abcIdMinLen int
	// ABC characters
	chars       string
	// base of ABC
	base        int
}

// Return new ABC with choosen chars and ID min length
func NewABC(chars string, abcIdMinLength int) *ABC {
	abc := ABC{
		abcIdMinLen: abcIdMinLength,
		chars:       chars,
	}
	abc.base = len(abc.chars)

	return &abc
}

// Encode decimal Id to ABC Id
// Read from right to left
func (a *ABC) Encode(abcId int) string {
	var quotient, remainder = abcId, 0

	digitsCount := a.getIdDigitsCount(abcId)
	chars       := make([]string, digitsCount)

	for i := 0; i < digitsCount; i++ {
		remainder = quotient % a.base
		quotient  = quotient / a.base

		chars[i] = string(a.chars[remainder])
	}

	return strings.Join(chars, EmptyStr)
}

// Encode ABC Id to decimal Id
func (a *ABC) Decode(abcId string) (decId int) {
	var abcDigit int

	for i := 0; i < len(abcId); i++ {
		abcDigit = a.IndexOf(abcId[i])

		if -1 == abcDigit {
			return
		}

		decId += abcDigit * int(math.Pow(float64(a.base), float64(i)))
	}

	return
}

// Return count of ABC digits for encode decimal id
func (a *ABC) getIdDigitsCount(decId int) int {
	// example: log62(Url.Id)
	l := math.Log2(float64(decId)) / math.Log2(float64(a.base))
	digitsCount := int(math.Ceil(l))

	if digitsCount < a.abcIdMinLen {
		digitsCount = a.abcIdMinLen
	}

	return digitsCount
}

// Return index of char in ABC
func (a *ABC) IndexOf(char byte) int {
	var index = -1

	for i := 0; i < a.base; i++ {
		if a.chars[i] == char {
			index = i
			break
		}
	}

	return index
}
