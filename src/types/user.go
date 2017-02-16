package types

import (
	"errors"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
)

type System_struct struct {
	OS string `bson:"OS"`
	UA string `bson:"UA"`
}

type Invite_struct struct {
	InviteCodes []string `bson:"InviteCodes"`
	InvitedBy   string   `bson:"InvitedBy"`
}

type Register_struct struct {
	IP     string        `bson:"IP",valid:"ipv4"`
	Region string        `bson:"Region"`
	System System_struct `bson:"System"`
	Date   time.Time     `bson:"Date"` // MongoDB is ISODate (UTC)
}

type UserChallenge struct {
	ChallengeID string `bson:"ChallengeID"`
	TemplateID  string `bson:"TemplateID"`
	Flag        string `bson:"Flag"`
	FinishTime  int    `bson:"FinishTime"`
	CreateTime  int    `bson:"CreateTime"`
}

type Challenge_types struct {
	Failed    []UserChallenge `bson:"Failed"`
	Finished  []UserChallenge `bson:"Finished"`
	InProcess []UserChallenge `bson:"InProcess"`
}

type Login_struct struct {
	LastLogins []Register_struct `bson:"LastLogins"`
	LoginTimes int               `bson:"LoginTimes"`
}
type User struct {
	UserID       string          `bson:"UserID"`
	UserNum      int             `bson:"UserNum"`
	UserName     string          `bson:"UserName"`
	UserURL      string          `bson:"UserURL"`
	Password     string          `bson:"Password"`
	Intro        string          `bson:"Intro"`
	EmailAddress string          `bson:"EmailAddress",valid:"email"`
	Challenges   Challenge_types `bson:"Challenges"`
	Following    []string        `bson:"Following"`
	Followers    []string        `bson:"Followers"`
	Login        Login_struct    `bson:"Login"`
	Quota        int             `bson:"Quota"`
	Register     Register_struct `bson:"Register"`
	IsActive     bool            `bson:"IsActive"`
	WebSite      string          `bson:"WebSite",valid:"url"`
	Invite       Invite_struct   `bson:"Invite"`
}

func (user *User) ValidateUser() error {
	// First check email address
	if !valid.IsEmail(user.EmailAddress) {
		return errors.New("wrong email format:" + user.EmailAddress)
	}

	// check website
	if user.WebSite != "" {
		if !valid.IsURL(user.WebSite) {
			return errors.New("wrong website format")
		}
	}

	if user.UserURL != "" {
		user.UserURL = strings.ToLower(user.UserURL)
	}

	if user.UserName == "" || user.Password == "" || user.EmailAddress == "" {
		return errors.New("username or password or email cannot be empty.")
	}

	user.UserName = strings.TrimSpace(user.UserName)
	if len(user.UserName) < 6 || len(user.UserName) > 66 {
		return errors.New("make sure [6 < username < 66]")
	}

	if len(user.Intro) > 423 {
		user.Intro = user.Intro[:423]
	}

	return nil
}
