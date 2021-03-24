package user

import (
	"database/sql"
	db "github.com/desfpc/Wishez_DB"
	"github.com/desfpc/Wishez_Type"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
)

/*item := make(types.JsonAnswerItem)
item["Login"] = "test"
body.Items = make([]types.JsonAnswerItem,0)
//body.Items[]=item
body.Items = append(body.Items, item)
var jsonItem, _ = json.Marshal(item)
log.Printf(string(jsonItem))*/

var dbres *sql.DB
var Errors types.Errors

func initDb(){
	dbres = db.Db()
}

//проверка пароля
func comparePasswords(hashedPwd string, plainPwd []byte) bool {

	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

//генерация хэша пароля
func hashAndSalt(pwd []byte) string {

	hash, err := bcrypt.GenerateFromPassword(pwd, 10)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

//роутер User
func Route(resp types.JsonRequest) types.JsonAnswerBody {

	var body types.JsonAnswerBody

	//проверяем метод
	switch resp.Action {
		//регистрация пользователя
		case "register":
			body = registerUser(resp)
	}

	return body
}

//TODO регистрация нового пользователя
func registerUser(resp types.JsonRequest) types.JsonAnswerBody {

	var body types.JsonAnswerBody
	var params = resp.Params

	//проверка на наличае логина
	var login, exists = params["login"]
	if(!exists){
		//body = make(types.JsonAnswerBody,0)
		Errors = make(types.Errors,0)
		Errors = append(Errors, "No login")

		return body

	} else {

		//TODO валидация login как email

		login = db.Escape(login) //для запроса в БД

		//проверка на наличае пароля
		var pass, exists = params["pass"]
		if(!exists){
			Errors = make(types.Errors,0)
			Errors = append(Errors, "No login")
			return body
		}

		//проверка пользователя в базе
		initDb()
		query := "SELECT count(id) count FROM users WHERE email = '"+login+"'"
		log.Printf("query: "+query)
		results, err := dbres.Query(query)
		db.CheckErr(err)

		count := db.CheckCount(results)

		if(count > 0){
			Errors = make(types.Errors,0)
			Errors = append(Errors, "Login email "+login+" is already in use")
			return body
		}

		//регистрация пользователя
		passHash := hashAndSalt([]byte(pass)) //хэш пароля

		res, err := dbres.Exec("INSERT INTO users (id, email, pass, fio, role, date_add) " +
			"VALUES (null, ?, ?, ?, ?, NOW())",
			login, passHash, "Unknown", "user")
		db.CheckErr(err)

		lastId, err := res.LastInsertId()
		db.CheckErr(err)

		item := make(types.JsonAnswerItem)
		item["Login"] = login
		item["Id"] = strconv.FormatInt(lastId, 10)

		body.Items = make([]types.JsonAnswerItem,0)
		body.Items = append(body.Items, item)

	}

	return body
}

//получение записи пользователя по id
func getUserByID(id int) types.JsonAnswerItem {
	initDb()
	var user types.User

	results, err := dbres.Query("SELECT * FROM users WHERE id = "+string(id))

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	//перебираем результаты
	for results.Next() {
		//пробуем все запихнуть в user-а
		err = results.Scan(&user.Id, &user.Email, &user.Pass, &user.Fio, &user.Sex, &user.Telegram, &user.Instagram, &user.Twitter, &user.Facebook,
			&user.Phone, &user.Role, &user.Avatar, &user.Google)

		db.CheckErr(err)
	}

	item := make(types.JsonAnswerItem,0)
	mapstructure.Decode(user, &item)

	return item
}