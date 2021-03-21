package user

import (
	"database/sql"
	db "github.com/desfpc/Wishez_DB"
)

var dbres *sql.DB

func initDb(){
	dbres = db.Db()
}

//тип для пользователя
type User struct {

}

//получение записи пользователя по id
func getUser(id int) User {
	var user User
	initDb()

	return user
}