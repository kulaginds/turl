package models

import (
	"database/sql"
	"fmt"
	"os"
)

const DbDriver = "mysql"

type DB struct {
	db *sql.DB
}

func NewDB(dsn string) (d *DB, ok bool) {
	var err error

	d         = &DB{}
	d.db, err = sql.Open(DbDriver, dsn)

	if nil != err {
		fmt.Fprintln(os.Stderr, "Can't connect to DB via dsn:", config.DSN())
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	err = d.db.Ping()

	if nil != err {
		fmt.Println("Can't ping DB")
		fmt.Println(err.Error())
		return
	}

	ok = true

	return
}

func (d *DB) GetUrlById(id int) (rowUrl string, err error) {
	row := d.db.QueryRow("SELECT url FROM urls WHERE id = ?", id)
	err  = row.Scan(&rowUrl)

	return
}

func (d *DB) AddUrl(rowUrl string) (lastInsertId int64, err error) {
	var result sql.Result

	stmt, err := d.db.Prepare("INSERT INTO urls VALUES(NULL, ?)")

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
