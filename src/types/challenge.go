package types

type Challenge struct {
	ID         string
	Name       string
	TemplateID string
	Flag       string
	Services   []string
	NetworkID  string
	UserID     string
	Time       Time_struct
}

/*
{
  "ID": "xxx-xxx-xxx",
  "Name": "hello",
  "Template": "xxx-xxx-xxx",
  "Flag": "flag{xxx-xxx-xxx}",
  "Services": [
    "xxx-xxx-xxxA",
    "xxx-xxx-xxxB"
  ],
  "NetworkID": "xxx-xxx-xxxNet",
  "UserID": "xxx-xxx-xxx",
  "Time": {
    "CreateTime": 1484476298,
    "FinishTime": 1484476298
  }
}
*/
