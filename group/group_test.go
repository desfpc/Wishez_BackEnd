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
	_, err := createGroup(request, auser)
	if len(err) > 0 {
		t.Error("Errors when creating group")
	}

	//TODO получение группы

	//TODO изменение группы

	//TODO удаление группы
}

/*func TestGetGroupAndCheckUserAdmin(t *testing.T) {
	isAdmin, group, errors := getGroupAndCheckUserAdmin("1")
}*/