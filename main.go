package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//структура запроса
type JsonRequest struct {
	Entity string //сущность (user, wish, group, badge, etc...)
	Id string //Идентификатор сущности (не обязательный)
	Action string //Действие (get, list, update, insert, etc...)
	Params map[string]string //Дополнительные параметры (page, sort, etc...) или поля entity (name, description, etc...)
}

//Item для JsonAnswerBody
type JsonAnswerItem map[string]string

//тело ответа
type JsonAnswerBody struct {
	Items []JsonAnswerItem
}

//ошибки парсинга
type Errors []string

//струкрута ответа
type JsonAnswer struct {
	Status string //статус (success, error)
	Answer JsonAnswerBody //тело ответа
	Response JsonRequest //запрашиваемые данные
	Errors Errors //ошибки запроса
}

var errors Errors

//возвращение тела запроса в виде строки
func getBody(w http.ResponseWriter, r *http.Request) string  {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		errors = append(errors, "Error reading Body")
		return ""
	}

	//log.Printf("Body: "+string(body))
	return string(body)
}

//получкение массива из JSON запроса
func getArrBody(w http.ResponseWriter, r *http.Request) JsonRequest {
	var body = getBody(w, r)
	var arr JsonRequest
	err := json.Unmarshal([]byte(body), &arr)
	if err != nil {
		log.Printf("Error reading JSON from body: %v", err)
		errors = append(errors, "Error reading JSON from body: "+body)
	}
	var _, _ = json.Marshal(arr)
	//log.Printf("JsonBody: "+string(resp))
	return arr
}

//проверка запроса на авторизацию
func checkAuthorize(){

}

//вывод JSON ответа
func answer(w http.ResponseWriter, status string, answer JsonAnswerBody, response JsonRequest, code int){

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	jsonAnswer := JsonAnswer{status, answer, response, errors}



	json.NewEncoder(w).Encode(jsonAnswer)
	errors = make(Errors,0)
}

//главный хандлер - все проверки и роутинг идут тут
func apiHandler(w http.ResponseWriter, r *http.Request) {

	resp := getArrBody(w, r)
	var anw JsonAnswerBody

	var code = 200

	/*anw.Items = make([]JsonAnswerItem,0)
	item := make(JsonAnswerItem)
	item["Test"] = "Value"
	anw.Items = append(anw.Items, item)
	anw.Items = append(anw.Items, item)*/

	var status = "success"

	//если есть ошибки - ставим error status и не проводим обработку
	if(len(errors) > 0){
		status = "error"
		code = 500
	} else {



	}

	answer(w, status, anw, resp, code)
}

//главная точнка входа - слушает все и выкидывает в apiHandler
func main() {
	errors = make(Errors,0)
	http.HandleFunc("/", apiHandler)
	log.Printf("Wishez server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}