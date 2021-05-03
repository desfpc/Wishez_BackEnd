package authorize

import (
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/desfpc/Wishez_Helpers"
	"github.com/desfpc/Wishez_Type"
	"github.com/desfpc/Wishez_User"
	"log"
	"regexp"
	"strconv"
	"time"
)

var key = []byte("Абдб%дв_3453453ы!всв^амвам_DFGVBdf*vdf43*453")

// MakeToken функция генерирует токен
func MakeToken(iss string, kind string, user types.User) string {

	//заголовок токена
	header := "{\"alg\":\"HS256\",\"typ\":\"JWT\"}"

	//ID пользователя
	id := strconv.Itoa(user.Id)

	//время жизни
	var lifetime string
	switch kind {
		case "access":
			lifetime = string(time.Now().Unix() + 1800)
		case "refresh":
			lifetime = string(time.Now().Unix() + 5184000)
	}

	//тело токена
	body := "{\"user_id\":"+id+"\"exp\":"+lifetime+"}"

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

	//паттерн для токена
	re := regexp.MustCompile("^{(.*)}{(.*)}(.*)$")

	//заполняем токен
	var deconcactedToken types.Token

	deconcactedToken.Head = "{" + re.ReplaceAllString(string(normalToken), "$1") + "}"
	deconcactedToken.Body = "{" + re.ReplaceAllString(string(normalToken), "$2") + "}"
	deconcactedToken.Signature = re.ReplaceAllString(string(normalToken), "$3")

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
func GetAuthorization(token string) (types.User, bool, bool) { //user, authorizeError, expireError
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
		return auser, true, false
	}

	if !CheckUserToken(token) {
		return auser, true, false
	}

	if body["exp"] == "" {
		return auser, false, true
	}

	//var exp int64
	exp, err := strconv.ParseInt(body["exp"], 10, 64)
	helpers.CheckErr(err)

	//если токен протух
	if exp < time.Now().Unix() {
		return auser, false, true
	}

	//получение данных пользователя
	auser = user.GetUserFromBD(body["user_id"])

	return auser, false, false
}