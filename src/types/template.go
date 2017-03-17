package types

type Template struct {
	ID      string      `bson:"ID"`
	Name    string      `bson:"Name"`
	Display bool        `bson:"Display"`
	Content interface{} `bson:"Content"`
}

/*
  "ID": "xxx-xxx-xxx",
  "Name": "templateName",
  "Content": "yaml file"
*/
