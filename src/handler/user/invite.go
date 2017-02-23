package user

import (
	"config"
	"errors"
	"fmt"

	// valid "github.com/asaskevich/govalidator"
	// pinyin "github.com/jmz331/gpinyin"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
	// "gopkg.in/mgo.v2/bson"
	// log "github.com/Sirupsen/logrus"
	db "utils/db"
)

func (user *User) GenerateInvitecodes(num int) error {
	codes := make([]string, num)

	for i := 0; i < num; i++ {
		uid, _ := uuid.NewV4()
		codes[i] = fmt.Sprintf("%v", uid)
	}

	update := bson.M{"$pushAll": bson.M{"Invite.InviteCodes": codes}}
	err := user.UpdateUser(update)
	if err != nil {
		return err
	}
	return nil
}

func GetInvitedBy(inviteCode string) (string, error) {
	if config.Conf.InviteMode == false {
		return "invite_off", nil
	}

	query := bson.M{"Invite.InviteCodes": inviteCode}
	selector := bson.M{}
	selector["UserID"] = 1
	user, err := db.MongoFind(C, query, selector)
	if err != nil {
		return "", err
	}
	if len(user) == 0 {
		return "", errors.New(fmt.Sprintf("Illegal InviteCode [%v]", inviteCode))
	}
	return user[0].(bson.M)["UserID"].(string), nil
}

func (user *User) RemoveInviteCode(code string) error {
	update := bson.M{"$pull": bson.M{"Invite.InviteCodes": code}}
	err := user.UpdateUser(update)
	if err != nil {
		return err
	}
	return nil
}
