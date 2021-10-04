package group

import (
	helpers "github.com/desfpc/Wishez_Helpers"
	types "github.com/desfpc/Wishez_Type"
)

// Route роутер Group
func Route(resp types.JsonRequest, auser types.User) (types.JsonAnswerBody, types.Errors, int) {
	var body types.JsonAnswerBody
	var err types.Errors
	var code = 200

	//запускаем обработку действий
	switch resp.Action {
	case "create":
		body, err = createGroup(resp, auser)
	default:
		err, code = helpers.NoRouteErrorAnswer()
	}

	return body, err, code
}

// createGroup создание нового листа желаний
//
// предпологаемый json запроса:
// {"entity":"group","action":"add","params":{"name":"GroupName","visible":"visibleString"}}
// entity string - сущность
// action string - действие
// params.name string - наименование новой группы
// params.visible string - видимость группы: строка из массива ['hidden','normal','public']
func createGroup(resp types.JsonRequest, auser types.User) (types.JsonAnswerBody, types.Errors) {

	var body types.JsonAnswerBody
	var params = resp.Params
	Errors := make(types.Errors,0)

	//проверка на наличае наименования
	var name, existsName = params["name"]
	if !existsName {
		Errors = append(Errors, "No name")
		return body, Errors
	}

	//проверка на наличае видимости
	var visible, existsVisible = params["visible"]
	if !existsVisible {
		Errors = append(Errors, "No visible")
		return body, Errors
	}

	

	return body, Errors
}