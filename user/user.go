package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/desfpc/Wishez_DB"
	"github.com/desfpc/Wishez_Helpers"
	"github.com/desfpc/Wishez_Type"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var key = []byte("Абдб%дв_3453453ы!всв^амвам_DFGVBdf*vdf43*453")
var dbres *sql.DB
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func initDb(){
	dbres = db.Db("", "")
}

// MakeToken функция генерирует токен
func MakeToken(kind string, user types.User) string {

	//заголовок токена
	header := "{\"alg\":\"HS256\",\"typ\":\"JWT\"}"

	//ID пользователя
	id := strconv.Itoa(user.Id)

	//время жизни
	var lifetime string
	switch kind {
	case "access":
		lifetime = strconv.FormatInt(time.Now().Unix() + 1800, 10)
	case "refresh":
		lifetime = strconv.FormatInt(time.Now().Unix() + 5184000, 10)
	}

	//тело токена
	body := "{\"user_id\":\""+id+"\",\"exp\":\""+lifetime+"\",\"kind\":\""+kind+"\"}"

	//подпись
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(header+body))
	signature := string(mac.Sum(nil))

	var token = b64.StdEncoding.EncodeToString([]byte(header+body+signature))

	return token
}

// deconcatToken преобразование токена в читабельный вид
func deconcatToken(token string) types.Token {

	//декодируем base64 токен
	normalToken, _ := b64.StdEncoding.DecodeString(token)
	tokenString := string(normalToken)

	//паттерн для токена
	re := regexp.MustCompile("(\\{.+?\\})(\\{.+?\\})(.*)")

	//заполняем токен
	var deconcactedToken types.Token

	deconcactedToken.Head = re.ReplaceAllString(tokenString, "$1")
	deconcactedToken.Body = re.ReplaceAllString(tokenString, "$2")
	deconcactedToken.Signature = re.ReplaceAllString(tokenString, "$3")

	return deconcactedToken
}

// checkToken проверка токена на валидность
func checkToken(token types.Token) bool {

	//подпись
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(token.Head + token.Body))
	signature := mac.Sum(nil)

	return hmac.Equal(signature, []byte(token.Signature))
}

// CheckUserToken проверка токена-строки на валидность
func CheckUserToken(token string) bool {
	var dToken = deconcatToken(token)
	return checkToken(dToken)
}

// GetAuthorization проверка авторизации, получение активного пользователя
func GetAuthorization(token string, kind string) (types.User, bool, bool) { //user, authorizeError, expireError
	dToken := deconcatToken(token)
	bodyString := dToken.Body
	body := make(types.TokenBody)
	var auser types.User

	if bodyString == "" {
		return auser, true, true
	}

	err := json.Unmarshal([]byte(bodyString), &body)
	if err != nil {
		log.Printf("Error reading JSON from token body: %v", err)
		log.Printf(bodyString)
		return auser, true, false
	}

	if !CheckUserToken(token) {
		return auser, true, false
	}

	if body["kind"] != kind {
		return auser, true, false
	}

	if body["exp"] == "" {
		return auser, false, true
	}

	//var exp int64
	exp, err := strconv.ParseInt(body["exp"], 10, 64)
	helpers.CheckErr(err)

	//получение данных пользователя
	auser = GetUserFromBD(body["user_id"])

	//если токен протух
	if exp < time.Now().Unix() {
		return auser, false, true
	}

	return auser, false, false
}

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

