package db

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
