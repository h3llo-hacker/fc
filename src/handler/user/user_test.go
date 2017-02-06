package handler

import (
	"testing"
	"time"
	"utils"

	"gopkg.in/mgo.v2/bson"
	"types"
)

func TestAddUser(t *testing.T) {
	var user types.User

	user.UserName = "大王111aa"
	user.Password = utils.Password("password")
	user.EmailAddress = "test6@test.com"
	user.Quota = 1
	user.IsActive = true
	user.Register = types.Register_struct{
		IP:     "8.8.8.8",
		Region: "China",
		System: types.System_struct{
			OS: "Windows",
			UA: "Chrome",
		},
		Date: time.Now(),
	}
	user.WebSite = "http://kfd.me"
	user.Intro = "一一"

	t.Log("Add a Validated user...")
	err := AddUser(user)
	if err != nil {
		t.Errorf("Add User Error: [%v]", err)
	}

	t.Log("Add a Unvalidated user...")
	user.EmailAddress = "mr@@kfd.me"
	err = AddUser(user)
	if err == nil {
		t.Errorf("Add Unvalidated Error %v", err)
	}
	t.Logf("Add User Error: [%v]", err)
}

func TestQueryUsersAll(t *testing.T) {
	emailAddr := "test2@test.com"
	users, err := QueryUsers(emailAddr, nil)
	if err != nil {
		t.Errorf("Query User Error: [%v]", err)
	}
	for _, user := range users {
		t.Log(user)
	}
}

func TestQueryUsersWithItem(t *testing.T) {
	emailAddr := "test2@test.com"
	items := []string{"UserID", "Username", "Password"}
	users, err := QueryUsers(emailAddr, items)
	if err != nil {
		t.Errorf("Query User Error: [%v]", err)
	}
	for _, user := range users {
		t.Log(user)
	}
}

func TestUpdateUser(t *testing.T) {
	emailAddr := "test@test.com"
	update := bson.M{"$set": bson.M{"EmailAddress": "test@test.com"}}
	err := UpdateUser(emailAddr, update)
	if err != nil {
		t.Errorf("Update User Errror: [%v]", err)
	} else {
		t.Log("Update User OK.")
	}
}

func TestRmUser(t *testing.T) {
	emailAddr := "test@test.com"
	err := RmUser(emailAddr)
	if err != nil {
		t.Errorf("Remove User Error: [%v]", err)
	} else {
		t.Log("Rm User OK.")
	}
}
