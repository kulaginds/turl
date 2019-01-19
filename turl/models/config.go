package models

import (
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	serviceUrl  *url.URL
	servicePort string
	abc         string
	abcIdMinLen int
}

func NewConfig() *Config {
	c := Config{}

	c.Initialize()

	return &c
}

func (c *Config) Initialize() {
	serviceUrl, err := url.ParseRequestURI(os.Getenv("SERVICE_URL"))
	servicePort := parseUint16(os.Getenv("SERVICE_PORT"))
	abc := os.Getenv("ABC")
	abcIdMinLen := parseUint16(os.Getenv("ABC_ID_MIN_LENGTH"))

	if nil != err {
		serviceUrl, _ = url.ParseRequestURI("http://localhost:8080")
	}

	if 0 == servicePort {
		servicePort = parseUint16(serviceUrl.Port())

		if 0 == servicePort {
			servicePort = 8080
		}
	}

	if 0 == len(abc) {
		abc = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	if 0 == abcIdMinLen {
		abcIdMinLen = 4
	}

	c.serviceUrl  = serviceUrl
	c.servicePort = ":" + strconv.FormatUint(uint64(servicePort), 10)
	c.abc         = abc
	c.abcIdMinLen = int(abcIdMinLen)
}

func (c *Config) ServiceUrl() *url.URL {
	return c.serviceUrl
}

func (c *Config) ServicePort() string {
	return c.servicePort
}

func (c *Config) ABC() string {
	return c.abc
}

func (c *Config) ABCIdMinLen() int {
	return c.abcIdMinLen
}

func parseUint16(portStr string) (port uint16) {
	p, err := strconv.ParseUint(portStr, 10, 16)

	if nil != err {
		port = uint16(p)
	}

	return
}
