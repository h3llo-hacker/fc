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

type Service struct {
	ServiceName   string `bson:"ServiceName"`
	TargetPort    int    `bson:"TargetPort"`
	PublishedPort int    `bson:"PublishedPort"`
}

type UserChallenge struct {
	ChallengeID string    `bson:"ChallengeID"`
	TemplateID  string    `bson:"TemplateID"`
	UrlPrefix   string    `bson:"UrlPrefix"`
	Flag        string    `bson:"Flag"`
	FinishTime  time.Time `bson:"FinishTime"`
	CreateTime  time.Time `bson:"CreateTime"`
	State       string    `bson:"State"`
	Services    []Service `bson:"Services"`
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
	Challenges   []UserChallenge `bson:"Challenges"`
	Following    []string        `bson:"Following"`
	Followers    []string        `bson:"Followers"`
	Login        Login_struct    `bson:"Login"`
	Quota        int             `bson:"Quota"`
	Register     Register_struct `bson:"Register"`
	IsActive     bool            `bson:"IsActive"`
	WebSite      string          `bson:"WebSite",valid:"url"`
	Invite       Invite_struct   `bson:"Invite"`
	Rank         int             `bson:"Rank"`
}

func (user *User) ValidateUser() error {

	if user.UserName == "" || user.Password == "" || user.EmailAddress == "" {
		return errors.New("username or password or email cannot be empty.")
	}

	err := user.ValidateFormat()
	if err != nil {
		return err
	}

	return nil
}

func (user *User) ValidateFormat() error {
	// First check email address
	if user.EmailAddress != "" && !valid.IsEmail(user.EmailAddress) {
		return errors.New("illegal email format: " + user.EmailAddress)
	}

	// check website
	if user.WebSite != "" {
		if !valid.IsURL(user.WebSite) {
			return errors.New("illegal website format: " + user.WebSite)
		}
	}

	if user.UserURL != "" {
		user.UserURL = strings.ToLower(user.UserURL)
	}

	if user.UserName != "" {
		user.UserName = strings.TrimSpace(user.UserName)
		if len(user.UserName) < 6 || len(user.UserName) > 66 {
			return errors.New("illegal username length, must be [6, 66]")
		}
	}

	if len(user.Intro) > 423 {
		user.Intro = user.Intro[:423]
	}

	return nil
}
