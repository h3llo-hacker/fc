package handler

import (
	"errors"
	"fmt"
	"types"

	log "github.com/Sirupsen/logrus"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
	db "handler/db"
)

func AddUser(user types.User) error {
	if hasSameEmailAddr(user.EmailAddress) {
		return errors.New("Email Address Has Already Used.")
	}
	uid, _ := uuid.NewV4()
	user.UserID = fmt.Sprintf("%v", uid)
	C := "user"
	err := db.MongoInsert(C, user)
	if err != nil {
		log.Errorf("Add User Error %v", err)
	}
	return err
}

func RmUser(emailAddr string) error {
	if !hasSameEmailAddr(emailAddr) {
		return errors.New("User Email Not Found")
	}
	C := "user"
	selector := bson.M{"EmailAddress": emailAddr}
	err := db.MongoRemove(C, selector)
	if err != nil {
		log.Errorf("Remove User Error: [%v]", err)
	}
	return err
}

func UpdateUser(emailAddr string, update bson.M) error {
	if !hasSameEmailAddr(emailAddr) {
		return errors.New("User Email Not Found")
	}
	// update := bson.M{"$set": bson.M{"EmailAddress": "mr@kfd.me"}}
	C := "user"
	selector := bson.M{"EmailAddress": emailAddr}
	err := db.MongoUpdate(C, selector, update)
	return err
}

func QueryUsers(emailAddr string) ([]types.User, error) {
	C := "user"
	selector := bson.M{"EmailAddress": emailAddr}
	var users []types.User
	err := db.MongoFindUsers(C, selector, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func hasSameEmailAddr(emailAddr string) bool {
	C := "user"
	selector := bson.M{"EmailAddress": emailAddr}
	var users []types.User
	db.MongoFindUsers(C, selector, &users)
	if len(users) == 0 {
		return false
	}
	return true
}
