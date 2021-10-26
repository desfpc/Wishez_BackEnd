package group

import (
	types "github.com/desfpc/Wishez_Type"
	"testing"
)

func TestWorkVsGroup(t *testing.T) {

	//создание группы
	request := types.JsonRequest{
		Entity: "group",
		Id:     "",
		Action: "add",
		Params: make(map[string]string),
	}
	auser := types.User{
		Id: 1,
		Email: "desfpc@gmail.com",
		Role: "user",
	}

	request.Params["name"] = "Тестовая группа 1"
	request.Params["visible"] = "hidden"
	answer, err := createGroup(request, auser)
	if len(err) > 0 {
		t.Error("Errors when creating group")
	}
	groupId := answer.Items[0]["Id"]

	//получение группы
	request = types.JsonRequest{
		Entity: "group",
		Id:     "",
		Action: "get",
		Params: make(map[string]string),
	}
	request.Params["groupId"] = groupId
	answer, err = getGroup(request, auser)
	if len(err) > 0 {
		t.Error("Errors when getting group")
	}
	answerGroupId := answer.Items[0]["Id"]
	if answerGroupId != groupId {
		t.Error("Wrong GroupId")
	}

	//изменение группы
	request = types.JsonRequest{
		Entity: "group",
		Id:     "",
		Action: "edit",
		Params: make(map[string]string),
	}
	request.Params["groupId"] = groupId
	request.Params["name"] = "Измененная группа"
	request.Params["visible"] = "public"
	answer, err = editGroup(request, auser)
	if len(err) > 0 {
		t.Error("Errors when updating group")
	}

	//TODO удаление группы
}

/*func TestGetGroupAndCheckUserAdmin(t *testing.T) {
	isAdmin, group, errors := getGroupAndCheckUserAdmin("1")
}*/