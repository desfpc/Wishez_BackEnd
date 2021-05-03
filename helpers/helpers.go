package helpers

import (
	"database/sql"
	"github.com/desfpc/Wishez_Type"
	"strconv"
)

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
func AuthErrorAnswer(authorizedError bool, expiredError bool) types.Errors {
	Errors := make(types.Errors,0)
	if authorizedError {
		Errors = append(Errors, "Authorization Required")
	}else if authorizedError {
		Errors = append(Errors, "Access Token is Expired")
	}

	return Errors
}