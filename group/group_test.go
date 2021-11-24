package group

import (
	"fmt"
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
		Id:     groupId,
		Action: "get",
		Params: make(map[string]string),
	}
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
		Id:     groupId,
		Action: "edit",
		Params: make(map[string]string),
	}
	request.Params["name"] = "Измененная группа"
	request.Params["visible"] = "public"
	answer, err = editGroup(request, auser)
	if len(err) > 0 {
		t.Error("Errors when updating group")
	}

	//TODO addUser Action: "addUser"

	//getGroupList Action: "list"
	request = types.JsonRequest{
		Entity: "group",
		Id:     "",
		Action: "list",
		Params: make(map[string]string),
	}
	request.Params["groupType"] = "own"
	answer, err = getGroupList(request, auser)
	if len(err) > 0 {
		t.Error("Errors when getting group list: " + err[0])
	}
	fmt.Println("Groups count:")
	t.Log(len(answer.Items))
	if answer.Items[(len(answer.Items) - 1)]["Name"] != "Измененная группа" {
		t.Error("Errors when getting group list")
	}


	//TODO getUserList Action: "userList"

	//TODO deleteUser Action: "deleteUser"

	//удаление группы
	request = types.JsonRequest{
		Entity: "group",
		Id:     groupId,
		Action: "delete",
		Params: make(map[string]string),
	}
	answer, err = deleteGroup(request, auser)
	if len(err) > 0 {
		t.Error("Errors when deleting group")
	}
}