package db

import (
	"database/sql"
	"github.com/desfpc/Wishez_Helpers"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// Db соединение с БД
func Db(driverName string, dataSourceName string) *sql.DB  {
	if driverName == "" {
		driverName = "mysql"
	}
	if dataSourceName == "" {
		dataSourceName = "root:root@/wishez"
	}
	//db, err := sql.Open("mysql", "root:025sergLBBK1&*@/wishez")
	db, err := sql.Open(driverName, dataSourceName)
	helpers.CheckErr(err)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

// CheckCount Возвращает кол-во строк запроса
func CheckCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err:= rows.Scan(&count)
		helpers.CheckErr(err)
	}
	return count
}