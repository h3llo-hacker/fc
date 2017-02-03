package handler

import (
	"errors"
	"fmt"
	"types"

	log "github.com/Sirupsen/logrus"
	valid "github.com/asaskevich/govalidator"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
	db "utils/db"
)

func AddUser(user types.User) error {
	if !ValidateUser(&user) {
		return errors.New("User Format Error.")
	}
	if EmailAddrExist(user.EmailAddress) {
		return errors.New("Email Address Has Already Used.")
	}
	uid, _ := uuid.NewV4()
	user.UserID = fmt.Sprintf("%v", uid)
	user.UserNum = getUserNum()
	user.IsActive = false
	C := "user"
	err := db.MongoInsert(C, user)
	if err != nil {
		log.Errorf("Add User Error %v", err)
	}
	return err
}

func RmUser(emailAddr string) error {
	if !EmailAddrExist(emailAddr) {
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
	if !EmailAddrExist(emailAddr) {
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

func EmailAddrExist(emailAddr string) bool {
	C := "user"
	selector := bson.M{"EmailAddress": emailAddr}
	var users []types.User
	db.MongoFindUsers(C, selector, &users)
	if len(users) == 0 {
		return false
	}
	return true
}

func ValidateUser(user *types.User) bool {
	if !valid.IsEmail(user.EmailAddress) {
		return false
	}
	if !valid.IsURL(user.WebSite) {
		return false
	}
	if len(user.Intro) > 423 {
		user.Intro = user.Intro[:423]
	}
	return true
}

func getUserNum() int {
	C := "user"
	selector := bson.M{}
	var users []types.User
	db.MongoFindUsers(C, selector, &users)
	// log.Debugf("There are %d users in total.", len(users))
	return len(users)
}
