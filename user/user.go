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
	"regexp"
	"strconv"
	"time"
)

var key = []byte("Абдб%дв_3453453ы!всв^амвам_DFGVBdf*vdf43*453")
var dbres *sql.DB

func initDb(){
	dbres = db.Db("", "")
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
	case "get":
		body, err = getUserByID(resp, auser)
	case "refreshToken":
		body, err = doRefreshToken(auser, refreshToken)
	//TODO case "addFriend":
	//	body, err = addFriend(resp)
	//TODO case "deleteFriend":
	//	body, err = deleteFriend(resp)
	//TODO case "confirmFriend":
	//	body, err = confirmFriend(resp)
	case "list":
		body, err = getUserList(resp, auser)
	default:
		err, code = helpers.NoRouteErrorAnswer()
	}
	return body, err, code
}

// getUserList TODO получение списка доступных пользователей (друзей)
//
// предполагаемый json запроса:
// {"Entity":"user","Action":"list","Params":{"type":"all","search":"вася"}}
// Entity string - сущность
// Action string - действие
// Params.type string - тип получаемых пользователей: строка из массива ['all','friend','request']
// Params.search string - строка для поиска пользователя по известным данным (имя, email) (необязательно)
func getUserList(resp types.JsonRequest, auser types.User) (types.JsonAnswerBody, types.Errors) {
	var body types.JsonAnswerBody
	Errors := make(types.Errors,0)



	return body, Errors
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

	//складываем, подписываеми и возвращаем токен
	return makeTokenFromStrings(header, body)
}

// makeTokenFromStrings подписывает токен, складывает и кодирует
func makeTokenFromStrings(header string, body string) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(header+body))
	// signature := string(mac.Sum(nil))
	signature := b64.StdEncoding.EncodeToString(mac.Sum(nil))
	return makeTokenFromStringsVsSignature(header, body, signature)
}

// makeTokenFromStringsVsSignature складывает и кодирует токен из строк заголовка, тела и подписи
func makeTokenFromStringsVsSignature(header string, body string, signature string) string {
	return b64.StdEncoding.EncodeToString([]byte(header+body+signature))
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
	signature, _ := b64.StdEncoding.DecodeString(re.ReplaceAllString(tokenString, "$3"))
	deconcactedToken.Signature = string(signature)

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

// comparePasswords проверка пароля
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

// hashAndSalt генерация хэша пароля
func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, 10)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// GetUserFromBD получение пользователя из БД по его ID
func GetUserFromBD(id string) types.User {
	initDb()
	var user types.User
	id = helpers.Escape(id)
	query := "SELECT * FROM users WHERE id = "+id
	results, err := dbres.Query(query)
	helpers.CheckErr(err)

	//перебираем результаты
	for results.Next() {
		//пробуем все запихнуть в user-а
		err = results.Scan(&user.Id, &user.Email, &user.Pass, &user.Fio, &user.Sex, &user.Telegram, &user.Instagram, &user.Twitter, &user.Facebook,
			&user.Phone, &user.Role, &user.Avatar, &user.Google, &user.DateAdd)
		helpers.CheckErr(err)
	}

	return user
}

// ToPublicJson формирование публичного JsonAnswerItem из User
func ToPublicJson (user types.User) types.JsonAnswerItem {
	item := make(types.JsonAnswerItem)
	item["Id"] = strconv.Itoa(user.Id)

	if item["Id"] != "0" {
		item["Fio"] = user.Fio
		item["Sex"] = user.Sex
		item["Avatar"] = helpers.MakeStringFromIntSQL(user.Avatar)
		item["DateAdd"] = user.DateAdd
	}

	return item
}

// ToJson формирование JsonAnswerItem из User
func ToJson (user types.User) types.JsonAnswerItem {
	item := make(types.JsonAnswerItem)
	item["Id"] = strconv.Itoa(user.Id)

	if item["Id"] != "0" {
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
		item["DateAdd"] = user.DateAdd
	}

	return item
}

