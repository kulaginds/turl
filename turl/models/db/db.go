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

func NewDB(dsn string, shardCount int) (d DB, ok bool) {
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

func CheckTable(db *sql.DB, table *string, autoIncrement string) (ok bool) {
	columns := make([]*TableColumnStruct, 2)

	rows, err := db.Query("DESCRIBE " + *table)

	if nil != err {
		me, ok := err.(*mysql.MySQLError)

		if !ok || 1146 != me.Number {
			fmt.Fprintln(os.Stderr, "Can't describe table " + *table)
			fmt.Fprintln(os.Stderr, err.Error())
			return false
		}

		return CreateUrlsTable(db, table, autoIncrement)
	}

	defer rows.Close()

	i := 0
	for rows.Next() {
		col := new(TableColumnStruct)
		err = rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)

		if nil != err {
			fmt.Fprintln(os.Stderr, "Can't scan columns from table " + *table)
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		columns[i] = col
		i++
	}

	if err = rows.Err(); nil != err {
		fmt.Fprintln(os.Stderr, "Some error in rows in table " + *table)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	return checkStruct(columns, table)
}

func CreateUrlsTable(db *sql.DB, table *string, autoIncrement string) (ok bool) {
	sql := "CREATE TABLE `" + *table + "` (\n" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT,\n" +
		"`url` text NOT NULL,\n" +
		"PRIMARY KEY (`id`)\n" +
		") ENGINE=InnoDB AUTO_INCREMENT=" + autoIncrement + " DEFAULT CHARSET=utf8;\n"

	fmt.Println("Create table urls")

	_, err := db.Exec(sql)

	if nil != err {
		fmt.Fprintln(os.Stderr, "Can't create table urls")
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	return true
}

func checkStruct(columns []*TableColumnStruct, table *string) (ok bool) {
	var field string

	if field, ok = checkStructId(columns[0]); !ok {
		fmt.Fprintln(os.Stderr, "Wrong structure of id field, column " + field + " in table " + *table)
		return
	}

	if field, ok = checkStructUrl(columns[1]); !ok {
		fmt.Fprintln(os.Stderr, "Wrong structure of url field, column " + field + " in table " + *table)
		return
	}

	return true
}

func checkStructId(idColumn *TableColumnStruct) (field string, ok bool) {
	if "id" != *idColumn.Field {
		field = "Field"
		return
	}

	if "int(10) unsigned" != *idColumn.Type {
		field = "Type"
		return
	}

	if "PRI" != *idColumn.Key {
		field = "Key"
		return
	}

	if "auto_increment" != *idColumn.Extra {
		field = "Extra"
		return
	}

	return "", true
}

func checkStructUrl(urlColumn *TableColumnStruct) (field string, ok bool) {
	if "url" != *urlColumn.Field {
		field = "Field"
		return
	}

	if "text" != *urlColumn.Type {
		field = "Type"
		return
	}

	return "", true
}
