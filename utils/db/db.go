package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/h3llo-hacker/fc/config"
	"github.com/h3llo-hacker/fc/types"

	log "github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongoConn *mgo.Session

func MongoConn() (*mgo.Session, string, error) {
	if config.Conf.MongoDB.DB == "" {
		config.LoadConfig()
	}
	if mongoConn != nil {
		if mongoConn.Ping() == nil {
			return mongoConn, config.Conf.MongoDB.DB, nil
		}
	}

	MongoDB := config.Conf.MongoDB
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
	mongo, err := mgo.DialWithTimeout(url, 3*time.Second)
	if err != nil {
		log.Errorf("Mongo Conn Error: [%v], Mongo ConnUrl: [%v]",
			err, url)
		errTextReturn := fmt.Sprintf("Mongo Conn Error: [%v]", err)
		return &mgo.Session{}, "", errors.New(errTextReturn)
	}
	mongoConn = mongo
	return mongoConn, db, nil
}

func MongoInsert(C string, data interface{}) error {
	mongo, dbName, err := MongoConn()
	if err != nil {
		log.Error(err)
		return err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Insert(data)
	if err != nil {
		return err
	}
	return nil
}

func MongoRemove(C string, selector bson.M) error {
	mongo, dbName, err := MongoConn()
	if err != nil {
		log.Error(err)
		return err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Remove(selector)
	if err != nil {
		return err
	}
	return nil
}

func MongoFind(C string, query, selector bson.M) ([]interface{}, error) {
	mongo, dbName, err := MongoConn()
	if err != nil {
		return nil, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	result := make([]interface{}, 0)
	err = collection.Find(query).Select(selector).All(&result)
	return result, err
}

func MongoFindUsers(C string, query, selector bson.M, result *[]types.User) error {
	mongo, dbName, err := MongoConn()
	if err != nil {
		return err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Find(query).Select(selector).All(result)
	return err
}

func MongoUpdate(C string, selector bson.M, update interface{}) error {
	mongo, dbName, err := MongoConn()
	if err != nil {
		log.Error(err)
		return err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	err = collection.Update(selector, update)
	if err != nil {
		return err
	}
	return nil
}

func MongoCount(C string) (int, error) {
	mongo, dbName, err := MongoConn()
	if err != nil {
		return 0, err
	}
	db := mongo.DB(dbName)
	collection := db.C(C)
	count, err := collection.Count()
	return count, err
}