// doRefreshToken обновление токенов
//
// предполагаемый json запроса:
// {"entity":"user","action":"refreshToken"}
// entity string - сущность
// action string - действие
func doRefreshToken(auser types.User, refreshToken string) (types.JsonAnswerBody, types.Errors) {
	var tuser types.User
	var body types.JsonAnswerBody
	errors := make(types.Errors,0)
	authorizeError := true
	expireError:= false
	tuser, authorizeError, expireError = GetAuthorization(refreshToken, "refresh")

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

// authorizeUser авторизация пользователя
//
// предполагаемый json запроса:
// {"Entity":"user","Action":"authorize","Params":{"login":"UserLogin","pass":"UserPassword"}}
// Entity string - сущность
// Action string - действие
// Params.login string - логин (email) пользователя
// Params.pass string - пароль пользователя
func authorizeUser(resp types.JsonRequest) (types.JsonAnswerBody, types.Errors) {
	var body types.JsonAnswerBody
	var params = resp.Params
	var exist bool
	Errors := make(types.Errors,0)

	//проверка на наличие логина
	var login string
	login, Errors, exist = helpers.ParamFromJsonRequest(params, "login", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на наличае пароля
	var pass string
	pass, Errors, exist = helpers.ParamFromJsonRequest(params, "pass", Errors)
	if !exist {
		return body, Errors
	}

	//получаем запись в БД по логину
	initDb()
	login = helpers.Escape(login)
	query := "SELECT * FROM users WHERE email = '"+login+"'"
	//log.Printf("query: "+query)
	results, err := dbres.Query(query)
	helpers.CheckErr(err)

	var user types.User

	//перебираем результаты
	for results.Next() {
		//пробуем все запихнуть в user-а
		err = results.Scan(&user.Id, &user.Email, &user.Pass, &user.Fio, &user.Sex, &user.Telegram, &user.Instagram, &user.Twitter, &user.Facebook,
			&user.Phone, &user.Role, &user.Avatar, &user.Google, &user.DateAdd)
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

// registerUser регистрация нового пользователя
//
// предполагаемый json запроса:
// {"Entity":"user","Action":"register","Params":{"login":"UserLogin","pass":"UserPassword"}}
// Entity string - сущность
// Action string - действие
// Params.login string - логин (email) пользователя
// Params.pass string - пароль пользователя
func registerUser(resp types.JsonRequest) (types.JsonAnswerBody, types.Errors) {

	var body types.JsonAnswerBody
	var params = resp.Params
	var exist bool
	Errors := make(types.Errors,0)

	//проверка на наличие логина
	var login string
	login, Errors, exist = helpers.ParamFromJsonRequest(params, "login", Errors)
	if !exist {
		return body, Errors
	}

	if !helpers.IsEmailValid(login) { //валидация login как email
		Errors = append(Errors, "Not valid login email")
		return body, Errors
	}

	//проверка на наличае пароля
	var pass string
	pass, Errors, exist = helpers.ParamFromJsonRequest(params, "pass", Errors)
	if !exist {
		return body, Errors
	}

	//проверка пользователя в базе
	initDb()
	login = helpers.Escape(login)
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

// getUserByID получение записи пользователя по id
//
// предполагаемый json запроса:
// {"Entity":"user","Action":"get","id":"1"}
// Entity string - сущность
// Action string - действие
// Id string - ID пользователя (число в виде строки)
func getUserByID(resp types.JsonRequest, auser types.User) (types.JsonAnswerBody, types.Errors) {
	var body types.JsonAnswerBody
	Errors := make(types.Errors,0)

	//проверка на наличие id
	if resp.Id == "" {
		Errors = append(Errors, "No Id in Request")
		return body, Errors
	}

	id := resp.Id
	user := GetUserFromBD(id)

	//проверка прав на просмотр пользователя для формирования коллекции (выводить все или короткие сведения)
	friendRights := true

	if user.Id != auser.Id {
		initDb()
		var counter int
		query := "SELECT count(*) FROM `users_friends` WHERE ((`user_id` = ? AND `friend_id` = ?) OR (`user_id` = ? AND `friend_id` = ?)) AND approved = 'Y'"
		dbres.QueryRow(query, user.Id, auser.Id, auser.Id, user.Id).Scan(&counter)
		if counter == 0 {
			friendRights = false
		}
	}
	var item types.JsonAnswerItem
	if friendRights {
		item = ToJson(user)
	} else {
		item = ToPublicJson(user)
	}

	if item["Id"] == "0" {
		Errors = append(Errors, "No user with Id: "+id)
		return body, Errors
	}

	body.Items = make([]types.JsonAnswerItem,0)
	body.Items = append(body.Items, item)

	return body, Errors
}