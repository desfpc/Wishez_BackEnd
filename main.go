package main

import (
	"encoding/json"
	"github.com/desfpc/Wishez_Type"
	"github.com/desfpc/Wishez_User"
	"io/ioutil"
	"log"
	"net/http"
)

var errors types.Errors

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
func getArrBody(w http.ResponseWriter, r *http.Request) types.JsonRequest {
	var body = getBody(w, r)
	var arr types.JsonRequest
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
func answer(w http.ResponseWriter, status string, answer types.JsonAnswerBody, response types.JsonRequest, code int){

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	jsonAnswer := types.JsonAnswer{status, answer, response, errors}



	json.NewEncoder(w).Encode(jsonAnswer)
	errors = make(types.Errors,0)
}

//главный хандлер - все проверки и роутинг идут тут
func apiHandler(w http.ResponseWriter, r *http.Request) {

	resp := getArrBody(w, r)
	var anw types.JsonAnswerBody
	var err types.Errors

	var code = 200

	var status = "success"

	switch resp.Entity {
		case "user":
			anw, err = user.Route(resp)
			if len(err) > 0 {
				errors = err
			}
	}

	//если есть ошибки - ставим error status и не проводим обработку
	if len(errors) > 0 {
		status = "error"
		code = 500
	}

	answer(w, status, anw, resp, code)
}

//главная точнка входа - слушает все и выкидывает в apiHandler (коммент для проверки автодеплоя)
func main() {
	errors = make(types.Errors,0)
	http.HandleFunc("/", apiHandler)
	log.Printf("Wishez server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}