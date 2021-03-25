package user

import (
	"database/sql"
	"github.com/desfpc/Wishez_DB"
	"github.com/desfpc/Wishez_Type"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

/*item := make(types.JsonAnswerItem)
item["Login"] = "test"
body.Items = make([]types.JsonAnswerItem,0)
//body.Items[]=item
body.Items = append(body.Items, item)
var jsonItem, _ = json.Marshal(item)
log.Printf(string(jsonItem))*/

var dbres *sql.DB

func initDb(){
	dbres = db.Db()
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}

//роутер User
func Route(resp types.JsonRequest) (types.JsonAnswerBody, types.Errors) {

	var body types.JsonAnswerBody
	var err types.Errors

	//проверяем метод
	switch resp.Action {
	//регистрация пользователя
	case "register":
		body, err = registerUser(resp)
	case "getById":
		body, err = getUserByID(resp)
	case "authorize":
		body, err = authorize(resp)

	}

	return body, err
}

//проверка пароля
func comparePasswords(hashedPwd string, plainPwd string) bool {

	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
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

//авторизация пользователя
func authorize(resp types.JsonRequest) (types.JsonAnswerBody, types.Errors) {
	var body types.JsonAnswerBody
	//var params = resp.Params
	Errors := make(types.Errors,0)


	return body, Errors
}

//регистрация нового пользователя
func registerUser(resp types.JsonRequest) (types.JsonAnswerBody, types.Errors) {

	var body types.JsonAnswerBody
	var params = resp.Params
	Errors := make(types.Errors,0)

	//проверка на наличае логина
	var login, existsLogin = params["login"]
	if(!existsLogin){
		//body = make(types.JsonAnswerBody,0)
		Errors := make(types.Errors,0)
		Errors = append(Errors, "No login")

		return body, Errors
	}

	login = db.Escape(login)  //для запроса в БД
	if !IsEmailValid(login) { //валидация login как email
		Errors = append(Errors, "Not valid login email")
		return body, Errors
	}


	//проверка на наличае пароля
	var pass, existsPass = params["pass"]
	if !existsPass {
		Errors = append(Errors, "No login")
		return body, Errors
	}

	//проверка пользователя в базе
	initDb()
	query := "SELECT count(id) count FROM users WHERE email = '"+login+"'"
	//log.Printf("query: "+query)
	results, err := dbres.Query(query)
	db.CheckErr(err)

	count := db.CheckCount(results)

	if count > 0 {
		Errors = append(Errors, "Login email "+login+" is already in use")
		return body, Errors
	}

	//регистрация пользователя
	passHash := hashAndSalt([]byte(pass)) //хэш пароля

	res, err := dbres.Exec("INSERT INTO users (id, email, pass, fio, role, date_add, sex) " +
		"VALUES (null, ?, ?, ?, ?, NOW(), ?)",
		login, passHash, "Unknown", "user", "other")
	db.CheckErr(err)

	lastId, err := res.LastInsertId()
	db.CheckErr(err)

	item := make(types.JsonAnswerItem)
	item["Login"] = login
	item["Id"] = strconv.FormatInt(lastId, 10)

	body.Items = make([]types.JsonAnswerItem,0)
	body.Items = append(body.Items, item)



	return body, Errors

}

//получение записи пользователя по id
func getUserByID(resp types.JsonRequest) (types.JsonAnswerBody, types.Errors) {

	//TODO сделать проверку на права пользователя (собственный id, id в списке друзей, права админа или модератора)

	var body types.JsonAnswerBody
	var params = resp.Params
	Errors := make(types.Errors,0)

	//проверка на наличае id
	var id, existsId = params["id"]
	if !existsId {
		Errors = append(Errors, "No user Id")

		return body, Errors
	}

	initDb()
	var user types.User

	query := "SELECT * FROM users WHERE id = "+id
	results, err := dbres.Query(query)
	//log.Printf("query: "+query)
	db.CheckErr(err)

	//перебираем результаты
	for results.Next() {
		//пробуем все запихнуть в user-а
		err = results.Scan(&user.Id, &user.Email, &user.Pass, &user.Fio, &user.Sex, &user.Telegram, &user.Instagram, &user.Twitter, &user.Facebook,
			&user.Phone, &user.Role, &user.Avatar, &user.Google, &user.CreatedAt)

		db.CheckErr(err)
	}

	item := make(types.JsonAnswerItem)
	item["id"] = strconv.Itoa(user.Id)

	if item["id"] == "0" {
		Errors = append(Errors, "No user with Id: "+id)

		return body, Errors
	}

	item["Email"] = user.Email
	item["Fio"] = user.Fio
	item["Sex"] = user.Sex

	if !user.Telegram.Valid {
		item["Telegram"] = ""
	} else {
		item["Telegram"] = user.Telegram.String
	}

	if !user.Instagram.Valid {
		item["Instagram"] = ""
	} else {
		item["Instagram"] = user.Instagram.String
	}

	if !user.Twitter.Valid {
		item["Twitter"] = ""
	} else {
		item["Twitter"] = user.Twitter.String
	}

	if !user.Facebook.Valid {
		item["Facebook"] = ""
	} else {
		item["Facebook"] = user.Facebook.String
	}

	if !user.Phone.Valid {
		item["Phone"] = ""
	} else {
		item["Phone"] = user.Phone.String
	}

	item["Role"] = user.Role

	if !user.Avatar.Valid {
		item["Avatar"] = ""
	} else {
		item["Avatar"] = strconv.FormatInt(user.Avatar.Int64, 10)
	}

	if !user.Google.Valid {
		item["Google"] = ""
	} else {
		item["Google"] = user.Phone.String
	}

	body.Items = make([]types.JsonAnswerItem,0)
	body.Items = append(body.Items, item)

	return body, Errors
}