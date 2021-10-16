package group

import (
	"database/sql"
	"fmt"
	db "github.com/desfpc/Wishez_DB"
	helpers "github.com/desfpc/Wishez_Helpers"
	types "github.com/desfpc/Wishez_Type"
	"strconv"
)

var dbres *sql.DB

func initDb(){
	dbres = db.Db("", "")
}

// Route роутер Group
func Route(resp types.JsonRequest, auser types.User) (types.JsonAnswerBody, types.Errors, int) {
	var body types.JsonAnswerBody
	var err types.Errors
	var code = 200

	//запускаем обработку действий
	switch resp.Action {
	case "create":
		body, err = createGroup(resp, auser)
	case "addUser":
		body, err = addUser(resp, auser)
	case "deleteUser":
		body, err = deleteUser(resp, auser)
	case "delete":
		body, err = deleteGroup(resp, auser)
	case "edit":
		body, err = editGroup(resp, auser)
	case "list":
		body, err = getGroupList(resp, auser)
	//TODO case "get":
	//	body, err = getGroup(resp, auser)
	default:
		err, code = helpers.NoRouteErrorAnswer()
	}

	return body, err, code
}

// getGroupAndCheckUserAdmin функция объединяет получение группы и проверка админских прав на нее у текущего пользователя
func getGroupAndCheckUserAdmin(groupId string, auser types.User) (bool, types.Group, types.Errors) {
	Errors := make(types.Errors,0)
	exist, group := checkGroupExistById(groupId)
	if !exist {
		Errors = append(Errors, "No group with Id: "+groupId)
		return false, group, Errors
	}
	if checkUserAdmin(auser, group) {
		return true, group, Errors
	}
	Errors = append(Errors, "No admin rights to change group with Id: "+groupId)
	return false, group, Errors
}

// getGroupAndCheckUser функция объединяет получение группы и проверка участия в ней текущего пользователя
func getGroupAndCheckUser(groupId string, auser types.User) (bool, types.Group, types.Errors) {
	Errors := make(types.Errors,0)
	exist, group := checkGroupExistById(groupId)
	if !exist {
		Errors = append(Errors, "No group with Id: "+groupId)
		return false, group, Errors
	}
	if checkUser(auser, group) {
		return true, group, Errors
	}
	Errors = append(Errors, "No rights to view group with Id: "+groupId)
	return false, group, Errors
}

// checkGroupExistById проверка на существование группы по string Id, если группа есть, выводим дополнительно ее данные
func checkGroupExistById(groupId string) (bool, types.Group) {
	query := "SELECT * FROM group WHERE id = "+groupId
	results, err := dbres.Query(query)
	helpers.CheckErr(err)

	var group types.Group
	count := 0

	//перебираем результаты
	for results.Next() {
		count += 1
		//пробуем все запихнуть в group-у
		err = results.Scan(&group.Id, &group.AuthorId, &group.Name, &group.Visible, &group.OpenSum, &group.ClosedSum, &group.DateAdd)
		helpers.CheckErr(err)
	}

	if count == 0 {
		return false, group
	}
	return true, group
}

// checkUserAdmin проверка прав пользователя на административные действия над группой
func checkUserAdmin(auser types.User, group types.Group) bool {
	if auser.Id == group.AuthorId {
		return true
	}

	query := "SELECT COUNT(user_id) count WHERE user_id = ? AND group_id = ? AND right = 'admin'"
	results, err := dbres.Query(query)
	helpers.CheckErr(err)

	count := db.CheckCount(results)
	if count > 0 {
		return true
	}

	return false
}

// checkUser проверка прав пользователя на участие в группе
func checkUser(auser types.User, group types.Group) bool {
	if auser.Id == group.AuthorId {
		return true
	}

	query := "SELECT COUNT(user_id) count WHERE user_id = ? AND group_id = ?"
	results, err := dbres.Query(query)
	helpers.CheckErr(err)

	count := db.CheckCount(results)
	if count > 0 {
		return true
	}

	return false
}

// addUserToGroup создает запись пользователя в группе
func addUserToGroup(groupId int, userId int, right string) bool {
	if right != "admin" {
		right = "user"
	}

	_, err := dbres.Exec("INSERT INTO group_users (group_id, user_id, right, date_add) " +
		"VALUES (?, ?, ?, NOW())",
		groupId, userId, right)

	if err != nil {
		return false
	}

	return true
}

