package handler

import (
	"config"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"testing"
)

func TestMongoConn(t *testing.T) {
	var c *mgo.Collection
	var db *mgo.Database
	MongoDB := config.MongoDB_Conf{
		Host: "172.17.0.1",
		Port: "27017",
		User: "muser",
		Pass: "mpass",
		DB:   "fc",
	}
	mongo, err := MongoConn(MongoDB)
	if err == nil {
		db = mongo.DB("test")
		c = db.C("test")
	} else {
		log.Fatal(err)
	}
	info := `{"name": "www","age": 18,"tags":["a","b","c"]}`
	dec := json.NewDecoder(strings.NewReader(info))
	var v interface{}
	err = dec.Decode(&v)
	if err == nil {
		err = c.Insert(v)
	} else {
		log.Error(err)
	}
	query := c.Find(bson.M{"name": "www"})
	num, err := query.Count()
	result := make([]interface{}, 0)
	query.All(&result)
	log.Info("query num: ", num)
	log.Info("result : ", result)

	t.Log("pass")
}
