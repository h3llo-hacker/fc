package handler

import (
	"config"
	// "encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "os"
	// "strings"
)

func MongoConn(MongoDB config.MongoDB_Conf) (*mgo.Session, error) {
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
	log.Info("MongoDB Conn:", url)
	mongo, err := mgo.Dial(url)
	if err == nil {
		return mongo, nil
	} else {
		log.Error("Mongo Conn Error: ", err)
		return &mgo.Session{}, err
	}
}
