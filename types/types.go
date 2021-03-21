package types

//======================================================================================================================
//Типы главного модуля для чтения запроса и выдачи ответа
//======================================================================================================================
//Item для JsonAnswerBody
type JsonAnswerItem map[string]string

//тело ответа
type JsonAnswerBody struct {
	Items []JsonAnswerItem
}

//ошибки парсинга
type Errors []string

//структура запроса
type JsonRequest struct {
	Entity string //сущность (user, wish, group, badge, etc...)
	Id string //Идентификатор сущности (не обязательный)
	Action string //Действие (get, list, update, insert, etc...)
	Params map[string]string //Дополнительные параметры (page, sort, etc...) или поля entity (name, description, etc...)
}

//струкрута ответа
type JsonAnswer struct {
	Status string //статус (success, error)
	Answer JsonAnswerBody //тело ответа
	Response JsonRequest //запрашиваемые данные
	Errors Errors //ошибки запроса
}


//======================================================================================================================
//Типы User
//======================================================================================================================
//тип для пользователя
type User struct {
	Id int
	Email string
	Pass string
	Fio string
	Sex string
	Telegram string
	Instagram string
	Twitter string
	Facebook string
	Phone string
	Role string
	Avatar int
	Google string
}