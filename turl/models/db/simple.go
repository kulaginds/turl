package db

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"os"
)

type SimpleDB struct {
	db *sql.DB
}

func NewSimpleDB(dsn string) (d *SimpleDB, ok bool) {
	var err error

	d = &SimpleDB{}
	d.db, err = sql.Open(dbDriver, dsn)

	if nil != err {
		fmt.Fprintln(os.Stderr, "Can't connect to DB via dsn:", dsn)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	err = d.db.Ping()

	if nil != err {
		fmt.Fprintln(os.Stderr, "Can't ping DB")
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	ok = true

	return
}

func (d *SimpleDB) Initialize() (ok bool) {
	columns := make([]*TableColumnStruct, 2)

	rows, err := d.db.Query("DESCRIBE " + urlsTable)

	if nil != err {
		me, ok := err.(*mysql.MySQLError)

		if !ok || 1146 != me.Number {
			fmt.Fprintln(os.Stderr, "Can't describe urls table")
			fmt.Fprintln(os.Stderr, err.Error())
			return false
		}

		return d.createUrlsTable()
	}

	defer rows.Close()

	i := 0
	for rows.Next() {
		col := new(TableColumnStruct)
		err = rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)

		if nil != err {
			fmt.Fprintln(os.Stderr, "Can't scan columns from urls table")
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		columns[i] = col
		i++
	}

	if err = rows.Err(); nil != err {
		fmt.Fprintln(os.Stderr, "Some error in rows")
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	return d.checkStruct(columns)
}

func (d *SimpleDB) createUrlsTable() (ok bool) {
	sql := "CREATE TABLE `urls` (\n" +
				"`id` int(10) unsigned NOT NULL AUTO_INCREMENT,\n" +
				"`url` text NOT NULL,\n" +
				"PRIMARY KEY (`id`)\n" +
			") ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;\n"

	fmt.Println("Create table urls")

	_, err := d.db.Exec(sql)

	if nil != err {
		fmt.Fprintln(os.Stderr, "Can't create table urls")
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	return true
}

func (d *SimpleDB) checkStruct(columns []*TableColumnStruct) (ok bool) {
	if ok = d.checkStructId(columns[0]); !ok {
		fmt.Fprintln(os.Stderr, "Wrong structure of id field in table urls")
		return
	}

	if ok = d.checkStructUrl(columns[1]); !ok {
		fmt.Fprintln(os.Stderr, "Wrong structure of url field in table urls")
		return
	}

	return true
}

func (d *SimpleDB) checkStructId(idColumn *TableColumnStruct) (ok bool) {
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

func (d *SimpleDB) checkStructUrl(urlColumn *TableColumnStruct) (ok bool) {
	if "url" != *urlColumn.Field {
		return
	}

	if "text" != *urlColumn.Type {
		return
	}

	return true
}

func (d *SimpleDB) GetUrlById(id int) (rowUrl string, err error) {
	tbl := urlsTable
	row := d.db.QueryRow("SELECT url FROM " + tbl + " WHERE id = ?", id)
	err  = row.Scan(&rowUrl)

	return
}

func (d *SimpleDB) AddUrl(rowUrl string) (lastInsertId int64, err error) {
	var result sql.Result

	tbl       := urlsTable
	stmt, err := d.db.Prepare("INSERT INTO " + tbl + " VALUES(NULL, ?)")

	if nil != err {
		return
	}

	defer stmt.Close()

	result, err = stmt.Exec(rowUrl)

	if nil != err {
		return
	}

	lastInsertId, err = result.LastInsertId()

	return
}
