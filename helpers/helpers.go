package helpers

import (
	"database/sql"
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