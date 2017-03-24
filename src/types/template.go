package types

type Template struct {
	ID      string      `bson:"ID"`
	Name    string      `bson:"Name"`
	Level   string      `bson:"Level"`
	Tips    string      `bson:"Tips"`
	Tags    []string    `bson:"Tags"`
	Enable  bool        `bson:"Enable"`
	Content interface{} `bson:"Content"`
	Details string      `bson:"Details"`
}

/*
  "ID": "xxx-xxx-xxx",
  "Name": "templateName",
  "Content": "yaml file"
*/
