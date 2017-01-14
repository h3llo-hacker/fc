package types

type System_struct struct {
	OS string `json:"OS"`
	UA string `json:"UA"`
}

type Register_struct struct {
	IP     string        `json:"IP"`
	Region string        `json:"Region"`
	System System_struct `json:"System"`
	Time   int64         `json:"Time"`
}

type Challenge struct {
	ChallengeID string
	ContainerID string
	Finishtime  int
	StartTime   int
}

type Challenge_types struct {
	Failed    []Challenge_struct
	Finished  []Challenge_struct
	InProcess []Challenge_struct
}

type Login_struct struct {
	LastLogins []Register_struct
	LoginTimes int
}
type User struct {
	Challenges   Challenge_types
	EmailAddress string `json:EmailAddress`
	Following    []string
	Followers    []string
	Login        Login_struct
	Password     string          `json:Password`
	Quota        int             `json:Quota`
	Username     string          `json:"Username"`
	UserID       string          `json:UserID`
	Register     Register_struct `json:"Register"`
}
