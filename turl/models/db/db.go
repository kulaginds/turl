package db

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
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

func CheckTable(db *sql.DB, table string) (ok bool) {
	columns := make([]*TableColumnStruct, 2)

	rows, err := db.Query("DESCRIBE " + table)

	if nil != err {
		me, ok := err.(*mysql.MySQLError)

		if !ok || 1146 != me.Number {
			fmt.Fprintln(os.Stderr, "Can't describe table " + table)
			fmt.Fprintln(os.Stderr, err.Error())
			return false
		}

		return createUrlsTable(db, table)
	}

	defer rows.Close()

	i := 0
	for rows.Next() {
		col := new(TableColumnStruct)
		err = rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)

		if nil != err {
			fmt.Fprintln(os.Stderr, "Can't scan columns from table " + table)
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		columns[i] = col
		i++
	}

	if err = rows.Err(); nil != err {
		fmt.Fprintln(os.Stderr, "Some error in rows in table " + table)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	return checkStruct(columns, table)
}

func createUrlsTable(db *sql.DB, table string) (ok bool) {
	sql := "CREATE TABLE `" + table + "` (\n" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT,\n" +
		"`url` text NOT NULL,\n" +
		"PRIMARY KEY (`id`)\n" +
		") ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;\n"

	fmt.Println("Create table urls")

	_, err := db.Exec(sql)

	if nil != err {
		fmt.Fprintln(os.Stderr, "Can't create table urls")
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	return true
}

func checkStruct(columns []*TableColumnStruct, table string) (ok bool) {
	if ok = checkStructId(columns[0]); !ok {
		fmt.Fprintln(os.Stderr, "Wrong structure of id field in table " + table)
		return
	}

	if ok = checkStructUrl(columns[1]); !ok {
		fmt.Fprintln(os.Stderr, "Wrong structure of url field in table " + table)
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
