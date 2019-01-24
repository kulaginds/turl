package db

import (
	"database/sql"
	"fmt"
	"os"
)

type ShardedDB struct {
	db *sql.DB
}

func NewShardedDB(dsn string, shardCount uint16) (d *ShardedDB, ok bool) {
	var err error

	d = &ShardedDB{}
	d.db, err = sql.Open(dbDriver, dsn)

	if nil != err {
		fmt.Fprintln(os.Stderr, "Can't connect to DB via dsn:", dsn)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	err = d.db.Ping()

	if nil != err {
		fmt.Println("Can't ping DB")
		fmt.Println(err.Error())
		return
	}

	d.Initialize()

	ok = true

	return
}

func (d *ShardedDB) Initialize() (ok bool) {
	// TODO: проверить структуру базы
	// если нет таблицы urlsTable, то создать
	// если количество шардов больше 1, найти lastInsertId,
	return true
}

func (d *ShardedDB) GetUrlById(id int) (rowUrl string, err error) {
	// TODO: если количество шардов больше 1, найти таблицу текущего id
	tbl := urlsTable
	row := d.db.QueryRow("SELECT url FROM " + tbl + " WHERE id = ?", id)
	err  = row.Scan(&rowUrl)

	return
}

func (d *ShardedDB) AddUrl(rowUrl string) (lastInsertId int64, err error) {
	var result sql.Result

	// TODO: если префикс не пустой и запись с lastInsertId
	//  существует, найти новый lastInsertId
	//  чтобы можно было создавать записи с ID=lastInsertId+1

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
