package user

import (
	"config"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"types"
	"utils"
	db "utils/db"
	"utils/email"

	log "github.com/Sirupsen/logrus"
	pinyin "github.com/jmz331/gpinyin"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
)

type User types.User

const C = "user"

func AddUser(user types.User) (string, error) {
	if err := user.ValidateUser(); err != nil {
		return "", err
	}

	// check if email address was taken
	query := bson.M{"EmailAddress": user.EmailAddress}
	users, _ := db.MongoFind(C, query, bson.M{"_id": 1})
	if len(users) != 0 {
		return "", errors.New("Email Address Has Already Used.")
	}

	uid, _ := uuid.NewV4()
	user.UserID = fmt.Sprintf("%v", uid)
	user.UserNum = getUserNum()
	user.UserURL = generateURL(user.UserName, user.UserNum)
	if user.Invite.InvitedBy != "invite_off" {
		user.IsActive = true
	} else {
		user.IsActive = false
	}

	err := db.MongoInsert(C, user)
	if err != nil {
		log.Errorf("MongoInsert Error %v", err)
		return "", err
	}

	go func() {
		region := utils.IP2Region(user.Register.IP)
		err := updateUserRegion(user.UserID, region)
		if err != nil {
			log.Errorf("Update User Region Error: %v, User: [%v]", err, user.EmailAddress)
		}

		u := User(user)
		err = u.GenerateInvitecodes()
		if err != nil {
			log.Errorf("generateInvitecodes error: [%v], UserEmail: [%v]", err, user.EmailAddress)
		}
	}()

	return user.UserID, nil
}

func (user *User) RmUser() error {
	if !user.UserExist() {
		return fmt.Errorf("User Not Found.")
	}
	e := bson.M{"EmailAddress": user.EmailAddress}
	u := bson.M{"UserURL": user.UserURL}
	i := bson.M{"UserID": user.UserID}
	query := bson.M{"$or": []bson.M{e, u, i}}
	err := db.MongoRemove(C, query)
	if err != nil {
		log.Errorf("MongoRemove Error: [%v]", err)
		return fmt.Errorf("MongoRemove Error: [%v]", err)
	}
	return nil
}

func (user *User) UpdateUser(update bson.M) error {
	e := bson.M{"EmailAddress": user.EmailAddress}
	u := bson.M{"UserURL": user.UserURL}
	i := bson.M{"UserID": user.UserID}
	query := bson.M{"$or": []bson.M{e, u, i}}
	log.Debugf("update user debug: [%v] [%v]", query, update)
	err := db.MongoUpdate(C, query, update)
	if err != nil {
		return fmt.Errorf("update user err: %v", err)
	}
	return nil
}

