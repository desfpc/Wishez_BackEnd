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
	Id int
	Email string
	Pass string
	Fio string
	Sex string
	Telegram string
	Instagram string
	Twitter string
	Facebook string
	Phone string
	Role string
	Avatar int
	Google string
}

//получение записи пользователя по id
func getUser(id int) User {
	var user User
	initDb()

	results, err := dbres.Query("SELECT * FROM users WHERE id = "+string(id))

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	//перебираем результаты
	for results.Next() {
		//пробуем все запихнуть в user-а
		err = results.Scan(&user.Id, &user.Email, &user.Pass, &user.Fio, &user.Sex, &user.Telegram, &user.Instagram, &user.Twitter, &user.Facebook,
			&user.Phone, &user.Role, &user.Avatar, &user.Google)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	return user
}