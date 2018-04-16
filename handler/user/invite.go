package user

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/h3llo-hacker/fc/config"
	db "github.com/h3llo-hacker/fc/utils/db"

	// valid "github.com/asaskevich/govalidator"
	// pinyin "github.com/jmz331/gpinyin"
	"github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func (user *User) GenerateInvitecodes() error {
	num := config.Conf.InviteCodes
	log.Infof("Fenerate %d invitecodes for [%v]", num, user.UserID)
	return user.AddInvitecodes(num)
}

func (user *User) AddInvitecodes(num int) error {
	codes := make([]string, num)
	inviteCodes := make([]string, num)

	for i := 0; i < num; i++ {
		uid, _ := uuid.NewV4()
		s := fmt.Sprintf("%v", uid)[9:23]
		codes[i] = strings.Replace(s, "-", "", -1)
	}
	for i, sc := range codes {
		hex := strconv.FormatInt(user.UserNum, 16)

		nCode := fmt.Sprintf("%v%v", sc[:len(sc)-len(hex)], hex)
		var Code string
		for i := 0; i < 3; i++ {
			s := 4 * i
			e := s + 4
			Code += nCode[s:e]
			Code += "-"
		}
		inviteCodes[i] = Code[:len(Code)-1]
	}

	update := bson.M{"$pushAll": bson.M{"Invite.InviteCodes": inviteCodes}}
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

	query := bson.M{"Invite.InviteCodes": strings.ToLower(inviteCode)}
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
