package db

import (
	"database/sql"
	"fmt"
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
	return CheckTable(d.db, urlsTable)
}

func (d *SimpleDB) GetUrlById(id int) (rowUrl string, err error) {
	row := d.db.QueryRow("SELECT url FROM " + urlsTable + " WHERE id = ?", id)
	err  = row.Scan(&rowUrl)

	return
}

func (d *SimpleDB) AddUrl(rowUrl string) (lastInsertId int64, err error) {
	var result sql.Result

	stmt, err := d.db.Prepare("INSERT INTO " + urlsTable + " VALUES(NULL, ?)")

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
