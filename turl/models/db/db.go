package db

import (
	"fmt"
	"os"
)

const (
	dbDriver = "mysql"
	urlsTable = "urls"
)

type TableColumnStruct struct {
	Field   *string
	Type    *string
	Null    *string
	Key     *string
	Default *string
	Extra   *string
}

type DB interface {
	Initialize() (bool)
	GetUrlById(int) (string, error)
	AddUrl(string) (int64, error)
}

func NewDB(dsn string, shardCount uint16) (d DB, ok bool) {
	if shardCount > 1 {
		d, ok = NewShardedDB(dsn, shardCount)
	} else {
		d, ok = NewSimpleDB(dsn)
	}

	if ok = d.Initialize(); !ok {
		return
	}

	return
}

func CheckStruct(columns []*TableColumnStruct) (ok bool) {
	if ok = checkStructId(columns[0]); !ok {
		fmt.Fprintln(os.Stderr, "Wrong structure of id field in table urls")
		return
	}

	if ok = checkStructUrl(columns[1]); !ok {
		fmt.Fprintln(os.Stderr, "Wrong structure of url field in table urls")
		return
	}

	return true
}

func checkStructId(idColumn *TableColumnStruct) (ok bool) {
	if "id" != *idColumn.Field {
		return
	}

	if "int(10) unsigned" != *idColumn.Type {
		return
	}

	if "PRI" != *idColumn.Key {
		return
	}

	if "auto_increment" != *idColumn.Extra {
		return
	}

	return true
}

func checkStructUrl(urlColumn *TableColumnStruct) (ok bool) {
	if "url" != *urlColumn.Field {
		return
	}

	if "text" != *urlColumn.Type {
		return
	}

	return true
}
