package types

type System_struct struct {
	OS string `json:"OS"`
	UA string `json:"UA"`
}

type Register_struct struct {
	IP     string        `json:"IP"`
	Region string        `json:"Region"`
	System System_struct `json:"System"`
	Date   int           `json:"Date"`
}

type UserChallenge struct {
	ChallengeID string `json:"ChallengeID"`
	TemplateID  string `json:"TemplateID"`
	Flag        string `json:"Flag"`
	FinishTime  int    `json:"FinishTime"`
	CreateTime  int    `json:"CreateTime"`
}

type Challenge_types struct {
	Failed    []UserChallenge
	Finished  []UserChallenge
	InProcess []UserChallenge
}

type Login_struct struct {
	LastLogins []Register_struct `json:"LastLogins"`
	LoginTimes int               `json:"LoginTimes"`
}
type User struct {
	Challenges   Challenge_types `json:"Challenges"`
	EmailAddress string          `json:"EmailAddress"`
	Following    []string        `json:"Following"`
	Followers    []string        `json:"Followers"`
	Login        Login_struct    `json:"Login"`
	Password     string          `json:"Password"`
	Quota        int             `json:"Quota"`
	Username     string          `json:"Username"`
	UserID       string          `json:"UserID"`
	Register     Register_struct `json:"Register"`
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
