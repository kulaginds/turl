package db

import (
	"database/sql"
	"log"
)

type SimpleDB struct {
	db *sql.DB
}

func NewSimpleDB(dsn string) (d *SimpleDB, ok bool) {
	var err error

	d = &SimpleDB{}
	d.db, err = sql.Open(dbDriver, dsn)

	if nil != err {
		log.Fatalln("Can't connect to DB via dsn:", dsn)
		log.Fatalln(err.Error())
		return
	}

	err = d.db.Ping()

	if nil != err {
		log.Fatalln("Can't ping DB")
		log.Fatalln(err.Error())
		return
	}

	ok = true

	return
}

func (d *SimpleDB) Initialize() (ok bool) {
	table := urlsTable
	return CheckTable(d.db, &table, "0")
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
