package db

import (
	"database/sql"
	"github.com/desfpc/Wishez_Helpers"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var res *sql.DB

// Close закрытие соединения с БД
func Close() {
	res.Close()
	res = nil
}

// Db соединение с БД
func Db(driverName string, dataSourceName string) *sql.DB  {

	if res != nil {
		return res
	}

	if driverName == "" {
		driverName = "mysql"
	}
	if dataSourceName == "" {
		dataSourceName = "root:root@/wishez"
	}
	//db, err := sql.Open("mysql", "root:025sergLBBK1&*@/wishez")
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		res = nil
	} else {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		res = db
	}
	return res
}

// CheckCount Возвращает кол-во строк запроса
func CheckCount(rows *sql.Rows) (count int) {
	if rows == nil {
		return 0
	}
	for rows.Next() {
		err:= rows.Scan(&count)
		helpers.CheckErr(err)
	}
	return count
}