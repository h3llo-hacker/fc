package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"types"

	log "github.com/Sirupsen/logrus"
	valid "github.com/asaskevich/govalidator"
	pinyin "github.com/jmz331/gpinyin"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
	db "utils/db"
)

func AddUser(user types.User) error {
	if !ValidateUser(&user) {
		return errors.New("User Format Error.")
	}
	if emailAddrExist(user.EmailAddress) {
		return errors.New("Email Address Has Already Used.")
	}
	uid, _ := uuid.NewV4()
	user.UserID = fmt.Sprintf("%v", uid)
	user.UserNum = getUserNum()
	user.UserURL = generateURL(user.UserName, user.UserNum)
	user.IsActive = false
	C := "user"
	err := db.MongoInsert(C, user)
	if err != nil {
		log.Errorf("MongoInsert Error %v", err)
		return err
	}
	return nil
}

func RmUser(emailAddr string) error {
	if !emailAddrExist(emailAddr) {
		return errors.New("User Email Not Found")
	}
	C := "user"
	query := bson.M{"EmailAddress": emailAddr}
	err := db.MongoRemove(C, query)
	if err != nil {
		log.Errorf("MongoRemove Error: [%v]", err)
		return err
	}
	return nil
}

func UpdateUser(emailAddr string, update bson.M) error {
	if !emailAddrExist(emailAddr) {
		return errors.New("User Email Not Found")
	}
	// update := bson.M{"$set": bson.M{"EmailAddress": "mr@kfd.me"}}
	C := "user"
	query := bson.M{"EmailAddress": emailAddr}
	err := db.MongoUpdate(C, query, update)
	if err != nil {
		return err
	}
	return nil
}

func QueryUsers(emailAddr string, items []string) ([]types.User, error) {
	var users []types.User

	C := "user"
	query := bson.M{"EmailAddress": emailAddr}
	if emailAddr == "" {
		query = bson.M{}
	}
	selector := make(bson.M, len(items))
	if items != nil {
		for _, item := range items {
			selector[item] = 1
		}
	}
	err := db.MongoFindUsers(C, query, selector, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func QueryUsersRaw(emailAddr string, items []string) ([]interface{}, error) {
	var userMap []interface{}

	C := "user"
	query := bson.M{"EmailAddress": emailAddr}
	if emailAddr == "" {
		query = bson.M{}
	}
	selector := make(bson.M, len(items))
	selector["_id"] = 0
	if items != nil {
		for _, item := range items {
			selector[item] = 1
		}
	}
	userMap, err := db.MongoFind(C, query, selector)
	if err != nil {
		return nil, err
	}
	return userMap, nil
}

func emailAddrExist(emailAddr string) bool {
	C := "user"
	query := bson.M{"EmailAddress": emailAddr}
	var users []types.User
	db.MongoFindUsers(C, query, bson.M{}, &users)
	if len(users) == 0 {
		return false
	}
	return true
}

func userUrlExist(userUrl string) bool {
	C := "user"
	query := bson.M{"UserURL": userUrl}
	var users []types.User
	db.MongoFindUsers(C, query, bson.M{"UserURL": 1}, &users)
	if len(users) == 0 {
		return false
	}
	return true
}

func ValidateUser(user *types.User) bool {
	if !valid.IsEmail(user.EmailAddress) {
		return false
	}
	if user.WebSite != "" {
		if !valid.IsURL(user.WebSite) {
			return false
		}
	}
	if user.UserName == "" || user.Password == "" || user.EmailAddress == "" {
		return false
	}
	if len(user.Intro) > 423 {
		user.Intro = user.Intro[:423]
	}
	user.UserName = strings.Trim(user.UserName, "#")
	if len(user.UserName) > 66 {
		user.UserName = user.UserName[:66]
	}
	return true
}

func getUserNum() int {
	C := "user"
	query := bson.M{}
	var users []types.User
	db.MongoFindUsers(C, query, bson.M{}, &users)
	// log.Debugf("There are %d users in total.", len(users))
	return len(users)
}

func generateURL(userName string, userNum int) string {
	URL := pinyin.ConvertToPinyinString(userName, "-", pinyin.PINYIN_WITHOUT_TONE)
	if !userUrlExist(URL) {
		return URL
	}
	source := rand.NewSource(int64(userNum))
	randNum := rand.New(source)
	uniqNum := randNum.Intn(999)
	URL += "-" + strconv.Itoa(uniqNum)
	for userUrlExist(URL) {
		URL = generateURL(URL, userNum)
	}
	return URL
}
