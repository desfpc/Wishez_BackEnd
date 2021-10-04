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

// CheckErr кидает панику, если есть ошибка
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// MakeStringFromSQL конвертирует sql.NullString в string
func MakeStringFromSQL(str sql.NullString) string {
	if !str.Valid {
		return ""
	}
	return str.String
}

// MakeStringFromIntSQL конвертирует sql.NullInt64 в string
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

// NoRouteErrorAnswer ответ при ошибочном роуте
func NoRouteErrorAnswer() (types.Errors, int) {
	Errors := make(types.Errors,0)
	Errors = append(Errors, "Entity and/or action not found")
	return Errors, 404
}

// IsEmailValid проверка валидности строки email
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

// Escape аналог real_escape_strings
func Escape(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]

		escape = 0

		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
			break
		case '\n': /* Must be escaped for logs */
			escape = 'n'
			break
		case '\r':
			escape = 'r'
			break
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		case '"': /* Better safe than sorry */
			escape = '"'
			break
		case '\032': //十进制26,八进制32,十六进制1a, /* This gives problems on Win32 */
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}

// ParamFromJsonRequest выводит string переменную из JsonRequest-а
func ParamFromJsonRequest(params map[string]string, paramName string, errors types.Errors) (string, types.Errors, bool) {
	param, exists := params[paramName]
	if !exists {
		errors = append(errors, "No "+paramName)
		param = ""
	}
	param = Escape(param)
	return param, errors, exists
}