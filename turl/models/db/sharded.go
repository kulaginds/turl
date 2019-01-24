package db

import (
	"database/sql"
	"fmt"
	"os"
)

type ShardedDB struct {
	db           *sql.DB
	shardCount   int
}

func NewShardedDB(dsn string, shardCount int) (d *ShardedDB, ok bool) {
	var err error

	d = &ShardedDB{shardCount:shardCount}
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

	d.Initialize()

	ok = true

	return
}

func (d *ShardedDB) Initialize() (ok bool) {
	maxShardRows := (0xFFFFFFFF - 1) / d.shardCount
	tables := make([]*string, d.shardCount)
	rows, err := d.db.Query("SHOW TABLES LIKE '" + urlsTable + "_%'")

	if nil != err {
		fmt.Fprintln(os.Stderr, "Can't get urls tables")
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	i := 0
	for rows.Next() {
		rows.Scan(tables[i])
		i++
	}

	if i != d.shardCount {
		if 0 != i {
			fmt.Fprintf(os.Stderr, "Wrong shard count: %d retrieve, expected %d\n", i, d.shardCount)
			return
		}

		for i = 0; i < d.shardCount; i++ {
			table := urlsTable + "_" + string(i)
			CreateUrlsTable(d.db, &table, string(maxShardRows * i))
		}

		return true
	}

	for i = 0; i < d.shardCount; i++ {
		if ok = CheckTable(d.db, tables[i], string(maxShardRows * i)); !ok {
			return
		}
	}

	return true
}

func (d *ShardedDB) GetUrlById(id int) (rowUrl string, err error) {
	tbl := d.getShardById(id)
	row := d.db.QueryRow("SELECT url FROM " + tbl + " WHERE id = ?", id)
	err  = row.Scan(&rowUrl)

	return
}

func (d *ShardedDB) getShardById(id int) (shard string) {
	return urlsTable + "_" + string(id % d.shardCount);
}

func (d *ShardedDB) AddUrl(rowUrl string) (lastInsertId int64, err error) {
	var result sql.Result

	// TODO: запускать в одном потоке
	//  чтобы d.getLastInsertId() вернуло актуальное значение
	//  на момент выполнения INSERT-а
	tbl       := d.getShardById(d.getLastInsertId())
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

func (d *ShardedDB) getLastInsertId() (lastInsertId int) {
	shardsMaxId := make([]int, d.shardCount)

	for i := 0; i < d.shardCount; i++ {
		shardsMaxId[i] = d.getShardMaxId(i)
	}

	maxId := maxInt(shardsMaxId)

	return maxId + 1
}

func (d *ShardedDB) getShardMaxId(index int) (maxId int) {
	return
}

func maxInt(arr []int) (maxInt int) {
	maxInt = arr[0]

	for i := 1; i < len(arr); i++ {
		if maxInt < arr[i] {
			maxInt = arr[i]
		}
	}

	return
}
