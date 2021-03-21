package user

import (
	"database/sql"
	db "github.com/desfpc/Wishez_DB"
	"github.com/desfpc/Wishez_Type"
	"github.com/mitchellh/mapstructure"
)

var dbres *sql.DB

func initDb(){
	dbres = db.Db()
}

//роутер User
func Route() {

}

//получение записи пользователя по id
func getUser(id int) types.JsonAnswerItem {
	var user types.User
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

	item := make(types.JsonAnswerItem,0)
	mapstructure.Decode(user, &item)

	return item
}