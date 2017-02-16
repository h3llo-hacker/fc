package user

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"types"
	"utils"

	log "github.com/Sirupsen/logrus"
	valid "github.com/asaskevich/govalidator"
	pinyin "github.com/jmz331/gpinyin"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
	db "utils/db"
)

type User types.User

func AddUser(user types.User) error {
	if err := user.ValidateUser(); err != nil {
		return err
	}
	if emailAddrExist(user.EmailAddress) {
		return errors.New("Email Address Has Already Used.")
	}
	uid, _ := uuid.NewV4()
	user.UserID = fmt.Sprintf("%v", uid)
	user.UserNum = getUserNum()
	user.UserURL = generateURL(user.UserName, user.UserNum)
	user.IsActive = false

	err := db.MongoInsert("user", user)
	if err != nil {
		log.Errorf("MongoInsert Error %v", err)
		return err
	}

	go func() {
		region := utils.IP2Region(user.Register.IP)
		err := updateUserRegion(user.UserID, region)
		if err != nil {
			log.Errorf("Update User Region Error: %v, User: [%v]", err, user.EmailAddress)
		}

		u := User(user)
		num := 5 // 5 invite code initially
		err = u.GenerateInvitecodes(num)
		if err != nil {
			log.Errorf("generateInvitecodes error: [%v], UserEmail: [%v]", err, user.EmailAddress)
		}
	}()

	return nil
}

func (user *User) RmUser() error {
	if !emailAddrExist(user.EmailAddress) {
		return errors.New("User Email Not Found")
	}
	C := "user"
	query := bson.M{"EmailAddress": user.EmailAddress}
	err := db.MongoRemove(C, query)
	if err != nil {
		log.Errorf("MongoRemove Error: [%v]", err)
		return err
	}
	return nil
}

func (user *User) UpdateUser(update bson.M) error {
	C := "user"
	e := bson.M{"EmailAddress": user.EmailAddress}
	u := bson.M{"UserURL": user.UserURL}
	i := bson.M{"UserID": user.UserID}
	query := bson.M{"$or": []bson.M{e, u, i}}
	err := db.MongoUpdate(C, query, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) QueryUser(items []string) ([]types.User, error) {
	var (
		users []types.User
		query = bson.M{}
		C     = "user"
	)
	e := bson.M{"EmailAddress": user.EmailAddress}
	i := bson.M{"UserID": user.UserID}
	u := bson.M{"UserURL": user.UserURL}
	query = bson.M{"$or": []bson.M{e, u, i}}

	if user.EmailAddress == "" && user.UserURL == "" {
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

func (user *User) QueryUserRaw(items []string) ([]interface{}, error) {
	var (
		userMap []interface{}
		query   = bson.M{}
		C       = "user"
	)
	if valid.IsEmail(user.EmailAddress) {
		query = bson.M{"EmailAddress": user.EmailAddress}
	} else {
		query = bson.M{"UserURL": user.UserURL}
	}
	if user.EmailAddress == "" && user.UserURL == "" {
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
	if len(userMap) == 0 {
		return nil, errors.New("404")
	}
	return userMap, nil
}

func emailAddrExist(emailAddr string) bool {
	C := "user"
	query := bson.M{"EmailAddress": emailAddr}
	users, _ := db.MongoFind(C, query, bson.M{"_id": 1})
	if len(users) == 0 {
		return false
	}
	return true
}

func userUrlExist(userUrl string) bool {
	C := "user"
	query := bson.M{"UserURL": userUrl}
	users, _ := db.MongoFind(C, query, bson.M{"_id": 1})
	if len(users) == 0 {
		return false
	}
	return true
}

func getUserNum() int {
	C := "user"
	num, err := db.MongoCount(C)
	if err != nil {
		return 0
	}
	return num
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

func GetUserID(str, t string) (string, error) {
	query := bson.M{}
	switch strings.ToLower(t) {
	case "email":
		query = bson.M{"EmailAddress": str}
	case "url":
		query = bson.M{"UserUrl": str}
	case "num":
		query = bson.M{"UserNum": str}
	}
	C := "user"
	var users []types.User
	db.MongoFindUsers(C, query, bson.M{}, &users)
	switch len(users) {
	case 0:
		return "", errors.New("404")
	case 1:
		return users[0].UserID, nil
	default:
		return "", errors.New("409")
	}
	return "", nil
}

func updateUserRegion(uid, region string) error {
	set := bson.M{"Register.Region": region}
	update := bson.M{"$set": set}
	query := bson.M{"UserID": uid}
	err := db.MongoUpdate("user", query, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) CheckLogin() bool {
	items := []string{"Password"}
	u, err := user.QueryUser(items)
	if err != nil || len(u) == 0 {
		return false
	}
	if utils.Password(user.Password) == u[0].Password {
		return true
	}
	return false
}

func (user *User) UpdateUserLogin(login types.Register_struct) error {
	update := bson.M{"$push": bson.M{"Login.LastLogins": login}}
	user.UpdateUser(update)
	update = bson.M{"$inc": bson.M{"Login.LoginTimes": 1}}
	user.UpdateUser(update)
	items := []string{"Login.LoginTimes"}
	u, e := user.QueryUserRaw(items)
	if e != nil {
		return e
	}
	loginTimes := u[0].(bson.M)["Login"].(bson.M)["LoginTimes"].(int)
	if loginTimes > 20 {
		update := bson.M{"$pop": bson.M{"Login.LastLogins": -1}}
		user.UpdateUser(update)
	}
	return nil
}
