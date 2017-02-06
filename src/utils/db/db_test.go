package utils

import (
	"config"
	// "encoding/json"
	// log "github.com/Sirupsen/logrus"
	// mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "strings"
	"testing"
	"types"
)

// func TestMongoConn(t *testing.T) {
// 	var c *mgo.Collection
// 	var db *mgo.Database
// 	MongoDB := config.MongoDB_Conf{
// 		Host: "172.17.0.1",
// 		Port: "27017",
// 		User: "muser",
// 		Pass: "mpass",
// 		DB:   "fc",
// 	}
// 	mongo, err := MongoConn(MongoDB)
// 	if err == nil {
// 		db = mongo.DB("test")
// 		c = db.C("test")
// 	} else {
// 		log.Fatal(err)
// 	}
// 	info := `{"name": "www","age": 18,"tags":["a","b","c"]}`
// 	dec := json.NewDecoder(strings.NewReader(info))
// 	var v interface{}
// 	err = dec.Decode(&v)
// 	if err == nil {
// 		err = c.Insert(v)
// 	} else {
// 		log.Error(err)
// 	}
// 	query := c.Find(bson.M{"name": "www"})
// 	num, err := query.Count()
// 	result := make([]interface{}, 0)
// 	query.All(&result)
// 	log.Info("query num: ", num)
// 	log.Info("result : ", result)

// 	t.Log("pass")
// }

func TestMongoInsert(t *testing.T) {
	var user types.User
	config.LoadConfig()
	user.UserID = "userID"
	user.Password = "ppppassword"
	collection := "user"
	err := MongoInsert(collection, user)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Insert OK")
	}
}

func TestMongoFind(t *testing.T) {
	query := bson.M{"UserID": "userID"}
	result, err := MongoFind("user", query, nil)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

func TestMongoUpdate(t *testing.T) {
	query := bson.M{"UserID": "userID"}
	result, err := MongoFind("user", query, nil)
	if err != nil {
		t.Error(err)
	}
	// t.Log(result[0].(types.User).Password)
	t.Log(len(result))
	update := bson.M{"$set": bson.M{"EmailAddress": "mr@kfd.mememememem"}}

	err = MongoUpdate("user", query, update)
	if err != nil {
		t.Error(err)
	} else {
		result, _ := MongoFind("user", query, nil)
		t.Log(result)
	}
}

func TestMongoRemove(t *testing.T) {
	selector := bson.M{"UserID": "userID"}
	err := MongoRemove("user", selector)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Remove OK")
	}
}
