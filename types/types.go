package types

import "database/sql"

// Token Тип токенов
type Token struct {
	Head string
	Body string
	Signature string
}
// JsonAnswerItem - Item для JsonAnswerBody
type JsonAnswerItem map[string]string

// JsonAnswerBody тело ответа
type JsonAnswerBody struct {
	Items []JsonAnswerItem
}

// Errors ошибки парсинга
type Errors []string

// JsonRequest структура запроса
type JsonRequest struct {
	Entity string //сущность (user, wish, group, badge, etc...)
	Id string //Идентификатор сущности (не обязательный)
	Action string //Действие (get, list, update, insert, etc...)
	Params map[string]string //Дополнительные параметры (page, sort, etc...) или поля entity (name, description, etc...)
}

// JsonAnswer струкрута ответа
type JsonAnswer struct {
	Status string //статус (success, error)
	Answer JsonAnswerBody //тело ответа
	Response JsonRequest //запрашиваемые данные
	Errors Errors //ошибки запроса
}

// User тип для пользователя
type User struct {
	Id int
	Email string
	Pass string
	Fio string
	Sex string
	Telegram sql.NullString
	Instagram sql.NullString
	Twitter sql.NullString
	Facebook sql.NullString
	Phone sql.NullString
	Role string
	Avatar sql.NullInt64
	Google sql.NullString
	CreatedAt string
}

// Group тип списка желаний
type Group struct {
	Id int
	Author int
	Name string
	Visible string
	OpenSum float32
	ClosedSum float32
	DateAt string
}

// TokenBody тип тела токена
type TokenBody map[string]string