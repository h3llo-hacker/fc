package types

import (
	"time"
)

type System_struct struct {
	OS string `bson:"OS"`
	UA string `bson:"UA"`
}

type Register_struct struct {
	IP     string        `bson:"IP"`
	Region string        `bson:"Region"`
	System System_struct `bson:"System"`
	Date   time.Time     `bson:"Date"`
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
	Challenges   Challenge_types `bson:"Challenges"`
	EmailAddress string          `bson:"EmailAddress"`
	Following    []string        `bson:"Following"`
	Followers    []string        `bson:"Followers"`
	Login        Login_struct    `bson:"Login"`
	Password     string          `bson:"Password"`
	Quota        int             `bson:"Quota"`
	Username     string          `bson:"Username"`
	UserID       string          `bson:"UserID"`
	Register     Register_struct `bson:"Register"`
}

// type System_struct struct {
// 	OS string
// 	UA string
// }

// type Register_struct struct {
// 	IP     string
// 	Region string
// 	System System_struct
// 	Date   int
// }

// type UserChallenge struct {
// 	ChallengeID string
// 	TemplateID  string
// 	Flag        string
// 	FinishTime  int
// 	CreateTime  int
// }

// type Challenge_types struct {
// 	Failed    []UserChallenge
// 	Finished  []UserChallenge
// 	InProcess []UserChallenge
// }

// type Login_struct struct {
// 	LastLogins []Register_struct
// 	LoginTimes int
// }
// type User struct {
// 	Challenges   Challenge_types
// 	EmailAddress string
// 	Following    []string
// 	Followers    []string
// 	Login        Login_struct
// 	Password     string
// 	Quota        int
// 	Username     string
// 	UserID       string
// 	Register     Register_struct
// }