func (user *User) QueryUserAll(items []string, limit, skip int) ([]types.User, error) {
	var (
		users []types.User
		query = bson.M{}
	)

	selector := make(bson.M, len(items))
	if items != nil {
		for _, item := range items {
			selector[item] = 1
		}
	}

	mongo, dbName, err := db.MongoConn()
	if err != nil {
		return users, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Find(query).Select(selector).Limit(limit).Skip(skip).All(&users)
	if err != nil {
		return nil, fmt.Errorf("MongoFindUsers Error: [%v]", err)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("404 user not found")
	}
	return users, nil
}

func (user *User) QueryUser(items []string) (types.User, error) {
	var (
		users []types.User
		query = bson.M{}
	)
	e := bson.M{"EmailAddress": user.EmailAddress}
	i := bson.M{"UserID": user.UserID}
	u := bson.M{"UserURL": user.UserURL}
	query = bson.M{"$or": []bson.M{e, u, i}}

	selector := make(bson.M, len(items))
	if items != nil {
		for _, item := range items {
			selector[item] = 1
		}
	} else {
		selector = nil
	}
	err := db.MongoFindUsers(C, query, selector, &users)
	if err != nil {
		return types.User{}, err
	}
	if len(users) == 0 {
		return types.User{}, errors.New("user not found")
	}
	return users[0], nil
}

func (user *User) QueryUserWithSelector(selector bson.M) (types.User, error) {
	var (
		users []types.User
		query = bson.M{}
	)
	e := bson.M{"EmailAddress": user.EmailAddress}
	i := bson.M{"UserID": user.UserID}
	u := bson.M{"UserURL": user.UserURL}
	query = bson.M{"$or": []bson.M{e, u, i}}

	err := db.MongoFindUsers(C, query, selector, &users)
	if err != nil {
		return types.User{}, err
	}
	if len(users) == 0 {
		return types.User{}, errors.New("user not found")
	}
	return users[0], nil
}

func (user *User) QueryUsersRaw(items []string, limit, skip int) ([]interface{}, error) {
	var (
		userMap []interface{}
		query   = bson.M{}
	)
	selector := make(bson.M, len(items)+1)
	selector["_id"] = 0
	if items != nil {
		for _, item := range items {
			selector[item] = 1
		}
	}
	mongo, dbName, err := db.MongoConn()
	if err != nil {
		return nil, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Find(query).Select(selector).Limit(limit).Skip(skip).All(&userMap)
	if err != nil {
		return nil, err
	}
	if len(userMap) == 0 {
		return nil, fmt.Errorf("404 user not found")
	}
	return userMap, nil
}

func (user *User) QueryUserRaw(selector bson.M) (interface{}, error) {
	var (
		users []interface{}
		query = bson.M{}
	)
	e := bson.M{"EmailAddress": user.EmailAddress}
	i := bson.M{"UserID": user.UserID}
	u := bson.M{"UserURL": user.UserURL}
	query = bson.M{"$or": []bson.M{e, u, i}}

	selector["_id"] = 0

	users, err := db.MongoFind(C, query, selector)
	if err != nil {
		return "", err
	}
	if len(users) == 0 {
		return "", errors.New("404")
	}
	return users[0], nil
}

func userUrlExist(userUrl string) bool {
	query := bson.M{"UserURL": userUrl}
	users, _ := db.MongoFind(C, query, bson.M{"_id": 1})
	if len(users) == 0 {
		return false
	}
	return true
}

func getUserNum() int64 {
	num, err := db.MongoCount(C)
	if err != nil {
		return 0
	}
	return int64(num)
}

func generateURL(userName string, userNum int64) string {
	URL := pinyin.ConvertToPinyinString(userName, "-", pinyin.PINYIN_WITHOUT_TONE)
	if !userUrlExist(URL) {
		return URL
	}
	source := rand.NewSource(userNum)
	randNum := rand.New(source)
	uniqNum := randNum.Intn(99)
	URL += "-" + strconv.Itoa(uniqNum)
	for userUrlExist(URL) {
		URL = generateURL(URL, userNum)
	}
	return URL
}

func GetUserID(str, t string) (string, error) {
	var users []types.User
	query := bson.M{}

	switch strings.ToLower(t) {
	case "email":
		query = bson.M{"EmailAddress": str}
	case "url":
		query = bson.M{"UserUrl": str}
	case "num":
		query = bson.M{"UserNum": str}
	}

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
	err := db.MongoUpdate(C, query, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) CheckLogin() (string, bool) {
	if !user.UserExist() {
		return "", false
	}
	items := []string{"Password", "UserID", "UserName"}
	u, err := user.QueryUser(items)
	if err != nil {
		log.Errorf("Query User Error: [%v]", err)
		return "", false
	}
	encPassword := utils.Password(user.Password)
	if encPassword != u.Password {
		return "fail", false
	}
	return u.UserID, true
}

// update login times and login infomation
func (user *User) UpdateUserLogin(login types.Register_struct) error {
	update := bson.M{"$push": bson.M{"Login.LastLogins": login}}
	user.UpdateUser(update)
	update = bson.M{"$inc": bson.M{"Login.LoginTimes": 1}}
	user.UpdateUser(update)
	selector := bson.M{"Login.LoginTimes": 1}
	u, e := user.QueryUserRaw(selector)
	if e != nil {
		return e
	}
	loginTimes := u.(bson.M)["Login"].(bson.M)["LoginTimes"].(int)
	if loginTimes > 20 {
		update := bson.M{"$pop": bson.M{"Login.LastLogins": -1}}
		user.UpdateUser(update)
	}
	return nil
}

func (user *User) UserExist() bool {
	e := bson.M{"EmailAddress": user.EmailAddress}
	i := bson.M{"UserID": user.UserID}
	u := bson.M{"UserURL": user.UserURL}
	query := bson.M{"$or": []bson.M{e, u, i}}
	users, _ := db.MongoFind(C, query, bson.M{"_id": 1})
	if len(users) == 0 {
		return false
	}
	return true
}

func (user *User) QueryUserChallenges(states []string, limit, skip int) ([]types.UserChallenge, error) {

	if !user.UserExist() {
		return nil, fmt.Errorf("User Not Found.")
	}

	e := bson.M{"EmailAddress": user.EmailAddress}
	i := bson.M{"UserID": user.UserID}
	u := bson.M{"UserURL": user.UserURL}
	query := bson.M{"$or": []bson.M{e, u, i}}

	selector := bson.M{"Challenges": 1}

	mongo, dbName, err := db.MongoConn()
	if err != nil {
		return nil, err
	}
	users := make([]types.User, 0)
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Find(query).Select(selector).All(&users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}

	// get the challenges whitch satisfied the states provided
	temp_user_challenges := make([]types.UserChallenge, 0)
	if states[0] == "all" {
		temp_user_challenges = users[0].Challenges
	} else {
		for _, uc := range users[0].Challenges {
			if contains(states, uc.State) {
				temp_user_challenges = append(temp_user_challenges, uc)
			}
		}
	}

	var userChallenges []types.UserChallenge
	// reversed user challenges
	RUCS := reverseChallenges(temp_user_challenges)
	if skip == 0 && limit == 0 {
		userChallenges = RUCS
	} else {
		start := skip
		end := skip + limit
		if start >= len(RUCS) {
			return nil, nil
		}
		if end >= len(RUCS) {
			end = len(RUCS)
		}
		userChallenges = RUCS[start:end]
	}

	return userChallenges, nil

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func reverseChallenges(cs []types.UserChallenge) []types.UserChallenge {
	for i, j := 0, len(cs)-1; i < j; i, j = i+1, j-1 {
		cs[i], cs[j] = cs[j], cs[i]
	}
	return cs
}

func (user *User) Active(active bool) error {
	if !user.UserExist() {
		return fmt.Errorf("User Not Found.")
	}

	update := bson.M{"$set": bson.M{"IsActive": active}}
	err := user.UpdateUser(update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) RestPwd() error {
	if !user.UserExist() {
		return fmt.Errorf("User Not Found.")
	}

	// generate reset code and update database
	uid, _ := uuid.NewV4()
	resetCode := strings.ToUpper(fmt.Sprintf("%v", uid)[9:23])
	s := bson.M{"ResetPwd.Code": resetCode,
		"ResetPwd.Expire": time.Now().Add(30 * time.Minute)}
	i := bson.M{"ResetPwd.Times": 1}
	update := bson.M{"$set": s, "$inc": i}
	err := user.UpdateUser(update)
	if err != nil {
		return err
	}

	qu, _ := user.QueryUser([]string{"UserName"})

	// Send email
	urlParam := user.EmailAddress + ":" + resetCode
	b64_urlParam := base64.StdEncoding.EncodeToString([]byte(urlParam))
	log.Debugln(b64_urlParam, utils.BB64(b64_urlParam))

	receiver := types.ValidateEmail{
		EmailAddr: user.EmailAddress,
		UserName:  qu.UserName,
		ClickURL:  config.Conf.Mail.Site + "/resetpwd.php?" + utils.BB64(b64_urlParam),
	}

	if err := email.SendResetPwdEmail(receiver); err != nil {
		log.Errorf("Send Email Error: [%v]", err)
		log.Errorf("Trying Again...")
		return email.SendResetPwdEmail(receiver)
	}

	return nil
}

func (user *User) DoRestPwd(passwd string) (string, error) {
	if !user.UserExist() {
		return "", fmt.Errorf("User Not Found.")
	}

	qu, err := user.QueryUser([]string{"ResetPwd", "UserID"})
	if err != nil {
		return "", err
	}
	if qu.ResetPwd.Code != user.ResetPwd.Code {
		return "", fmt.Errorf("Code didn't match.")
	}
	if time.Now().After(qu.ResetPwd.Expire) {
		return "", fmt.Errorf("Code Expired.")
	}

	go func() {
		newPwd := utils.Password(passwd)
		now := time.Now()
		u := bson.M{"Password": newPwd, "ResetPwd.Expire": now}
		update := bson.M{"$set": u}
		err = user.UpdateUser(update)
		if err != nil {
			log.Errorf("ResetPassword.UpdateUser Error: [%v]", err)
		}
	}()

	return qu.UserID, nil
}

func (user *User) SendVerifyEmail() error {
	if !user.UserExist() {
		return fmt.Errorf("User Not Found.")
	}

	qu, _ := user.QueryUser([]string{"UserName", "Verify"})
	if qu.Verify.Verification {
		return fmt.Errorf("You have already verified your email.")
	}

	rt := qu.Verify.RequestTime.Add(1 * time.Minute)
	if !time.Now().After(rt) {
		return fmt.Errorf("Wait 1 minute and try.")
	}

	// generate verify code and update database
	uid, _ := uuid.NewV4()
	Code := strings.ToUpper(fmt.Sprintf("%v", uid)[9:23])
	s := bson.M{"Verify.Code": Code, "Verify.RequestTime": time.Now()}
	update := bson.M{"$set": s}
	err := user.UpdateUser(update)
	if err != nil {
		return err
	}

	// Send email
	urlParam := user.EmailAddress + ":" + Code
	b64_urlParam := base64.StdEncoding.EncodeToString([]byte(urlParam))
	log.Debugln(b64_urlParam, utils.BB64(b64_urlParam))
	receiver := types.ValidateEmail{
		EmailAddr: user.EmailAddress,
		UserName:  qu.UserName,
		ClickURL:  config.Conf.Mail.Site + "/verification.php?" + utils.BB64(b64_urlParam),
	}

	if err := email.SendVerifyEmail(receiver); err != nil {
		log.Errorf("Send Email Error: [%v]", err)
		log.Errorf("Trying Again...")
		return email.SendVerifyEmail(receiver)
	}

	return nil
}

func (user *User) VerifyEmail() error {
	if !user.UserExist() {
		return fmt.Errorf("User Not Found.")
	}

	qu, err := user.QueryUser([]string{"Verify"})
	if err != nil {
		return err
	}
	if qu.Verify.Code != user.Verify.Code {
		return fmt.Errorf("Code didn't match.")
	}

	if qu.Verify.Verification {
		return fmt.Errorf("You've already verified your email.")
	}

	u := bson.M{"Verify.Verification": true}
	update := bson.M{"$set": u}
	err = user.UpdateUser(update)
	if err != nil {
		log.Errorf("VerifyEmail.UpdateUser Error: [%v]", err)
		return err
	}

	return nil
}
