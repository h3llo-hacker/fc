package utils

import (
	"config"
	// "encoding/json"
	"fmt"
	// "os"
	// "strings"

	log "github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"types"
)

var mongoConn *mgo.Session

func MongoConn(MongoDB config.MongoDB_Conf) (*mgo.Session, error) {
	if mongoConn != nil {
		return mongoConn, nil
	}

	if config.Conf.MongoDB.Host == "" {
		config.LoadConfig()
		MongoDB = config.Conf.MongoDB
	}
	user := MongoDB.User
	pass := MongoDB.Pass
	host := MongoDB.Host
	port := MongoDB.Port
	db := MongoDB.DB
	if host == "" || port == "" || db == "" {
		log.Fatal("Host or port or db is nil")
	}
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, pass, host, port, db)
	if user == "" {
		url = fmt.Sprintf("mongodb://%s:%s/%s", host, port, db)
	}
	log.Debugln("MongoDB Conn:", url)
	mongo, err := mgo.Dial(url)
	if err != nil {
		log.Errorf("Mongo Conn Error: [%v], Mongo ConnUrl: [%v]", err, url)
		return &mgo.Session{}, err
	}
	mongoConn = mongo
	return mongoConn, nil
}

func MongoInsert(C string, data interface{}) error {
	MongoDB := config.Conf.MongoDB

	mongo, err := MongoConn(MongoDB)
	if err != nil {
		log.Fatal(err)
	}
	db := mongo.DB(MongoDB.DB)
	collection := db.C(C)
	err = collection.Insert(data)
	if err != nil {
		return err
	}
	return nil
}

func MongoRemove(C string, selector bson.M) error {
	MongoDB := config.Conf.MongoDB

	mongo, err := MongoConn(MongoDB)
	if err != nil {
		log.Fatal(err)
	}
	db := mongo.DB(MongoDB.DB)
	collection := db.C(C)
	err = collection.Remove(selector)
	if err != nil {
		return err
	}
	return nil
}

func MongoFind(C string, query bson.M) ([]interface{}, error) {
	MongoDB := config.Conf.MongoDB

	mongo, err := MongoConn(MongoDB)
	if err != nil {
		return nil, err
	}
	db := mongo.DB(MongoDB.DB)
	collection := db.C(C)
	result := make([]interface{}, 0)
	collection.Find(query).All(&result)
	return result, nil
}

func MongoFindUsers(C string, query bson.M, result *[]types.User) error {
	MongoDB := config.Conf.MongoDB
	mongo, err := MongoConn(MongoDB)
	if err != nil {
		return err
	}
	db := mongo.DB(MongoDB.DB)
	collection := db.C(C)
	err = collection.Find(query).All(result)
	return err
}

func MongoUpdate(C string, selector bson.M, update interface{}) error {
	MongoDB := config.Conf.MongoDB

	mongo, err := MongoConn(MongoDB)
	if err != nil {
		log.Fatal(err)
	}
	db := mongo.DB(MongoDB.DB)
	collection := db.C(C)
	err = collection.Update(selector, update)
	if err != nil {
		return err
	}
	return nil
}
