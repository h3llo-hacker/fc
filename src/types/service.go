package types

type Port struct {
	Src int
	Dst int
}

type Service struct {
	ID           string
	BelongToUser string
	TemplateID   string
	NetworkID    string
	Time         Time_struct
	Status       string
	Ports        []Port
}

/*
{
  "ID": "xxx-xxx-xxx",
  "BelongToChallenge": "xxx-xxx-xxx",
  "NetworkID":"xxx-xxx-xxx-Net",
  "TemplateID": "xxx-xxx-xxx",
  "Time": {
    "CreateTime": 1484476298,
    "FinishTime": 1484476298
  },
  "Status": "Running",
  "Ports": [
    {
      "Src": 30082,
      "Dst": 80
    },
    {
      "Src": 30083,
      "Dst": 22
    }
  ]
}
*/
