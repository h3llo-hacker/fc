package types

import (
	"time"
)

type Time_struct struct {
	CreateTime time.Time `bson:"CreateTime"`
	FinishTime time.Time `bson:"FinishTime"`
}

type Challenge struct {
	ID         string      `bson:"ID"`
	Name       string      `bson:"Name"`
	TemplateID string      `bson:"TemplateID"`
	Flag       string      `bson:"Flag"`
	StackID    string      `bson:"StackID"`
	UserID     string      `bson:"UserID"`
	Time       Time_struct `bson:"Time"`
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
