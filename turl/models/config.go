package models

import (
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	serviceUrl  *url.URL
	servicePort string
	abc         string
	abcIdMinLen int
	abcIdPrefix string
	dsn         string
	shardCount  uint16
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
	abcIdPrefix := os.Getenv("ABC_ID_PREFIX")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	shardCount := parseUint16(os.Getenv("SHARD_COUNT"))

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

	if 0 == len(dbHost) {
		dbHost = "127.0.0.1"
	}

	if 0 == len(dbName) {
		dbName = "test"
	}

	if 0 == len(dbUser) {
		dbUser = "root"
	}

	if 0 == shardCount {
		shardCount = 1
	}

	c.serviceUrl  = serviceUrl
	c.servicePort = ":" + strconv.FormatUint(uint64(servicePort), 10)
	c.abc         = abc
	c.abcIdMinLen = int(abcIdMinLen)
	c.abcIdPrefix = abcIdPrefix
	c.dsn         = prepareDsn(dbHost, dbName, dbUser, dbPassword)
	c.shardCount  = shardCount
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

func (c *Config) ABCIdPrefix() string {
	return c.abcIdPrefix
}

func (c *Config) DSN() string {
	return c.dsn
}

func (c *Config) ShardCount() uint16 {
	return c.shardCount
}

func parseUint16(input string) (output uint16) {
	p, err := strconv.ParseUint(input, 10, 16)

	if nil != err {
		output = uint16(p)
	}

	return
}

func prepareDsn(host, name, user, password string) string {
	dsn := make([]string, 8)

	dsn[0] = user
	dsn[1] = ":"
	dsn[2] = password
	dsn[3] = "@"
	dsn[4] = "tcp("
	dsn[5] = host
	dsn[6] = ")/"
	dsn[7] = name

	return strings.Join(dsn, EmptyStr)
}
