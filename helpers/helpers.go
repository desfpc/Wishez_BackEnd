package helpers

import (
	"database/sql"
	"github.com/desfpc/Wishez_Type"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MakeStringFromSQL(str sql.NullString) string {
	if !str.Valid {
		return ""
	}
	return str.String
}

func MakeStringFromIntSQL(str sql.NullInt64) string {
	if !str.Valid {
		return ""
	}
	return strconv.FormatInt(str.Int64, 10)
}

// AuthErrorAnswer ответ при ошибке авторизации или протухании токена
func AuthErrorAnswer(authorizedError bool, expiredError bool) (types.Errors, int) {
	code := 200
	Errors := make(types.Errors,0)
	if authorizedError {
		Errors = append(Errors, "Authorization Required")
		code = 401
	}else if expiredError {
		Errors = append(Errors, "Access Token is Expired")
		code = 401
	}
	return Errors, code
}

func NoRouteErrorAnswer() (types.Errors, int) {
	Errors := make(types.Errors,0)
	Errors = append(Errors, "Entity and/or action not found")
	return Errors, 404
}

func IsEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
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