package utils

import (
	"config"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "os"
	"strings"
)

func Conn(MongoDB config.MongoDB_Conf) (*mgo.Session, error) {
	user := MongoDB.User
	pass := MongoDB.Pass
	host := MongoDB.Host
	port := MongoDB.Port
	db := config.C.MongoDB.DB
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, pass, host, port, db)
	if user == "" {
		url = fmt.Sprintf("mongodb://%s:%s/%s", host, port, db)
	}
	log.Debug("MongoDB Conn:", url)
	mongo, err := mgo.Dial(url)
	if err == nil {
		return mongo, nil
	} else {
		log.Error("Mongo Conn Error: ", err)
		return &mgo.Session{}, err
	}
}

func TestDB() {
	var c *mgo.Collection
	var db *mgo.Database

	mongo, err := Conn(config.C.MongoDB)
	if err == nil {
		db = mongo.DB("test")
		c = db.C("test")
	} else {
		log.Error(err)
	}
	info := `{"name": "wrfly","age": 18,"tags":["a","b","c"]}`
	dec := json.NewDecoder(strings.NewReader(info))
	var v interface{}
	err = dec.Decode(&v)
	if err == nil {
		err = c.Insert(v)
	} else {
		log.Error(err)
	}
	query := c.Find(bson.M{"name": "wrfly"})
	num, err := query.Count()
	result := make([]interface{}, 0)
	log.Errorln(query.All(&result))
	log.Info("query num: ", num)
	log.Info("result : ", result)

	log.Error("test done")
}