// Route роутер User
func Route(resp types.JsonRequest, auser types.User, refreshToken string) (types.JsonAnswerBody, types.Errors, int) {
	var body types.JsonAnswerBody
	var err types.Errors
	var code = 200

	//проверяем метод
	switch resp.Action {
	case "register":
		body, err = registerUser(resp)
	case "authorize":
		body, err = authorizeUser(resp)
	case "getById":
		body, err = getUserByID(resp)
	case "refreshToken":
		body, err = doRefreshToken(auser, refreshToken)
	default:
		err, code = helpers.NoRouteErrorAnswer()
	}
	return body, err, code
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

//обновление токенов
func doRefreshToken(auser types.User, refreshToken string) (types.JsonAnswerBody, types.Errors) {
	var tuser types.User
	var body types.JsonAnswerBody
	errors := make(types.Errors,0)
	authorizeError := true
	expireError:= false
	tuser, authorizeError, expireError = GetAuthorization(refreshToken, "refresh")
	//tuser, _, _ = GetAuthorization(refreshToken)
	if !authorizeError && !expireError && tuser.Id == auser.Id {
		item := make(types.JsonAnswerItem)
		item["accessToken"] = MakeToken("access", auser)
		item["refreshToken"] = MakeToken("refresh", auser)
		body.Items = make([]types.JsonAnswerItem,0)
		body.Items = append(body.Items, item)
	} else {
		errors = append(errors, "Wrong refreshToken")
	}
	return body, errors
}

//авторизация пользователя
func authorizeUser(resp types.JsonRequest) (types.JsonAnswerBody, types.Errors) {
	var body types.JsonAnswerBody
	var params = resp.Params
	Errors := make(types.Errors,0)

	//проверка на наличае логина
	var login, existsLogin = params["login"]
	if !existsLogin {
		Errors = append(Errors, "No login")
		return body, Errors
	}
	login = db.Escape(login)  //для запроса в БД

	//проверка на наличае пароля
	var pass, existsPass = params["pass"]
	if !existsPass {
		Errors = append(Errors, "No pass")
		return body, Errors
	}

	//получаем запись в БД по логину
	initDb()
	query := "SELECT * FROM users WHERE email = '"+login+"'"
	//log.Printf("query: "+query)
	results, err := dbres.Query(query)
	helpers.CheckErr(err)

	var user types.User

	//перебираем результаты
	for results.Next() {
		//пробуем все запихнуть в user-а
		err = results.Scan(&user.Id, &user.Email, &user.Pass, &user.Fio, &user.Sex, &user.Telegram, &user.Instagram, &user.Twitter, &user.Facebook,
			&user.Phone, &user.Role, &user.Avatar, &user.Google, &user.CreatedAt)
		helpers.CheckErr(err)
	}

	if comparePasswords(user.Pass, pass) {
		item := make(types.JsonAnswerItem)
		item["accessToken"] = MakeToken("access", user)
		item["refreshToken"] = MakeToken("refresh", user)
		body.Items = make([]types.JsonAnswerItem,0)
		body.Items = append(body.Items, item)
	} else {
		Errors = append(Errors, "Wrong user login or password")
		return body, Errors
	}

	return body, Errors
}

//регистрация нового пользователя
func registerUser(resp types.JsonRequest) (types.JsonAnswerBody, types.Errors) {

	var body types.JsonAnswerBody
	var params = resp.Params
	Errors := make(types.Errors,0)

	//проверка на наличае логина
	var login, existsLogin = params["login"]
	if !existsLogin {
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
		Errors = append(Errors, "No pass")
		return body, Errors
	}

	//проверка пользователя в базе
	initDb()
	query := "SELECT count(id) count FROM users WHERE email = '"+login+"'"
	//log.Printf("query: "+query)
	results, err := dbres.Query(query)
	helpers.CheckErr(err)

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
	helpers.CheckErr(err)

	lastId, err := res.LastInsertId()
	helpers.CheckErr(err)

	item := make(types.JsonAnswerItem)
	item["Login"] = login
	item["Id"] = strconv.FormatInt(lastId, 10)

	body.Items = make([]types.JsonAnswerItem,0)
	body.Items = append(body.Items, item)

	return body, Errors
}

// GetUserFromBD получение пользователя из БД по его ID
func GetUserFromBD(id string) types.User {
	initDb()
	var user types.User

	query := "SELECT * FROM users WHERE id = "+id
	results, err := dbres.Query(query)
	helpers.CheckErr(err)

	//перебираем результаты
	for results.Next() {
		//пробуем все запихнуть в user-а
		err = results.Scan(&user.Id, &user.Email, &user.Pass, &user.Fio, &user.Sex, &user.Telegram, &user.Instagram, &user.Twitter, &user.Facebook,
			&user.Phone, &user.Role, &user.Avatar, &user.Google, &user.CreatedAt)
		helpers.CheckErr(err)
	}

	return user
}

// ToJson формирование JsonAnswerItem из User
func ToJson (user types.User) types.JsonAnswerItem {

	item := make(types.JsonAnswerItem)
	item["id"] = strconv.Itoa(user.Id)

	if item["id"] != "0" {
		item["Email"] = user.Email
		item["Fio"] = user.Fio
		item["Sex"] = user.Sex
		item["Telegram"] = helpers.MakeStringFromSQL(user.Telegram)
		item["Instagram"] = helpers.MakeStringFromSQL(user.Instagram)
		item["Twitter"] = helpers.MakeStringFromSQL(user.Twitter)
		item["Facebook"] = helpers.MakeStringFromSQL(user.Facebook)
		item["Phone"] = helpers.MakeStringFromSQL(user.Phone)
		item["Role"] = user.Role
		item["Avatar"] = helpers.MakeStringFromIntSQL(user.Avatar)
		item["Google"] = helpers.MakeStringFromSQL(user.Google)
	}

	return item
}

//получение записи пользователя по id
func getUserByID(resp types.JsonRequest) (types.JsonAnswerBody, types.Errors) {

	var body types.JsonAnswerBody
	var params = resp.Params
	Errors := make(types.Errors,0)

	//проверка на наличае id
	var id, existsId = params["id"]
	if !existsId {
		Errors = append(Errors, "No user Id")

		return body, Errors
	}

	user := GetUserFromBD(id)
	item := ToJson(user)

	if item["id"] == "0" {
		Errors = append(Errors, "No user with Id: "+id)
		return body, Errors
	}

	body.Items = make([]types.JsonAnswerItem,0)
	body.Items = append(body.Items, item)

	return body, Errors
}