// getGroupList получение списка доступных групп
//
// предпологаемый json запроса:
// {"entity":"group","action":"list","params":{"groupType":"all","userId":1,"search":"подарок"}}
// entity string - сущность
// action string - действие
// params.groupType string - тип получаемых групп: строка из массива ['own','all']
// params.userId string - id пользователя, если надо получить его публичные группы (необязательно)
// params.search string - строка для поиска (необязательно)
func getGroupList(resp types.JsonRequest, auser types.User) (types.JsonAnswerBody, types.Errors) {
	var body types.JsonAnswerBody
	var params = resp.Params
	Errors := make(types.Errors,0)

	//проверка на наличие типа группы
	groupType, Errors, exist := helpers.ParamFromJsonRequest(params, "groupType", Errors)
	if !exist {
		return body, Errors
	}

	if groupType != "all" {
		groupType = "own"
	}

	//проверка на наличие ID пользователя для получения группы
	userId, Errors, existUser := helpers.ParamFromJsonRequest(params, "userId", Errors)
	userId = helpers.Escape(userId)

	//проверка на существование строки поиска
	search, Errors, existSearch := helpers.ParamFromJsonRequest(params, "search", Errors)
	search = helpers.Escape(search)

	//формируем запрос в БД
	query := "SELECT * FROM group WHERE "

	if groupType == "own" {
		query += "author = " + strconv.Itoa(auser.Id) + " "
	} else {
		query += "visible = 'public' "
		if existUser {
			query += "AND author = " + userId + " "
		}
	}
	if existSearch {
		query += "AND name LIKE ('%" + search + "%')"
	}

	initDb()
	results, err := dbres.Query(query)
	helpers.CheckErr(err)

	body.Items = make([]types.JsonAnswerItem,0)

	for results.Next() {
		var group types.Group
		err = results.Scan(&group.Id, &group.AuthorId, &group.Name, &group.Visible, &group.OpenSum, &group.ClosedSum, &group.DateAdd)
		helpers.CheckErr(err)
		item := make(types.JsonAnswerItem)
		item["Id"] = strconv.Itoa(group.Id)
		item["AuthorId"] = strconv.Itoa(group.AuthorId)
		item["Name"] = group.Name
		item["Visible"] = group.Visible
		item["OpenSum"] = fmt.Sprintf("%.2f", group.OpenSum)
		item["ClosedSum"] = fmt.Sprintf("%.2f", group.ClosedSum)
		item["DateAdd"] = group.DateAdd
		body.Items = append(body.Items, item)
	}

	return body, Errors
}

