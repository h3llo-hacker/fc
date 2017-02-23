package routes

import (
	"handler/template"
	"handler/user"
	db "utils/db"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"ping": "pong",
	})
}

func validateAuth(auth string) bool {
	C := "auth"
	selector := bson.M{}
	auths, err := db.MongoFind(C, nil, selector)
	if err != nil {
		return false
	}
	for _, Auth := range auths {
		if auth == Auth.(bson.M)["token"] {
			return true
		}
	}
	return false
}

func validateUser(uid string) bool {
	return user.ExistUser(uid)
}

func validateTemplate(templateID string) bool {
	return template.ValidateTemplate(templateID)
}
