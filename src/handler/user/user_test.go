package handler

import (
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
	"handler/password"
	"types"
)

func TestAddUser(t *testing.T) {
	var user types.User

	user.Username = "testUser"
	user.Password = handler.Password("password")
	user.EmailAddress = "mr@kfd.me"
	user.Quota = 9
	user.Register = types.Register_struct{
		IP:     "8.8.8.8",
		Region: "China",
		System: types.System_struct{
			OS: "Windows",
			UA: "Chrome",
		},
		Date: time.Now(),
	}

	err := AddUser(user)
	if err != nil {
		t.Errorf("Add User Error: [%v]", err)
		return
	}
	t.Log("Add User OK.")

}

func TestRmUser(t *testing.T) {
	emailAddr := "mr@kfd.me"
	err := RmUser(emailAddr)
	if err != nil {
		t.Errorf("Remove User Error: [%v]", err)
	} else {
		t.Log("Rm User OK.")
	}
}

func TestQueryUsers(t *testing.T) {
	emailAddr := "test@test.com"
	users, err := QueryUsers(emailAddr)
	if err != nil {
		t.Errorf("Query User Error: [%v]", err)
	}
	for _, user := range users {
		t.Log(user.EmailAddress)
	}
}

func TestUpdateUser(t *testing.T) {
	emailAddr := "test@test.com"
	update := bson.M{"$set": bson.M{"EmailAddress": "mr@kfd.me"}}
	err := UpdateUser(emailAddr, update)
	if err != nil {
		t.Errorf("Update User Errror: [%v]", err)
	} else {
		t.Log("Update User OK.")
	}
}