// deleteGroup удаление группы
//
// предпологаемый json запроса:
// {"entity":"group","action":"delete","params":{"groupId":"GroupId"}}
// entity string - сущность
// action string - действие
// params.groupId string - ID группы для удаления
func deleteGroup(resp types.JsonRequest, auser types.User) (types.JsonAnswerBody, types.Errors) {
	var body types.JsonAnswerBody
	var params = resp.Params
	Errors := make(types.Errors,0)

	//проверка на наличие ID группы
	groupId, Errors, exist := helpers.ParamFromJsonRequest(params, "groupId", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на существование группы и прав пользователя на ее изменение
	initDb()
	exist, _, Errors = getGroupAndCheckUserAdmin(groupId, auser)
	if !exist {
		return body, Errors
	}

	//удаляем группу и плачем
	_, err := dbres.Exec("DELETE FROM group WHERE id = ?",
		groupId)
	helpers.CheckErr(err)

	return body, Errors
}

// deleteUser удаление пользователя из группы
//
// предпологаемый json запроса:
// {"entity":"group","action":"deleteUser","params":{"groupId":"GroupId","userId":"UserId"}}
// entity string - сущность
// action string - действие
// params.groupId string - ID группы, куда нужно добавить пользователя
// params.userId string - ID пользователя, добавляемого в группу
func deleteUser(resp types.JsonRequest, auser types.User) (types.JsonAnswerBody, types.Errors) {
	var body types.JsonAnswerBody
	var params = resp.Params
	var exist bool
	Errors := make(types.Errors,0)

	//проверка на наличие ID группы
	var groupId string
	groupId, Errors, exist = helpers.ParamFromJsonRequest(params, "groupId", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на наличие ID пользователя
	var userId string
	userId, Errors, exist = helpers.ParamFromJsonRequest(params, "userId", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на существование группы и прав пользователя на ее изменение
	initDb()
	exist, _, Errors = getGroupAndCheckUserAdmin(groupId, auser)
	if !exist {
		return body, Errors
	}

	//удаляем пользователя из группы
	_, err := dbres.Exec("DELETE FROM group_users WHERE group_id = ? AND user_id = ?",
		groupId, userId)
	helpers.CheckErr(err)

	return body, Errors
}

// addUser добавление пользователя в группу
//
// предпологаемый json запроса:
// {"entity":"group","action":"addUser","params":{"groupId":"GroupId","userId":"UserId","right":"admin"}}
// entity string - сущность
// action string - действие
// params.groupId string - ID группы, куда нужно добавить пользователя
// params.userId string - ID пользователя, добавляемого в группу
// params.right string - права пользователя в группе: строка из массива ['admin','user']
func addUser(resp types.JsonRequest, auser types.User) (types.JsonAnswerBody, types.Errors) {
	var body types.JsonAnswerBody
	var params = resp.Params
	var exist bool
	Errors := make(types.Errors,0)

	//проверка на наличие ID группы
	var groupId string
	groupId, Errors, exist = helpers.ParamFromJsonRequest(params, "groupId", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на наличие ID пользователя
	var userId string
	userId, Errors, exist = helpers.ParamFromJsonRequest(params, "userId", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на наличие прав пользователя
	var right string
	right, Errors, exist = helpers.ParamFromJsonRequest(params, "right", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на существование группы и прав пользователя на ее изменение
	initDb()
	exist, _, Errors = getGroupAndCheckUserAdmin(groupId, auser)
	if !exist {
		return body, Errors
	}

	//добавляем пользователя в группу
	groupIdInt, err := strconv.Atoi(groupId)
	helpers.CheckErr(err)
	userIdInt, err := strconv.Atoi(userId)
	helpers.CheckErr(err)
	addedUser := addUserToGroup(groupIdInt, userIdInt, right)

	if !addedUser {
		Errors = append(Errors, "Error when adding user ("+userId+") to group ("+groupId+")")
		return body, Errors
	}

	item := make(types.JsonAnswerItem)
	item["UserId"] = userId
	item["GroupId"] = groupId
	item["right"] = right

	body.Items = make([]types.JsonAnswerItem,0)
	body.Items = append(body.Items, item)

	return body, Errors
}

// editGroup изменение группы
//
// предпологаемый json запроса:
// {"entity":"group","action":"edit","params":{"name":"GroupName","visible":"visibleString", "groupId":"GroupId"}}
// entity string - сущность
// action string - действие
// params.groupId string - ID группы для редактирования
// params.name string - новое наименование группы (необязательно)
// params.visible string - видимость группы: строка из массива ['hidden','normal','public'] (необязательно)
func editGroup(resp types.JsonRequest, auser types.User) (types.JsonAnswerBody, types.Errors) {
	var body types.JsonAnswerBody
	var params = resp.Params
	var exist bool
	Errors := make(types.Errors,0)

	//проверка на наличие ID группы
	var groupId string
	groupId, Errors, exist = helpers.ParamFromJsonRequest(params, "groupId", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на существование группы и прав пользователя на ее изменение
	initDb()
	exist, _, Errors = getGroupAndCheckUserAdmin(groupId, auser)
	if !exist {
		return body, Errors
	}

	//проверка на наличие наименования
	var name, existsName = params["name"]

	//проверка на наличие видимости
	var visible, existsVisible = params["visible"]

	//формируем запрос в зависимости от переданных значений на изменение
	if !existsName && ! existsVisible {
		Errors = append(Errors, "No values to change")
		return body, Errors
	}

	if existsName && existsVisible {
		_, err := dbres.Exec("UPDATE group SET name = ?, visible = ? WHERE id = ?",
			name, visible, groupId)
		helpers.CheckErr(err)
	} else if existsName {
		_, err := dbres.Exec("UPDATE group SET name = ? WHERE id = ?",
			name, groupId)
		helpers.CheckErr(err)
	} else {
		_, err := dbres.Exec("UPDATE group SET visible = ? WHERE id = ?",
			visible, groupId)
		helpers.CheckErr(err)
	}

	return body, Errors
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

	//проверка на наличие наименования
	var name, existsName = params["name"]
	if !existsName {
		Errors = append(Errors, "No name")
		return body, Errors
	}

	//проверка на наличие видимости
	var visible, existsVisible = params["visible"]
	if !existsVisible {
		Errors = append(Errors, "No visible")
		return body, Errors
	}

	//пробуем создать группу
	initDb()
	res, err := dbres.Exec("INSERT INTO group (id, author, name, visible, open_sum, closed_sum, date_add) " +
		"VALUES (null, ?, ?, ?, 0, 0, NOW())",
		auser.Id, name, visible)
	helpers.CheckErr(err)

	lastId, err := res.LastInsertId()
	helpers.CheckErr(err)

	addedUser := addUserToGroup(int(lastId), auser.Id, "admin")

	if !addedUser {
		res, err = dbres.Exec("DELETE FROM group WHERE id = ?",
			lastId)
		helpers.CheckErr(err)
	}

	item := make(types.JsonAnswerItem)
	item["Name"] = name
	item["Id"] = strconv.FormatInt(lastId, 10)

	body.Items = make([]types.JsonAnswerItem,0)
	body.Items = append(body.Items, item)

	return body, Errors
}