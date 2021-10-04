package group

import (
	"database/sql"
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
	//TODO case "deleteUser":
	//	body, err = deleteUser(resp, auser)
	//TODO case "delete":
	//	body, err = deleteGroup(resp, auser)
	//TODO case "edit":
	//	body, err = editGroup(resp, auser)
	//TODO case "list":
	//	body, err = getGroupList(resp, auser)
	//TODO case "get":
	//	body, err = getGroup(resp, auser)
	default:
		err, code = helpers.NoRouteErrorAnswer()
	}

	return body, err, code
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

	//проверка на наличае ID группы
	var groupId string
	groupId, Errors, exist = helpers.ParamFromJsonRequest(params, "groupId", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на наличае ID пользователя
	var userId string
	userId, Errors, exist = helpers.ParamFromJsonRequest(params, "userId", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на наличае прав пользователя
	var right string
	right, Errors, exist = helpers.ParamFromJsonRequest(params, "right", Errors)
	if !exist {
		return body, Errors
	}

	//проверка на существование группы
	initDb()
	query := "SELECT * FROM group WHERE id = "+groupId
	results, err := dbres.Query(query)
	helpers.CheckErr(err)

	var group types.Group
	count := 0

	//перебираем результаты
	for results.Next() {
		count += 1
		//пробуем все запихнуть в user-а
		err = results.Scan(&group.Id, &group.AuthorId, &group.Name, &group.Visible, &group.OpenSum, &group.ClosedSum, &group.DateAdd)
		helpers.CheckErr(err)
	}

	if count == 0 {
		Errors = append(Errors, "No group with Id: "+groupId)
		return body, Errors
	}

	//проверяем права пользователя на возможэность добавить другого пользователя
	if !checkUserAdmin(auser, group) {
		Errors = append(Errors, "No admin rights to change group with Id: "+groupId)
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