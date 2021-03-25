package authorize

import (
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"github.com/desfpc/Wishez_Type"
	"regexp"
	"strconv"
	"time"
)

var Key = []byte("Абдб%дв_3453453ы!всв^амвам_DFGVBdf*vdf43*453")

//функция генерирует токен
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
	mac := hmac.New(sha256.New, Key)
	mac.Write([]byte(header+body))
	signature := string(mac.Sum(nil))

	var token = b64.StdEncoding.EncodeToString([]byte(header+body+signature))

	return token

}

//преобразование токена в читабельный вид
func DeconcatToken(token string) types.Token {

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

//проверка токена на валидность
func CheckToken(token types.Token) bool {

	//подпись
	mac := hmac.New(sha256.New, Key)
	mac.Write([]byte(token.Head + token.Body))
	signature := mac.Sum(nil)

	return hmac.Equal(signature, []byte(token.Signature))
}